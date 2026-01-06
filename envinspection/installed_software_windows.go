//go:build windows

package envinspection

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

type listSubKeysError struct {
	e       error
	message string
}

func (l *listSubKeysError) Unwrap() error {
	return l.e
}

func (l *listSubKeysError) Error() string {
	return fmt.Sprintf("listSubkeys: %s", l.e)
}

func listSubKeys(ctx context.Context, key registry.Key, path string) ([]string, error) {
	k, e := registry.OpenKey(key, path, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
	if e != nil {
		return nil, &listSubKeysError{e, "open key failed"}
	}
	defer k.Close()
	ki, e := k.Stat()
	if e != nil {
		return nil, &listSubKeysError{e, "get stat failed"}
	}
	skList, e := k.ReadSubKeyNames(int(ki.SubKeyCount))
	if e != nil {
		return nil, &listSubKeysError{e, "read subkey names failed"}
	}
	var r = make([]string, 0, len(skList))
	for _, s := range skList {
		r = append(r, filepath.Join(path, s))
	}
	return r, nil
}

func getWindowsVersion() model.Component {
	return model.Component{
		CompName:    "Windows",
		CompVersion: fmt.Sprintf("%d.%d.%d", windows.RtlGetVersion().MajorVersion, windows.RtlGetVersion().MinorVersion, windows.RtlGetVersion().BuildNumber),
	}
}

func listInstalledSoftwareWindows(ctx context.Context) ([]model.DependencyItem, error) {
	var logger = logctx.Use(ctx).Sugar()
	var searchDirs = []struct {
		Key  registry.Key
		Path string
	}{
		{registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall"},
		{registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall"},
		{registry.LOCAL_MACHINE, "SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall"},
		{registry.CURRENT_USER, "SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall"},
	}
	var r []model.DependencyItem
	for _, rKey := range searchDirs {
		paths, e := listSubKeys(ctx, rKey.Key, rKey.Path)
		if e != nil {
			return nil, e
		}
		for _, p := range paths {
			k, e := registry.OpenKey(rKey.Key, p, registry.READ)
			if e != nil {
				continue
			}
			displayName, _, _ := k.GetStringValue("DisplayName")
			displayVersion, _, _ := k.GetStringValue("DisplayVersion")
			_ = k.Close()
			if displayName == "" {
				continue
			}
			r = append(r, model.DependencyItem{
				Component: model.Component{
					CompName:    displayName,
					CompVersion: displayVersion,
				},
			})
		}
	}
	cmd := exec.CommandContext(ctx, "powershell")
	cmd.Stdin = bytes.NewReader([]byte(`
Get-AppxPackage | Select-Object @{Name="Name";Expression={ "[APPX] " + $_.Name }}, Version

`))
	var output bytes.Buffer
	cmd.Stdout = &output
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	e := cmd.Run()
	if e != nil {
		logger.Errorf("powershell failed: %s", e)
		return r, nil
	}
	stdoutBytes := output.Bytes()
	stderrBytes := stderr.Bytes()
	converter := consoleDecoder()
	if converter != nil {
		stdoutBytes, _ = converter.Bytes(stdoutBytes)
		stderrBytes, _ = converter.Bytes(stderrBytes)
	}
	stderrText := strings.TrimSpace(string(stderrBytes))
	if stderrText != "" {
		logger.Errorf("powershell failed: %s", stderrText)
	}
	var pattern = regexp.MustCompile(`^\[APPX\] (\S+)\s*(\S+)`)
	for line := range strings.SplitSeq(strings.TrimSpace(string(stdoutBytes)), "\n") {
		line = strings.TrimSpace(line)
		matches := pattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		r = append(r, model.DependencyItem{
			Component: model.Component{
				CompName:    matches[1],
				CompVersion: matches[2],
			},
		})
	}
	return r, nil
}

func consoleDecoder() *encoding.Decoder {
	var cp, e = windows.GetConsoleCP()
	if e != nil {
		return nil
	}
	switch cp {
	case 936:
		return simplifiedchinese.GBK.NewDecoder()
	case 932:
		return japanese.ShiftJIS.NewDecoder()
	case 950:
		return traditionalchinese.Big5.NewDecoder()
	case 949:
		return korean.EUCKR.NewDecoder()
	case 437:
		return charmap.CodePage437.NewDecoder()
	case 850:
		return charmap.CodePage850.NewDecoder()
	case 852:
		return charmap.CodePage852.NewDecoder()
	case 866:
		return charmap.CodePage866.NewDecoder()
	case 1250:
		return charmap.Windows1250.NewDecoder()
	case 1251:
		return charmap.Windows1251.NewDecoder()
	case 1252:
		return charmap.Windows1252.NewDecoder()
	default:
		return nil // assume UTF-8
	}
}

func listPendingPatch(ctx context.Context) []string {
	var logger = logctx.Use(ctx).Sugar()
	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()
	cmd := exec.CommandContext(ctx, "powershell")
	cmd.Stdin = bytes.NewReader([]byte(`
$session = New-Object -ComObject Microsoft.Update.Session
$searcher = $session.CreateUpdateSearcher()
$searcher.Online = $true
$resultAvailable = $searcher.Search("IsInstalled=0")
function Get-KBString($update) {
    if ($update.KBArticleIDs -and $update.KBArticleIDs.Count -gt 0) {
        return "KB$($update.KBArticleIDs[0])"
    }
    return ""
}
function Print-Updates($updates) {
    foreach ($update in $updates) {
        $kb = Get-KBString $update
        Write-Host "KBNumber: $kb"
    }
}
Print-Updates $resultAvailable.Updates


`))
	var kbPattern = regexp.MustCompile(`KBNumber: (KB\d+)`)
	var output bytes.Buffer
	cmd.Stdout = &output
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	e := cmd.Run()
	if e != nil {
		logger.Errorf("powershell failed: %s", e)
		return nil
	}
	stdoutBytes := output.Bytes()
	stderrBytes := stderr.Bytes()
	converter := consoleDecoder()
	if converter != nil {
		stdoutBytes, _ = converter.Bytes(stdoutBytes)
		stderrBytes, _ = converter.Bytes(stderrBytes)
	}
	stderrText := strings.TrimSpace(string(stderrBytes))
	if stderrText != "" {
		logger.Errorf("powershell failed: %s", stderrText)
	}
	var kbNumbers []string
	for line := range strings.SplitSeq(strings.TrimSpace(string(stdoutBytes)), "\n") {
		var kbNumber = kbPattern.FindStringSubmatch(line)
		if kbNumber != nil {
			kbNumbers = append(kbNumbers, kbNumber[1])
		}
	}
	return kbNumbers
}

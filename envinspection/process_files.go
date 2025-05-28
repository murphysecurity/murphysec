package envinspection

import (
	"archive/zip"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/repeale/fp-go"
	"github.com/samber/lo"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type HashResult struct {
	MD5    [md5.Size]byte
	SHA1   [sha1.Size]byte
	SHA256 [sha256.Size]byte
}

// Retrieve information about files related to running processes
func listProcessFiles(ctx context.Context) ([]model.Module, error) {
	LOG := logctx.Use(ctx).Sugar()
	LOG.Info("Fetching process-related file information")

	var moduleList []model.Module
	procDir := "/proc"

	// Check if the /proc directory exists
	if _, err := os.Stat(procDir); os.IsNotExist(err) {
		LOG.Warnf("Directory %s does not exist, skipping process file collection", procDir)
		return moduleList, nil
	}

	entries, err := os.ReadDir(procDir)
	if err != nil {
		LOG.Warnf("Failed to read /proc directory: %v", err)
		return moduleList, err
	}

	for _, entry := range entries {
		if !entry.IsDir() || !isNumeric(entry.Name()) {
			continue
		}

		pid := entry.Name()
		pidN, _ := strconv.Atoi(pid)
		if pidN == os.Getpid() {
			continue
		}
		exePath := filepath.Join(procDir, pid, "exe")
		cmdlinePath := filepath.Join(procDir, pid, "cmdline")
		statusPath := filepath.Join(procDir, pid, "status")

		// Get process name
		processName, err := os.Readlink(exePath)
		if err != nil {
			continue
		}

		// Filter system processes
		if isSystemProcess(pid, filepath.Base(processName), cmdlinePath, statusPath) {
			continue
		}

		// Create a module for each process
		moduleName := filepath.Base(processName)
		modulePath := processName
		LOG.Infof("Processing module: %s (PID: %s)", moduleName, modulePath)

		var hashes []HashResult

		// Process regular ELF executable
		target, err := os.Readlink(exePath)
		if err != nil {
			continue
		}

		// Calculate exe file information
		if isRelevantFile(target) {
			if hash, e := calculateFileHashes(target); e == nil {
				hashes = append(hashes, hash)
			}
		}

		// Get file descriptor information
		fdDir := filepath.Join(procDir, pid, "fd")
		fdFiles := getFdFiles(fdDir)

		// Get dynamic library information for the executable
		lddFiles := getLddFiles(target)

		// Merge file paths
		allFiles := append(fdFiles, lddFiles...)

		// Process the list of files
		for _, file := range allFiles {
			if isRelevantFile(file) {
				if strings.HasSuffix(file, ".jar") || strings.HasSuffix(file, ".war") {
					// Calculate hash for the JAR/WAR file itself
					hash, e := calculateFileHashes(file)
					if e == nil {
						hashes = append(hashes, hash)
					}

					// Extract and process JAR/WAR files
					_, hashes0 := extractJarFilesWithHashes(ctx, file)
					hashes = append(hashes, hashes0...)
				} else {
					// Process regular files
					if hash, e := calculateFileHashes(file); e == nil {
						hashes = append(hashes, hash)
					}
				}
			}
		}
		// Create a module for this process
		processModule := model.Module{
			ModuleName:     moduleName,
			ModulePath:     modulePath + fmt.Sprintf("[%s]", pid),
			PackageManager: "process",
			Dependencies:   []model.DependencyItem{},
			ScanStrategy:   model.ScanStrategyNormal,
			MD5Hashes:      lo.Uniq(fp.Map(func(it HashResult) model.MD5Hash { return it.MD5 })(hashes)),
			SHA1Hashes:     lo.Uniq(fp.Map(func(it HashResult) model.SHA1Hash { return it.SHA1 })(hashes)),
			SHA256Hashes:   lo.Uniq(fp.Map(func(it HashResult) model.SHA256Hash { return it.SHA256 })(hashes)),
		}

		moduleList = append(moduleList, processModule)
	}

	return moduleList, nil
}

// Calculate MD5 and SHA256 hashes for a file
func calculateFileHashes(filePath string) (HashResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return HashResult{}, err
	}
	defer func() { _ = file.Close() }()

	// Use streaming hash calculation
	return computeHashesFromReader(file)

}

// Get file paths from file descriptors
func getFdFiles(fdDir string) []string {
	var files []string
	entries, err := os.ReadDir(fdDir)
	if err != nil {
		return files
	}

	for _, entry := range entries {
		fdPath := filepath.Join(fdDir, entry.Name())

		// Check if it's a device file
		if isDeviceFile(fdPath) {
			continue
		}

		target, err := os.Readlink(fdPath)
		if err == nil {
			files = append(files, target)
		}
	}
	return files
}

// Get dynamic library paths using ldd
func getLddFiles(exePath string) []string {
	var files []string

	cmd := exec.Command("ldd", exePath)
	output, err := cmd.Output()
	if err != nil {
		return files
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) > 2 && strings.HasPrefix(parts[2], "/") {
			// Skip anon_inode
			if strings.HasPrefix(parts[2], "anon_inode:") {
				continue
			}

			realPath, err := filepath.EvalSymlinks(parts[2])
			if err != nil {
				continue
			}
			files = append(files, realPath)
		}
	}
	return files
}

// Check if a file is a device file
func isDeviceFile(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	// Check if it's a device file
	if info.Mode()&os.ModeDevice != 0 {
		return true
	}
	return false
}

// Check if a process is a system process
func isSystemProcess(pid, processName, cmdlinePath, statusPath string) bool {
	// Filter by PID range
	pidInt, err := strconv.Atoi(pid)
	if err != nil || pidInt <= 1000 {
		return true
	}

	// Filter by process name
	if isSystemProcessByName(processName) {
		return true
	}

	// Filter by command line
	if isSystemProcessByCmdline(cmdlinePath) {
		return true
	}

	// Filter by process status
	if isSystemProcessByStatus(statusPath) {
		return true
	}

	return false
}

// Check if a process is a system process by name
func isSystemProcessByName(processName string) bool {
	systemProcesses := []string{
		"systemd", "kthreadd", "kworker", "rcu_sched", "migration",
		"watchdog", "kswapd", "ksoftirqd", "init", "bioset",
		"dbus-daemon", "dbus", "cupsd", "snapd", "sshd",
		"gnome-shell", "gnome-session-binary", "gnome-keyring-daemon",
		"gsd-a11y-settings", "gsd-color", "gsd-datetime", "gsd-housekeeping",
		"gsd-keyboard", "gsd-media-keys", "gsd-power", "gsd-printer",
		"gsd-print-notifications", "gsd-rfkill", "gsd-screensaver-proxy",
		"gsd-sharing", "gsd-smartcard", "gsd-sound", "gsd-wacom",
		"gvfs-afc-volume-monitor", "gvfsd", "gvfsd-fuse", "gvfs-goa-volume-monitor",
		"gvfs-gphoto2-volume-monitor", "gvfs-mtp-volume-monitor",
		"gvfs-udisks2-volume-monitor", "plasmashell", "kwin_x11",
		"xdg-document-portal", "xdg-permission-store",
	}
	for _, name := range systemProcesses {
		if strings.Contains(processName, name) {
			return true
		}
	}
	return false
}

// Check if a process is a system process by command line
func isSystemProcessByCmdline(cmdlinePath string) bool {
	cmdlineData, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return true
	}

	// If cmdline is empty, it might be a kernel thread or system process
	if len(cmdlineData) == 0 {
		return true
	}

	// Check if it contains system service keywords
	systemKeywords := []string{
		"/usr/lib/systemd", "/lib/systemd", "init", "kthreadd", "kworker",
	}
	cmdline := string(cmdlineData)
	for _, keyword := range systemKeywords {
		if strings.Contains(cmdline, keyword) {
			return true
		}
	}

	return false
}

// Check if a process is a system process by status
func isSystemProcessByStatus(statusPath string) bool {
	data, err := os.ReadFile(statusPath)
	if err != nil {
		return true
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Check if the process name is a kernel thread
		if strings.HasPrefix(line, "Name:") {
			fields := strings.Fields(line)
			if len(fields) > 1 && strings.HasPrefix(fields[1], "[") && strings.HasSuffix(fields[1], "]") {
				return true
			}
		}
	}
	return false
}

// Check if a string is numeric
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// Check if a file should be skipped
func shouldSkipFile(filePath string) bool {
	// Skip files in /proc and /sys filesystems
	if strings.HasPrefix(filePath, "/proc/") || strings.HasPrefix(filePath, "/sys/") {
		return true
	}

	// Skip anon_inode files
	if strings.HasPrefix(filePath, "anon_inode:") {
		return true
	}

	// Check file type
	info, err := os.Stat(filePath)
	if err != nil {
		return true
	}

	// Skip FIFO (named pipe) files
	if info.Mode()&os.ModeNamedPipe != 0 {
		return true
	}

	// Skip device files
	if info.Mode()&os.ModeDevice != 0 {
		return true
	}

	// Skip socket files
	if info.Mode()&os.ModeSocket != 0 {
		return true
	}

	return false
}

// Check if a file is relevant for processing
func isRelevantFile(filePath string) bool {
	// Use encapsulated skip check method
	if shouldSkipFile(filePath) {
		return false
	}

	allowedExtensions := []string{".jar", ".war", ".sh", ".py", ".pl", ".so", ".conf", ".ini"}
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(filePath, ext) {
			return true
		}
	}

	return isELFFile(filePath)
}

// Check if a file is an ELF file
func isELFFile(filePath string) bool {
	// Use encapsulated skip check method
	if shouldSkipFile(filePath) {
		return false
	}

	realPath, err := filepath.EvalSymlinks(filePath)
	if err != nil {
		return false
	}

	file, err := os.Open(realPath)
	if err != nil {
		return false
	}
	defer func() { _ = file.Close() }()

	magic := make([]byte, 4)
	_, err = file.Read(magic)
	if err != nil {
		return false
	}

	// Check for ELF file magic number
	return magic[0] == 0x7F && magic[1] == 'E' && magic[2] == 'L' && magic[3] == 'F'
}

// Stream-based hash calculation for MD5, SHA1, and SHA256
func computeHashesFromReader(reader io.Reader) (HashResult, error) {
	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()

	multiWriter := io.MultiWriter(md5Hash, sha1Hash, sha256Hash)

	_, err := io.Copy(multiWriter, reader)
	if err != nil {
		return HashResult{}, err
	}

	var result HashResult
	copy(result.MD5[:], md5Hash.Sum(nil))
	copy(result.SHA1[:], sha1Hash.Sum(nil))
	copy(result.SHA256[:], sha256Hash.Sum(nil))

	return result, nil
}

// Extract JAR/WAR files and return nested JAR files and their hashes
func extractJarFilesWithHashes(ctx context.Context, filePath string) ([]string, []HashResult) {
	LOG := logctx.Use(ctx).Sugar()
	var jarFiles []string
	var hashes []HashResult

	// Open ZIP file
	r, err := zip.OpenReader(filePath)
	if err != nil {
		LOG.Warnf("Failed to open JAR/WAR file %s: %v", filePath, err)
		return jarFiles, hashes
	}
	defer func() { _ = r.Close() }()

	// Iterate through ZIP contents, only process first-level JAR/WAR files
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".jar") || strings.HasSuffix(f.Name, ".war") {
			// Add JAR/WAR filename to the list
			jarFiles = append(jarFiles, f.Name)

			// Open nested file
			rc, err := f.Open()
			if err != nil {
				LOG.Warnf("Failed to open nested file %s: %v", f.Name, err)
				continue
			}

			// Stream hash calculation directly from reader without loading entire content to memory
			hash, err := computeHashesFromReader(rc)
			if err != nil {
				LOG.Warnf("Failed to calculate hash for nested file %s: %v", f.Name, err)
			} else {
				// Add hash values to results
				hashes = append(hashes, hash)
			}
			_ = rc.Close() // Ignore errors when closing reader
		}
	}

	return jarFiles, hashes
}

package inspector

import (
	"bufio"
	"github.com/google/uuid"
	"io"
	"io/fs"
	"murphysec-cli-simple/logger"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var pyImportPattern = regexp.MustCompile("^import +([A-Za-z_-][A-Za-z_0-9-]*)(?:\\.[A-Za-z_-][A-Za-z_0-9-]*)*|^from +([A-Za-z_-][A-Za-z_0-9-]*)(?:\\.[A-Za-z_-][A-Za-z_0-9-]*)* +import")
var PythonUUID = uuid.Must(uuid.Parse("fab5210e-2a75-4f89-875f-f05544264f50"))

func ScanPythonImport(dir string) map[string]string {
	componentMap := map[string]string{}
	requirementsFiles := map[string]struct{}{}
	ignoreSet := map[string]struct{}{}
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			ignoreSet[d.Name()] = struct{}{}
			return nil
		}
		if (filepath.Ext(path) == ".txt" || filepath.Ext(path) == "") && strings.HasPrefix(d.Name(), "requirements") {
			requirementsFiles[path] = struct{}{}
			return nil
		}
		if filepath.Ext(path) != ".py" {
			return nil
		}
		f, e := os.Open(path)
		if e != nil {
			return e
		}
		defer f.Close()
		scanner := bufio.NewScanner(io.LimitReader(f, 4*1024*1024))
		scanner.Split(bufio.ScanLines)
		scanner.Buffer(make([]byte, 16*1024), 16*1024)
		for scanner.Scan() {
			if scanner.Err() != nil {
				return nil
			}
			t := strings.TrimSpace(scanner.Text())
			m := pyImportPattern.FindStringSubmatch(t)
			if m == nil {
				continue
			}
			pkg := m[1]
			if pkg == "" {
				pkg = m[2]
			}
			if pyPkgBlackList[pkg] {
				continue
			}
			componentMap[pkg] = ""
		}
		return nil
	})
	for fp := range requirementsFiles {
		for k, v := range parsePythonRequirements(fp) {
			componentMap[k] = v
		}
	}
	for s := range ignoreSet {
		delete(componentMap, s)
	}
	return componentMap
}

var pyRequirementsPattern = regexp.MustCompile("^([A-Za-z0-9_-]+)==([^= \\n\\r]+)$")

func parsePythonRequirements(p string) map[string]string {
	rs := map[string]string{}
	f, e := os.Open(p)
	if e != nil {
		logger.Warn.Println("Open file failed.", e.Error(), p)
		return nil
	}
	defer f.Close()
	scanner := bufio.NewScanner(io.LimitReader(f, 4*1024*1024))
	for scanner.Scan() {
		if scanner.Err() != nil {
			logger.Warn.Println("read file failed.", e.Error(), p)
			return nil
		}
		t := strings.TrimSpace(scanner.Text())
		m := pyRequirementsPattern.FindStringSubmatch(t)
		if m == nil {
			continue
		}
		rs[m[1]] = m[2]
	}
	return rs
}

var pyPkgBlackList = map[string]bool{
	"string":          true,
	"re":              true,
	"difflib":         true,
	"textwrap":        true,
	"unicodedata":     true,
	"stringprep":      true,
	"readline":        true,
	"rlcompleter":     true,
	"struct":          true,
	"codecs":          true,
	"datetime":        true,
	"zoneinfo":        true,
	"calendar":        true,
	"collections":     true,
	"heapq":           true,
	"bisect":          true,
	"array":           true,
	"weakref":         true,
	"types":           true,
	"copy":            true,
	"pprint":          true,
	"reprlib":         true,
	"enum":            true,
	"graphlib":        true,
	"numbers":         true,
	"math":            true,
	"cmath":           true,
	"decimal":         true,
	"fractions":       true,
	"random":          true,
	"statistics":      true,
	"itertools":       true,
	"functools":       true,
	"operator":        true,
	"pathlib":         true,
	"fileinput":       true,
	"stat":            true,
	"filecmp":         true,
	"tempfile":        true,
	"glob":            true,
	"fnmatch":         true,
	"linecache":       true,
	"shutil":          true,
	"pickle":          true,
	"copyreg":         true,
	"shelve":          true,
	"marshal":         true,
	"dbm":             true,
	"sqlite3":         true,
	"zlib":            true,
	"gzip":            true,
	"bz2":             true,
	"lzma":            true,
	"zipfile":         true,
	"tarfile":         true,
	"csv":             true,
	"configparser":    true,
	"netrc":           true,
	"xdrlib":          true,
	"plistlib":        true,
	"hashlib":         true,
	"hmac":            true,
	"secrets":         true,
	"os":              true,
	"io":              true,
	"time":            true,
	"argparse":        true,
	"getopt":          true,
	"logging":         true,
	"getpass":         true,
	"curses":          true,
	"platform":        true,
	"errno":           true,
	"ctypes":          true,
	"threading":       true,
	"multiprocessing": true,
	"concurrent":      true,
	"subprocess":      true,
	"sched":           true,
	"queue":           true,
	"contextvars":     true,
	"_thread":         true,
	"asyncio":         true,
	"socket":          true,
	"ssl":             true,
	"select":          true,
	"selectors":       true,
	"asyncore":        true,
	"asynchat":        true,
	"signal":          true,
	"mmap":            true,
	"email":           true,
	"json":            true,
	"mailcap":         true,
	"mailbox":         true,
	"mimetypes":       true,
	"base64":          true,
	"binhex":          true,
	"binascii":        true,
	"quopri":          true,
	"uu":              true,
	"html":            true,
	"xml":             true,
	"webbrowser":      true,
	"cgi":             true,
	"cgitb":           true,
	"wsgiref":         true,
	"urllib":          true,
	"http":            true,
	"ftplib":          true,
	"poplib":          true,
	"imaplib":         true,
	"nntplib":         true,
	"smtplib":         true,
	"smtpd":           true,
	"telnetlib":       true,
	"uuid":            true,
	"socketserver":    true,
	"xmlrpc":          true,
	"ipaddress":       true,
	"audioop":         true,
	"aifc":            true,
	"sunau":           true,
	"wave":            true,
	"chunk":           true,
	"colorsys":        true,
	"imghdr":          true,
	"sndhdr":          true,
	"ossaudiodev":     true,
	"gettext":         true,
	"locale":          true,
	"turtle":          true,
	"cmd":             true,
	"shlex":           true,
	"tkinter":         true,
	"typing":          true,
	"pydoc":           true,
	"doctest":         true,
	"unittest":        true,
	"2to3":            true,
	"test":            true,
	"bdb":             true,
	"faulthandler":    true,
	"pdb":             true,
	"timeit":          true,
	"trace":           true,
	"tracemalloc":     true,
	"distutils":       true,
	"ensurepip":       true,
	"venv":            true,
	"zipapp":          true,
	"sys":             true,
	"sysconfig":       true,
	"builtins":        true,
	"__main__":        true,
	"warnings":        true,
	"dataclasses":     true,
	"contextlib":      true,
	"abc":             true,
	"atexit":          true,
	"traceback":       true,
	"__future__":      true,
	"gc":              true,
	"inspect":         true,
	"site":            true,
	"code":            true,
	"codeop":          true,
	"zipimport":       true,
	"pkgutil":         true,
	"modulefinder":    true,
	"runpy":           true,
	"importlib":       true,
	"ast":             true,
	"symtable":        true,
	"token":           true,
	"keyword":         true,
	"tokenize":        true,
	"tabnanny":        true,
	"pyclbr":          true,
	"py_compile":      true,
	"compileall":      true,
	"dis":             true,
	"pickletools":     true,
	"msilib":          true,
	"msvcrt":          true,
	"winreg":          true,
	"winsound":        true,
	"posix":           true,
	"pwd":             true,
	"spwd":            true,
	"grp":             true,
	"crypt":           true,
	"termios":         true,
	"tty":             true,
	"pty":             true,
	"fcntl":           true,
	"pipes":           true,
	"resource":        true,
	"nis":             true,
	"optparse":        true,
	"imp":             true,
}

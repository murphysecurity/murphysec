package chunkupload

import (
	"encoding/binary"
	"io"
	"io/fs"
	"os"
	"strings"
)

type FilterVote int

const (
	_ FilterVote = iota
	FilterAdd
	FilterSkip
	FilterSkipDir
)

type Filter func(path string, entry fs.DirEntry) (FilterVote, error)

var DiscardDot Filter = func(path string, entry fs.DirEntry) (FilterVote, error) {
	var name = entry.Name()
	if strings.HasPrefix(name, ".") || name == "node_modules" {
		if entry.IsDir() {
			return FilterSkipDir, nil
		} else {
			return FilterSkip, nil
		}
	}
	return FilterAdd, nil
}

var BinaryOnly Filter = func(path string, entry fs.DirEntry) (FilterVote, error) {
	if entry.IsDir() {
		return FilterAdd, nil
	}
	eFormat, _ := detectFileFormat(path)
	if eFormat != "" {
		return FilterAdd, nil
	}
	return FilterSkip, nil
}

var uploadAll Filter = func(path string, entry fs.DirEntry) (FilterVote, error) { return FilterAdd, nil }

func detectFileFormat(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	magic := make([]byte, 4)
	_, err = io.ReadFull(file, magic)
	if err != nil {
		return "", err
	}

	// Check ELF (0x7f 'E' 'L' 'F')
	if magic[0] == 0x7f && magic[1] == 'E' && magic[2] == 'L' && magic[3] == 'F' {
		return "ELF", nil
	}

	// Check Mach-O (32-bit or 64-bit, big or little endian)
	magicNum := binary.BigEndian.Uint32(magic)
	switch magicNum {
	case 0xfeedface, 0xcefaedfe, // 32-bit Mach-O
		0xfeedfacf, 0xcffaedfe, // 64-bit Mach-O
		0xcafebabe, 0xbebafeca: // FAT binary
		return "Mach-O", nil
	}

	// Check PE (first two bytes 'M' 'Z')
	if magic[0] == 'M' && magic[1] == 'Z' {
		return "PE", nil
	}

	return "", nil
}

package utils

import (
	"crypto/sha256"
	"io"
)

// ComputeSHA256AndLF reads from reader until EOF and returns:
// 1. the SHA256 of the original bytes
// 2. the SHA256 after normalizing CRLF to LF
func ComputeSHA256AndLF(reader io.Reader) (raw [sha256.Size]byte, lf [sha256.Size]byte, err error) {
	rawHash := sha256.New()
	lfHash := sha256.New()
	lfWriter := Dos2UnixWriter(lfHash)

	if _, err = io.Copy(io.MultiWriter(rawHash, lfWriter), reader); err != nil {
		_ = lfWriter.Close()
		return raw, lf, err
	}
	if err = lfWriter.Close(); err != nil {
		return raw, lf, err
	}

	copy(raw[:], rawHash.Sum(nil))
	copy(lf[:], lfHash.Sum(nil))
	return raw, lf, nil
}

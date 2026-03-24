package utils

import (
	"crypto/sha256"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestComputeSHA256AndLF(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		chunkSizes []int
		wantLF     string
	}{
		{
			name:   "empty",
			input:  "",
			wantLF: "",
		},
		{
			name:   "plain text",
			input:  "alpha\nbeta\ngamma",
			wantLF: "alpha\nbeta\ngamma",
		},
		{
			name:       "crlf normalized across chunks",
			input:      "a\r\nb\r\nc\r\n",
			chunkSizes: []int{2, 1, 2, 1, 1},
			wantLF:     "a\nb\nc\n",
		},
		{
			name:       "lone cr preserved",
			input:      "a\rb\r\nc\r",
			chunkSizes: []int{1, 1, 2, 1, 1},
			wantLF:     "a\rb\nc\r",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reader := io.Reader(strings.NewReader(tt.input))
			if len(tt.chunkSizes) > 0 {
				reader = &chunkedReader{
					data:       []byte(tt.input),
					chunkSizes: tt.chunkSizes,
				}
			}

			raw, lf, err := ComputeSHA256AndLF(reader)
			if err != nil {
				t.Fatalf("ComputeSHA256AndLF failed: %v", err)
			}

			wantRaw := sha256.Sum256([]byte(tt.input))
			wantLF := sha256.Sum256([]byte(tt.wantLF))
			if raw != wantRaw {
				t.Fatalf("raw hash mismatch: got %x want %x", raw, wantRaw)
			}
			if lf != wantLF {
				t.Fatalf("LF hash mismatch: got %x want %x", lf, wantLF)
			}
		})
	}
}

func TestComputeSHA256AndLF_ReadError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	_, _, err := ComputeSHA256AndLF(errorReader{err: wantErr})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

type chunkedReader struct {
	data       []byte
	chunkSizes []int
	offset     int
	index      int
}

func (r *chunkedReader) Read(p []byte) (int, error) {
	if r.offset >= len(r.data) {
		return 0, io.EOF
	}

	size := len(p)
	if r.index < len(r.chunkSizes) && r.chunkSizes[r.index] < size {
		size = r.chunkSizes[r.index]
	}
	remaining := len(r.data) - r.offset
	if size > remaining {
		size = remaining
	}

	copy(p, r.data[r.offset:r.offset+size])
	r.offset += size
	r.index++
	return size, nil
}

type errorReader struct {
	err error
}

func (r errorReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

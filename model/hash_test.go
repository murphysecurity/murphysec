package model

import (
	"encoding/json"
	"testing"
)

func TestMD5Hash_JSON(t *testing.T) {
	const jsonInput = `"5eb63bbbe01eeed093cb22bb8f5acdc3"` // "hello world" 的 MD5

	var h MD5Hash
	if err := json.Unmarshal([]byte(jsonInput), &h); err != nil {
		t.Fatalf("Unmarshal MD5Hash failed: %v", err)
	}

	out, err := json.Marshal(h)
	if err != nil {
		t.Fatalf("Marshal MD5Hash failed: %v", err)
	}

	if string(out) != jsonInput {
		t.Errorf("MD5Hash JSON round-trip mismatch: got %s, want %s", string(out), jsonInput)
	}
}

func TestSHA256Hash_JSON(t *testing.T) {
	const jsonInput = `"b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"` // "hello world" 的 SHA256

	var h SHA256Hash
	if err := json.Unmarshal([]byte(jsonInput), &h); err != nil {
		t.Fatalf("Unmarshal SHA256Hash failed: %v", err)
	}

	out, err := json.Marshal(h)
	if err != nil {
		t.Fatalf("Marshal SHA256Hash failed: %v", err)
	}

	if string(out) != jsonInput {
		t.Errorf("SHA256Hash JSON round-trip mismatch: got %s, want %s", string(out), jsonInput)
	}
}

func TestMD5Hash_InvalidHex(t *testing.T) {
	const invalidJSON = `"not-a-valid-hex"`

	var h MD5Hash
	if err := json.Unmarshal([]byte(invalidJSON), &h); err == nil {
		t.Fatal("expected error for invalid hex MD5 string, got nil")
	}
}

func TestMD5Hash_InvalidLength(t *testing.T) {
	const tooShort = `"deadbeef"`                          // 明显长度不足
	const tooLong = `"5eb63bbbe01eeed093cb22bb8f5acdc300"` // 多出一位

	var h MD5Hash
	if err := json.Unmarshal([]byte(tooShort), &h); err == nil {
		t.Error("expected error for short MD5 hash, got nil")
	}
	if err := json.Unmarshal([]byte(tooLong), &h); err == nil {
		t.Error("expected error for long MD5 hash, got nil")
	}
}

func TestSHA256Hash_InvalidHex(t *testing.T) {
	const invalidJSON = `"not-a-valid-hex"`

	var h SHA256Hash
	if err := json.Unmarshal([]byte(invalidJSON), &h); err == nil {
		t.Fatal("expected error for invalid hex SHA256 string, got nil")
	}
}

func TestSHA256Hash_InvalidLength(t *testing.T) {
	const tooShort = `"deadbeef"`
	const tooLong = `"b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde900"`

	var h SHA256Hash
	if err := json.Unmarshal([]byte(tooShort), &h); err == nil {
		t.Error("expected error for short SHA256 hash, got nil")
	}
	if err := json.Unmarshal([]byte(tooLong), &h); err == nil {
		t.Error("expected error for long SHA256 hash, got nil")
	}
}
func TestSHA1Hash_JSON(t *testing.T) {
	const jsonInput = `"2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"` // "hello world" 的 SHA1

	var h SHA1Hash
	if err := json.Unmarshal([]byte(jsonInput), &h); err != nil {
		t.Fatalf("Unmarshal SHA1Hash failed: %v", err)
	}

	out, err := json.Marshal(h)
	if err != nil {
		t.Fatalf("Marshal SHA1Hash failed: %v", err)
	}

	if string(out) != jsonInput {
		t.Errorf("SHA1Hash JSON round-trip mismatch: got %s, want %s", string(out), jsonInput)
	}
}

func TestSHA1Hash_InvalidHex(t *testing.T) {
	const invalidJSON = `"not-a-valid-hex"`

	var h SHA1Hash
	if err := json.Unmarshal([]byte(invalidJSON), &h); err == nil {
		t.Fatal("expected error for invalid hex SHA1 string, got nil")
	}
}

func TestSHA1Hash_InvalidLength(t *testing.T) {
	const tooShort = `"deadbeef"`
	const tooLong = `"2aae6c35c94fcfb415dbe95f408b9ce91ee846ed00"`

	var h SHA1Hash
	if err := json.Unmarshal([]byte(tooShort), &h); err == nil {
		t.Error("expected error for short SHA1 hash, got nil")
	}
	if err := json.Unmarshal([]byte(tooLong), &h); err == nil {
		t.Error("expected error for long SHA1 hash, got nil")
	}
}

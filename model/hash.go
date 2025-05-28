package model

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// MD5Hash 表示一个 MD5 哈希值（16字节）
type MD5Hash [md5.Size]byte

// SHA256Hash 表示一个 SHA256 哈希值（32字节）
type SHA256Hash [sha256.Size]byte

// String 将 MD5Hash 转换为十六进制字符串
func (h MD5Hash) String() string {
	return hex.EncodeToString(h[:])
}

// MarshalJSON 实现 json.Marshaler 接口
func (h MD5Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (h *MD5Hash) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}

	// 从十六进制字符串解析
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}

	if len(b) != md5.Size {
		return fmt.Errorf("invalid MD5 hash length: got %d, want %d", len(b), md5.Size)
	}

	copy(h[:], b)
	return nil
}

// String 将 SHA256Hash 转换为十六进制字符串
func (h SHA256Hash) String() string {
	return hex.EncodeToString(h[:])
}

// MarshalJSON 实现 json.Marshaler 接口
func (h SHA256Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (h *SHA256Hash) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}

	// 从十六进制字符串解析
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}

	if len(b) != sha256.Size {
		return fmt.Errorf("invalid SHA256 hash length: got %d, want %d", len(b), sha256.Size)
	}

	copy(h[:], b)
	return nil
}

// ParseMD5Hash 从十六进制字符串解析 MD5Hash
func ParseMD5Hash(hexStr string) (MD5Hash, error) {
	var h MD5Hash
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return h, err
	}

	if len(b) != md5.Size {
		return h, fmt.Errorf("invalid MD5 hash length: got %d, want %d", len(b), md5.Size)
	}

	copy(h[:], b)
	return h, nil
}

// ParseSHA256Hash 从十六进制字符串解析 SHA256Hash
func ParseSHA256Hash(hexStr string) (SHA256Hash, error) {
	var h SHA256Hash
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return h, err
	}

	if len(b) != sha256.Size {
		return h, fmt.Errorf("invalid SHA256 hash length: got %d, want %d", len(b), sha256.Size)
	}

	copy(h[:], b)
	return h, nil
}

// ComputeMD5Hash 计算数据的 MD5 哈希值
func ComputeMD5Hash(data []byte) MD5Hash {
	return md5.Sum(data)
}

// ComputeSHA256Hash 计算数据的 SHA256 哈希值
func ComputeSHA256Hash(data []byte) SHA256Hash {
	return sha256.Sum256(data)
}

// SHA1Hash 表示一个 SHA1 哈希值（20字节）
type SHA1Hash [sha1.Size]byte

// String 将 SHA1Hash 转换为十六进制字符串
func (h SHA1Hash) String() string {
	return hex.EncodeToString(h[:])
}

// MarshalJSON 实现 json.Marshaler 接口
func (h SHA1Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (h *SHA1Hash) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}

	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}

	if len(b) != sha1.Size {
		return fmt.Errorf("invalid SHA1 hash length: got %d, want %d", len(b), sha1.Size)
	}

	copy(h[:], b)
	return nil
}

// ParseSHA1Hash 从十六进制字符串解析 SHA1Hash
func ParseSHA1Hash(hexStr string) (SHA1Hash, error) {
	var h SHA1Hash
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return h, err
	}
	if len(b) != sha1.Size {
		return h, fmt.Errorf("invalid SHA1 hash length: got %d, want %d", len(b), sha1.Size)
	}
	copy(h[:], b)
	return h, nil
}

// ComputeSHA1Hash 计算数据的 SHA1 哈希值
func ComputeSHA1Hash(data []byte) SHA1Hash {
	return sha1.Sum(data)
}

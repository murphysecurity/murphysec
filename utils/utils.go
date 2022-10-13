package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/murphysecurity/murphysec/utils/must"
	"go.uber.org/zap"
	"io"
	"net"
	"os"
	"strings"
)

func InStringSlice(slice []string, s string) bool {
	for _, it := range slice {
		if it == s {
			return true
		}
	}
	return false
}

func JoinStringAny(sep string, a ...fmt.Stringer) string {
	var s []string
	for i := range a {
		s = append(s, a[i].String())
	}
	return strings.Join(s, sep)
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func DistinctStringSlice(s []string) []string {
	m := map[string]struct{}{}
	var rs []string
	for _, it := range s {
		if _, ok := m[it]; !ok {
			rs = append(rs, it)
			m[it] = struct{}{}
		}
	}
	return rs
}

func base64UrlEncode(s string) string {
	b := new(bytes.Buffer)
	w := base64.NewEncoder(base64.URLEncoding, b)
	must.A(w.Write([]byte(s)))
	must.Close(w)
	return b.String()
}

func ReadFileLimited(p string, maxRead int64) ([]byte, error) {
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	return io.ReadAll(io.LimitReader(f, maxRead))
}

func CloseLogErrZap(closer io.Closer, logger *zap.Logger) {
	if e := closer.Close(); e != nil {
		logger.Error("Close error", zap.Error(e))
	}
}

func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func GetOutBoundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return ""
	}
	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok{
		return ""
	}
	return localAddr.IP.String()
}

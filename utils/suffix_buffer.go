package utils

type SuffixBuffer struct {
	b         []byte
	pos       int
	truncated bool
}

func NewSuffixBuffer(size int) *SuffixBuffer {
	return &SuffixBuffer{
		b:         make([]byte, size),
		pos:       0,
		truncated: false,
	}
}

func (s *SuffixBuffer) Truncated() bool {
	return s.truncated
}

func (s *SuffixBuffer) Write(data []byte) (n int, e error) {
	p := 0
	if len(data) > len(s.b) {
		p = len(data) - len(s.b)
		s.truncated = true
		s.pos = 0
		copy(s.b, data[p:])
		return
	}
	w := MinInt(len(s.b)-s.pos, len(data))
	copy(s.b[s.pos:], data[p:p+w])
	p += w
	if p == len(data) {
		s.pos += w
		return
	}
	s.truncated = true
	s.pos = len(data) - p
	copy(s.b, data[p:])
	return len(data), nil
}

func (s *SuffixBuffer) Bytes() []byte {
	if !s.truncated {
		r := make([]byte, s.pos)
		copy(r, s.b[:s.pos])
		return r
	}
	r := make([]byte, len(s.b))
	copy(r, s.b[s.pos:])
	copy(r[len(r)-(s.pos):], s.b[:s.pos])
	return r
}

func (s *SuffixBuffer) String() string {
	return string(s.Bytes())
}

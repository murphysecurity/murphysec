package suffixbuf

import "github.com/murphysecurity/murphysec/utils"

func NewSize(capacity int) *Buf {
	return &Buf{
		data:     make([]byte, capacity),
		pos:      0,
		overflow: false,
	}
}

type Buf struct {
	data     []byte
	pos      int
	overflow bool
}

func (r *Buf) Bytes() (rs []byte) {
	if r.overflow {
		rs = make([]byte, len(r.data))
		copy(rs, r.data[r.pos:])
		copy(rs[len(r.data)-r.pos:], r.data[:r.pos])
	} else {
		rs = make([]byte, r.pos)
		copy(rs, r.data[:r.pos])
	}
	return
}

func (r *Buf) Write(data []byte) (int, error) {
	r.write(data)
	return len(data), nil
}

func (r *Buf) write(input []byte) {
	srcPos := 0
	if len(input) > len(r.data) {
		srcPos = len(input) - len(r.data)
	}
	// part 1
	k := utils.MinInt(len(r.data)-r.pos, len(input)-srcPos)
	copy(r.data[r.pos:r.pos+k], input[srcPos:srcPos+k])
	r.pos += k
	srcPos += k
	if srcPos == len(input) {
		return
	}
	// part 2
	r.pos = 0
	r.overflow = true
	copy(r.data[:len(input)-srcPos], input[srcPos:])
	r.pos += len(input) - srcPos
}

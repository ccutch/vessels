package memory

import (
	"bytes"
)

type Buffer struct{ bytes.Buffer }

func (v *Buffer) Close() (err error) {
	v.Reset()
	return
}

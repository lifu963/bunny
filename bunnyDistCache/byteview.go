package bunnyDistCache

type ByteView struct {
	B []byte
}

func cloneBytes(bytes []byte) []byte {
	copyBytes := make([]byte, len(bytes))
	copy(copyBytes, bytes)
	return copyBytes
}

func (v ByteView) Len() int {
	return len(v.B)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.B)
}

func (v ByteView) String() string {
	return string(v.B)
}

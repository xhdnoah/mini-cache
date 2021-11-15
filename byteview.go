package cache

// 抽象一个只读数据结构表示缓存值
// A ByteView holds an immutable view of bytes.
type ByteView struct {
	// byte 能支持任意数据类型的存储，如字符串、图片
	b []byte
}

// Len returns the view's length
func (v ByteView) Len() int { return len(v.b) }

// 返回拷贝，防止缓存值被外部修改
// ByteSlice returns a copy of the data as a byte slice.
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String returns the data as a string, making a copy if necessary.
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

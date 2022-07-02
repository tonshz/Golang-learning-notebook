package convert

import "strconv"

type StrTo string

func (s StrTo) String() string {
	return string(s)
}
func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

// 强制转换为 int，不输出错误信息
func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

// uint32: 32位无符号整数
func (s StrTo) UInt32() (uint32, error) {
	v, err := strconv.Atoi(s.String())
	return uint32(v), err
}

// 强制转换为 uint32, 不输出错误信息
func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}

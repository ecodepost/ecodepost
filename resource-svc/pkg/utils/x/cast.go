package x

type Enum interface {
	~int | ~int16 | ~int32 | ~int64
}

// E2I32 把Enum类型转换为int32
func E2I32[T Enum](enum T) int32 {
	return int32(enum)
}

// Es2I32s 把Enum列表转换为[]int32
func Es2I32s[T Enum](enums []T) []int32 {
	res := make([]int32, 0, len(enums))
	for _, v := range enums {
		res = append(res, int32(v))
	}
	return res
}

// E2I64 把Enum类型转换为int64j
func E2I64[T Enum](enum T) int64 {
	return int64(enum)
}

// Es2I64s 把Enum列表转换为[]int64
func Es2I64s[T Enum](enums []T) []int64 {
	res := make([]int64, 0, len(enums))
	for _, v := range enums {
		res = append(res, int64(v))
	}
	return res
}

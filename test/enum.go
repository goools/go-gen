package test

//go:generate ./../go-gen -v gen enum EnumA EnumB
type EnumA int
type EnumB int

const (
	// aaaa
	EnumANone EnumA = 0 // vvvv
)

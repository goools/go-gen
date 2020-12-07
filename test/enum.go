package test

//go:generate go run ../main.go -v gen enum EnumA EnumB
type EnumA int
type EnumB int

const (
	// comment None
	EnumANone EnumA = 0
)

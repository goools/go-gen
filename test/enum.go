package test

//go:generate go run ../main.go -v gen enum EnumA EnumB
type EnumA int
type EnumB int

const (
	// comment EnumANone
	// 
	EnumANone EnumA = iota
	EnumANone2
)

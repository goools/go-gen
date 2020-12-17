package enum

//go:generate go-gen -v gen enum EnumA EnumB
type EnumA int
type EnumB int

const (
	// comment EnumANone
	//
	EnumANone EnumA = iota
	EnumANone2
)

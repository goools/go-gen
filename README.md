# go-gen
> generate golang source code by go generate command

## install go-gen
``` shell script
go install github.com/goools/go-gen
```

## Enum
### [example](./example/enum)
``` go
package enum

//go:generate go run go-gen -v gen enum EnumA EnumB
type EnumA int
type EnumB int

const (
	// comment EnumANone
	//
	EnumANone EnumA = iota
	EnumANone2
)
```

#### generate file content
##### enum_a_enum_generate.go
```go
// Code generated by go-gen DO NOT EDIT. visit https://github.com/goools/go-gen

package enum

import "fmt"

func (e *EnumA) String() string {
	switch *e {
	case EnumANone:
		return "None"
	case EnumANone2:
		return "None2"
	}
	return "UNKNOWN"
}

func (e *EnumA) MarshalText() ([]byte, error) {
	switch *e {
	case EnumANone:
		return []byte("None"), nil
	case EnumANone2:
		return []byte("None2"), nil
	}
	return []byte("UNKNOWN"), fmt.Errorf("not found EnumA, value: %v", *e)
}

func (e *EnumA) UnmarshalText(enumBytes []byte) error {
	enumString := string(enumBytes)
	switch enumString {
	case "None":
		*e = EnumANone
		return nil
	case "None2":
		*e = EnumANone2
		return nil
	}
	return fmt.Errorf("not found EnumA, value: %v", enumString)
}

func (e *EnumA) Comment() string {
	switch *e {
	case EnumANone:
		return "comment EnumANone"
	case EnumANone2:
		return "None2"
	default:
		return "UNKNOWN"
	}
}
```

##### enum_b_enum_generate.go
```go
// Code generated by go-gen DO NOT EDIT. visit https://github.com/goools/go-gen

package enum

import "fmt"

func (e *EnumB) String() string {
	switch *e {
	}
	return "UNKNOWN"
}

func (e *EnumB) MarshalText() ([]byte, error) {
	switch *e {
	}
	return []byte("UNKNOWN"), fmt.Errorf("not found EnumB, value: %v", *e)
}

func (e *EnumB) UnmarshalText(enumBytes []byte) error {
	enumString := string(enumBytes)
	switch enumString {
	}
	return fmt.Errorf("not found EnumB, value: %v", enumString)
}

func (e *EnumB) Comment() string {
	switch *e {
	default:
		return "UNKNOWN"
	}
}
```
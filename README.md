# go-gen
> generate golang source code by go generate command

## install go-gen
``` shell script
go install github.com/goools/go-gen
```

## Enum
### [example](./example/enum/enum.go)
``` go
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
```

## Sync Map
### [example](./example/syncmap/gen.go)
```go
package syncmap

//go:generate go-gen -v gen syncmap Pill<int,int> Pill2<string,int> Pill3<string,encoding/json.Number>

```

## Sync Pool
### [example](./example/syncpool/gen.go)
```go
package syncpool

//go:generate go-gen -v gen syncpool Pill<int> Pill2<string> Pill3<encoding/json.Number> Pill4<Pill>

```
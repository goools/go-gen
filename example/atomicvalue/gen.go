package atomicvalue

//go:generate go run ../../main.go -v gen atomicvalue Pill<int> Pill2<string> Pill3<encoding/json.Number> Pill4<Pill> Pill5<*int> Pill6<*encoding/json.Number>

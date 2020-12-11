package test

import (
	"encoding/json"
	"testing"
)

type Demo struct {
	A EnumA `json:"a"`
}

func TestEnumA_MarshalText(t *testing.T) {
	demo := &Demo{A: EnumANone2}
	bds, err := json.Marshal(demo)
	if err != nil {
		t.Fatal("marshal have an err:", err)
	}
	t.Logf("bds: %s", string(bds))
}

func TestEnumA_UnmarshalText(t *testing.T) {
	jsonStr := []byte(`{"a":"None2"}`)
	demo := &Demo{}
	err := json.Unmarshal(jsonStr, demo)
	if err != nil {
		t.Fatal("err: ", err)
	}
	t.Logf("demo: %#v, enum: %s", demo, demo.A.String())
}

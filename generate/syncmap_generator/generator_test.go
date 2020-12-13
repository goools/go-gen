package syncmap_generator

import (
	"strings"
	"testing"
)

func Test_syncMapDefRegexp(t *testing.T) {
	var res []string
	res = syncMapDefRegexp.FindStringSubmatch("Pill<int,int>")
	t.Logf("res: %#v", res)
}

func Test_split(t *testing.T) {
	index := strings.LastIndex("a.b.c", ".")
	t.Logf("index: %#v, %#v, %#v", index, "a.b.c"[:index], "a.b.c"[index+1:])
}

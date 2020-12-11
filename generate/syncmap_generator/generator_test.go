package syncmap_generator

import (
	"testing"
)

func Test_syncMapDefRegexp(t *testing.T) {
	var res []string
	res = syncMapDefRegexp.FindStringSubmatch("Pill<int,int>")
	t.Logf("res: %#v", res)
}

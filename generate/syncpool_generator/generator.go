package syncpool_generator

import (
	"fmt"
	"regexp"

	"github.com/goools/go-gen/generate"
	"github.com/sirupsen/logrus"
)

var (
	syncPoolDefRegexp = regexp.MustCompile(`([A-Za-z_][A-Za-z0-9]*)<([A-Za-z][A-Za-z0-9/.]*)>`)
)

const (
	syncPoolDefRegexpTotal          = 3
	syncPoolDefRegexpNameIndex      = 1
	syncPoolDefRegexpValueTypeIndex = 2
	syncPoolFuncThisName            = "p"
)

type SyncPool struct {
	PkgPath      string
	Name         string
	Value        *generate.Type
	funcTypeName string
	funcThisName string
}

func NewSyncPool(pkgPath, syncPoolDef string) *SyncPool {
	syncPool := &SyncPool{
		PkgPath: pkgPath,
	}
	syncPool.parseDef(syncPoolDef)
	syncPool.funcThisName = syncPoolFuncThisName
	syncPool.funcTypeName = syncPool.Name
	return syncPool
}

func (syncPool *SyncPool) parseDef(syncPoolDef string) {
	res := syncPoolDefRegexp.FindStringSubmatch(syncPoolDef)
	if len(res) != syncPoolDefRegexpTotal {
		panic(fmt.Errorf("cannot parse sync pool def: %s", syncPoolDef))
	}
	syncPool.Name = res[syncPoolDefRegexpNameIndex]
	valueType := res[syncPoolDefRegexpValueTypeIndex]
	syncPool.Value = generate.NewType(valueType)
}

func (syncPool *SyncPool) WriteToFile() {
	logrus.Infof("begin generate sync pool: %s", syncPool.Name)
	syncPoolSnackName := generate.ToSnakeCase(syncPool.Name)
	syncPoolFileName := fmt.Sprintf("%s_syncpool_generate.go", syncPoolSnackName)
	logrus.Infof("file: %s", syncPoolFileName)
}

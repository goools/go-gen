package gen

import (
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/generate/syncmap_generator"
	"github.com/spf13/cobra"
)

var (
	cmdGenSyncMap = &cobra.Command{
		Use:   "syncmap",
		Short: "generate interfaces of sync map",
		Run: func(cmd *cobra.Command, args []string) {
			generate.RunGenerator(syncmap_generator.NewSyncMapGenerator, args)
		},
	}
)

func init() {
	CmdGen.AddCommand(cmdGenSyncMap)
}

package gen

import (
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/generate/syncpool_generator"
	"github.com/spf13/cobra"
)

var (
	cmdGenSyncPool = &cobra.Command{
		Use:   "syncpool",
		Short: "generate interfaces of sync pool",
		Run: func(cmd *cobra.Command, args []string) {
			generate.RunGenerator(syncpool_generator.NewSyncPoolGenerator, args)
		},
	}
)

func init() {
	CmdGen.AddCommand(cmdGenSyncPool)
}

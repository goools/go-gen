package gen

import (
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/generate/deepcopy_generator"
	"github.com/spf13/cobra"
)

var (
	cmdDeepCopyEnum = &cobra.Command{
		Use:   "deepcopy",
		Short: "generate interfaces of deep copy",
		Run: func(cmd *cobra.Command, args []string) {
			generate.RunGenerator(deepcopy_generator.NewDeepCopyGenerator, args)
		},
	}
)

func init() {
	CmdGen.AddCommand(cmdDeepCopyEnum)
}

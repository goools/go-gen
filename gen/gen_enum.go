package gen

import (
	"github.com/goools/go-gen/generate"
	"github.com/spf13/cobra"
)

var (
	cmdGenEnum = &cobra.Command{
		Use:   "enum",
		Short: "generate interfaces of enumeration",
		Run: func(cmd *cobra.Command, args []string) {
			generate.RunGenerator(generate.NewEnumGenerator, args)
		},
	}
)

func init() {
	CmdGen.AddCommand(cmdGenEnum)
}

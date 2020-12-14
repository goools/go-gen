package gen

import (
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/generate/atomicvalue_generator"
	"github.com/spf13/cobra"
)

var (
	cmdGenAtomicValue = &cobra.Command{
		Use:   "atomicvalue",
		Short: "generate interfaces of atomic value",
		Run: func(cmd *cobra.Command, args []string) {
			generate.RunGenerator(atomicvalue_generator.NewAtomicValueGenerator, args)
		},
	}
)

func init() {
	CmdGen.AddCommand(cmdGenAtomicValue)
}

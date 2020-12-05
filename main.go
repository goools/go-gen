package main

import (
	"github.com/goools/go-gen/gen"
	"github.com/goools/go-gen/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	verbose = false
	cmdRoot = &cobra.Command{
		Use:     "",
		Version: version.Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				logrus.SetLevel(logrus.InfoLevel)
			}
		},
	}
)

func init() {
	cmdRoot.PersistentFlags().BoolVarP(&verbose, "verbose", "v", verbose, "")

	cmdRoot.AddCommand(gen.CmdGen)
}

func main() {
	if err := cmdRoot.Execute(); err != nil {
		logrus.Fatalf("cmd root execute have an err: %v", err)
	}
}

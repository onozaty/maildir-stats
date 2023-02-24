package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "maildir-stats",
		Short: "Report maildir statistics",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newUserCmd())
	rootCmd.AddCommand(newAllCmd())
	rootCmd.AddCommand(newUserListCmd())

	for _, c := range rootCmd.Commands() {
		// フラグ以外は受け付けないように
		c.Args = func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("only flags can be specified")
			}
			return nil
		}
		c.Flags().SortFlags = false
		c.InheritedFlags().SortFlags = false
	}

	return rootCmd
}

func Execute() {

	rootCmd := newRootCmd()
	cobra.CheckErr(rootCmd.Execute())
}

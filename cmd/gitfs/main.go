package main

import (
	"fmt"
	"os"

	"github.com/mjpitz/gitfs/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "gitfs",
		Run: func(this *cobra.Command, args []string) {
			// run the start command by default
			cmd.StartCommand.Run(this, args)
		},
	}

	rootCmd.AddCommand(cmd.StartCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

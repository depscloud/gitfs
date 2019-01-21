package main

import (
	"fmt"
	"os"

	"github.com/mjpitz/gitfs/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	logrus.SetOutput(os.Stdout)

	rootCmd := &cobra.Command{
		Use: "gitfs",
	}

	// all commands need config
	rootCmd.PersistentFlags().StringVar(&cmd.ConfigPath, "config",
		cmd.ConfigPath,
		"Specify the configuration path")

	rootCmd.AddCommand(cmd.StartCommand)
	rootCmd.AddCommand(cmd.StopCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

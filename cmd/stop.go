package cmd

import (
	"os"
	"syscall"

	"github.com/mjpitz/gitfs/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

// StopCommand defines the cobra.Command used to stop the system.
var StopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stops the file system server.",
	Run: func(cmd *cobra.Command, args []string) {
		props := ConfigPath
		if len(props) == 0 {
			fail("missing config file")
		}

		cfg, err := config.Load(os.ExpandEnv(props))
		if err != nil {
			fail("failed to parse configuration: %v", err)
		}

		mountpoint := os.ExpandEnv(cfg.Mount)

		logrus.Infof("Unmounting %s", mountpoint)

		_ = syscall.Unmount(mountpoint, unix.MNT_FORCE)
		// ignore the error since it's likely not mounted
	},
}

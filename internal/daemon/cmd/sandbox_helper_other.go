//go:build !linux

package cmd

import (
	"fmt"
	"os"
)

// runLandlockHelper is a no-op on non-Linux platforms.
// The --sandbox-exec flag is only used by the Landlock helper on Linux.
func runLandlockHelper(args []string) {
	fmt.Fprintf(os.Stderr, "watchfired --sandbox-exec is only supported on Linux\n")
	os.Exit(1)
}

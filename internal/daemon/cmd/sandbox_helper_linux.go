//go:build linux

package cmd

import "github.com/watchfire-io/watchfire/internal/daemon/agent"

// runLandlockHelper delegates to the agent package's Landlock helper.
func runLandlockHelper(args []string) {
	agent.RunLandlockHelper(args)
}

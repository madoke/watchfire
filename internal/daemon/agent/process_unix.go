//go:build !windows

package agent

import (
	"syscall"
	"time"
)

// Stop terminates the agent process. Sends SIGTERM, waits 5 seconds, then SIGKILL.
func (p *Process) Stop() {
	if p.cmd.Process == nil {
		return
	}

	// Send SIGTERM
	_ = p.cmd.Process.Signal(syscall.SIGTERM)

	// Wait up to 5 seconds for graceful exit
	select {
	case <-p.done:
		p.Cleanup()
		return
	case <-time.After(5 * time.Second):
	}

	// Force kill
	_ = p.cmd.Process.Kill()
	<-p.done
	p.Cleanup()
}

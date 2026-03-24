//go:build windows

package agent

// Stop terminates the agent process.
// Windows has no SIGTERM equivalent, so we kill the process directly.
func (p *Process) Stop() {
	if p.cmd.Process == nil {
		return
	}

	_ = p.cmd.Process.Kill()
	<-p.done
	p.Cleanup()
}

//go:build windows

package cli

import (
	"context"

	pb "github.com/watchfire-io/watchfire/proto"
)

// watchWindowResize is a no-op on Windows.
// Windows terminals handle resize via the PTY library directly.
func watchWindowResize(ctx context.Context, client pb.AgentServiceClient, projectID string) {
	// Block until context is done so the goroutine doesn't exit immediately
	<-ctx.Done()
}

//go:build !windows

package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"

	pb "github.com/watchfire-io/watchfire/proto"
)

// watchWindowResize handles SIGWINCH signals and sends resize requests.
func watchWindowResize(ctx context.Context, client pb.AgentServiceClient, projectID string) {
	sigwinchCh := make(chan os.Signal, 1)
	signal.Notify(sigwinchCh, syscall.SIGWINCH)
	for range sigwinchCh {
		cols, rows, err := term.GetSize(int(os.Stdin.Fd()))
		if err == nil {
			_, _ = client.Resize(ctx, &pb.ResizeRequest{
				ProjectId: projectID,
				Rows:      int32(rows),
				Cols:      int32(cols),
			})
		}
	}
}

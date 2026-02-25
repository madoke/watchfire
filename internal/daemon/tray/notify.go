package tray

import (
	"fmt"
	"log"

	"github.com/gen2brain/beeep"
)

func notifyAgentDone(projectName, mode string) {
	title := "Watchfire"
	msg := fmt.Sprintf("%s — %s completed", projectName, mode)
	if err := beeep.Notify(title, msg, ""); err != nil {
		log.Printf("Failed to send notification: %v", err)
	}
}

func notifyAgentError(projectName, mode string) {
	title := "Watchfire"
	msg := fmt.Sprintf("%s — %s stopped", projectName, mode)
	if err := beeep.Notify(title, msg, ""); err != nil {
		log.Printf("Failed to send notification: %v", err)
	}
}

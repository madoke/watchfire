// Package analytics provides a thin wrapper around PostHog for event tracking.
package analytics

import (
	"log"

	"github.com/posthog/posthog-go"
)

var (
	client         posthog.Client
	installationID string
	version        string
)

// Init creates the PostHog client. No-op if apiKey is empty.
func Init(apiKey, instID, ver string) {
	if apiKey == "" {
		return
	}
	installationID = instID
	version = ver

	var err error
	client, err = posthog.NewWithConfig(apiKey, posthog.Config{
		Endpoint:  "https://eu.i.posthog.com",
		BatchSize: 1,
		Verbose:   true,
	})
	if err != nil {
		log.Printf("[analytics] Failed to create PostHog client: %v", err)
		client = nil
	}
}

// Track sends an event to PostHog. No-op if not initialized.
func Track(event string, properties posthog.Properties) {
	if client == nil {
		return
	}
	if properties == nil {
		properties = posthog.NewProperties()
	}
	properties.Set("version", version)

	if err := client.Enqueue(posthog.Capture{
		DistinctId: installationID,
		Event:      event,
		Properties: properties,
	}); err != nil {
		log.Printf("[analytics] Failed to track %s: %v", event, err)
	}
}

// Close flushes pending events and shuts down the client.
func Close() {
	if client != nil {
		_ = client.Close()
	}
}

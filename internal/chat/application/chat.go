package chat

import "context"

// Client defines the application-level contract for chat operations.
type Client interface {
	Connect(ctx context.Context) error
	Say(channel, message string)
}

// Stream defines how chat messages are rendered to a local output.
type Stream interface {
	Message(author, colorHex, content string)
	System(content string)
}

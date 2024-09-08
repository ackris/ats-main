package common

import (
	"io"
)

// MessageFormatter defines an interface for formatting and displaying records read by a consumer instance.
type MessageFormatter interface {
	// Configure sets up the MessageFormatter with the provided configuration.
	Configure(configs map[string]interface{})

	// WriteTo formats the provided ConsumerRecord for display and writes it to the output.
	WriteTo(record *ConsumerRecord, output io.Writer)

	// Close cleans up any resources used by the formatter.
	Close() error
}

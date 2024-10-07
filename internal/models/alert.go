package models

type (
	AlertInterface[T any] interface {
		// Init initializes the alert with the provided configuration.
		ConfigValidation() error
		// Render renders the alert message with the provided data.
		Render() (string, error)
		// GetSpec returns the alert configuration.
		GetSpec() T
	}
)

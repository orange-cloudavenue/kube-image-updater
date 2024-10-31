package models

type (
	ActionName string
)

// String returns the string representation of the action name.
func (n ActionName) String() string {
	return string(n)
}

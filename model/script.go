// Package model defines the domain model.
package model

// Command represents a command to execute against a server.
type Command string

// Script defines a collection of Commands.
type Script struct {
	Commands []Command
	Enabled bool
	RequiresSudo bool
}

// NewScript allocates a new Script.
func NewScript() Script {
	return Script{make([]Command, 0), true, false}
}

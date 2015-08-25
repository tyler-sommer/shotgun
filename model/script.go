package model

type Command string

type Script struct {
	Commands []Command
	Enabled bool
	RequiresSudo bool
}

func NewScript() Script {
	return Script{make([]Command, 0), true, false}
}

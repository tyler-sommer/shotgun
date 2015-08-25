package model

type Command string

type Script struct {
	Commands []Command
	Enabled bool
}

func NewScript() Script {
	return Script{make([]Command, 0), false}
}

type Server struct {
	Host string
	User string
	Scripts []Script
	key string
}

func NewServer() Server {
	return Server{}
}

func (s *Server) Key() string {
	return s.key
}

func (s *Server) SetKey(key string) {
	s.key = key
}

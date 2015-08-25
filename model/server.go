package model

// Server defines an SSH endpoint and a list of scripts to run.
type Server struct {
	Host string
	User string
	Scripts []Script
	key string
}

// NewServer allocates a new Server.
func NewServer() Server {
	return Server{}
}

// Key returns the Server's key.
func (s *Server) Key() string {
	return s.key
}

// SetKey sets the Server's key.
func (s *Server) SetKey(key string) {
	s.key = key
}

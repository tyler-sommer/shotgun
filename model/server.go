package model

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

package http_server

// Option - option
type Option func(server *Server)

// WithAddress - set address for server
func WithAddress(address string) Option {
	return func(s *Server) {
		s.server.Addr = address
	}
}

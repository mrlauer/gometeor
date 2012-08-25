package meteor

import (
	"github.com/mrlauer/gosockjs"
)

type Server struct {
}

func (s *Server) sockjsHandler(c *gosockjs.Conn) {
	// Create and start a Stream, a thin layer on top of sockjs.
	// Also create a session.
	// Hand control to the stream and session.
}

func (s *Server) HandleHTTP(baseUrl string) {
	// The protocol is built on SockJS, so install SockJS handlers.
	gosockjs.Install(baseUrl, func(c *gosockjs.Conn) {
		s.sockjsHandler(c)
	})
}

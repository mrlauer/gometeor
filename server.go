package meteor

import (
	"github.com/mrlauer/gosockjs"
	"log"
)

type Server struct {
	uuid string
}

func (s *Server) sockjsHandler(c *gosockjs.Conn) {
	// Acknowledge.
	ack, err := ToJSON(struct{ Server_id string }{s.uuid})
	if err != nil {
		log.Print(err)
		return
	}
	c.Write(ack)

	// We should probably do this with callbacks.
	// Create and start a Stream, a thin layer on top of sockjs.
	// Also create a session.
	// Hand control to the stream and session.

	session := newSession(s, c)
	session.Run()
}

func (s *Server) HandleHTTP(baseUrl string) {
	// The protocol is built on SockJS, so install SockJS handlers.
	gosockjs.Install(baseUrl, func(c *gosockjs.Conn) {
		s.sockjsHandler(c)
	})
}

func NewServer() *Server {
	s := &Server{}
	s.uuid = uuid()
	return s
}

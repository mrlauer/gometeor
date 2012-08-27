package meteor

import (
	"errors"
	"fmt"
	"github.com/mrlauer/gosockjs"
	"log"
	"reflect"
	"sync"
)

type Server struct {
	uuid string

	functions    map[string]interface{}
	functionLock sync.RWMutex
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
	s.functions = make(map[string]interface{})
	return s
}

func (s *Server) Methods(methods map[string]interface{}) error {
	s.functionLock.Lock()
	defer s.functionLock.Unlock()

	for name := range methods {
		if _, ok := s.functions[name]; ok {
			return fmt.Errorf("There is already a method named %s\n", name)
		}
	}

	for name, function := range methods {
		val := reflect.ValueOf(function)
		if val.Kind() != reflect.Func {
			return errors.New("Not a function.")
		}
		if val.Type().NumOut() > 1 {
			return errors.New("Registered functions may return only one result.")
		}
		s.functions[name] = function
	}
	return nil
}

func (s *Server) GetFunction(name string) interface{} {
	s.functionLock.RLock()
	defer s.functionLock.RUnlock()
	return s.functions[name]
}

package gometeor

import (
	"encoding/json"
	"io"
	"log"
)

type SessionStatus int

const (
	Disconnected SessionStatus = iota
	Connected
	Closed
)

type Session struct {
	server *Server
	conn   io.ReadWriteCloser
	uuid   string
	status SessionStatus
}

func newSession(server *Server, conn io.ReadWriteCloser) *Session {
	s := &Session{}
	s.server = server
	s.uuid = uuid()
	s.conn = conn
	return s
}

func (s *Session) Run() {
	// Read loop
	dec := json.NewDecoder(s.conn)
	for {
		// Read a string from the reader.
		var msgStr json.RawMessage
		err := dec.Decode(&msgStr)
		// Is there a way to distinguish connection errors from
		// json errors?
		if err != nil {
			log.Println(err)
			return
		}
		// Now try to turn it into an object
		var raw RawMessage
		err = json.Unmarshal(msgStr, &raw)
		if err != nil {
			log.Printf("Error %v decoding %v\n", err, msgStr)
		}
		s.process(raw)
	}
}

func (s *Session) SendObject(obj interface{}) error {
	// First turn it into json
	b, err := ToJSON(obj)
	if err != nil {
		log.Println(err)
		return err
	}
	go s.conn.Write(b)
	_, err = s.conn.Write(b)
	return err
}

type msgHeader struct {
	Msg string
}

type msgMethod struct {
	Msg    string
	Method string
	Params []json.RawMessage
	Id     string
}

type msgMethodAck struct {
	Msg     string
	Methods []string
}

type MeteorError struct {
	Error  int    `json:"error"`
	Reason string `json:"reason"`
}

type msgResult struct {
	Msg    string
	Id     string
	Result interface{}
}

type msgResultError struct {
	Msg   string
	Id    string
	Error MeteorError
}

type msgError struct {
	Msg    string
	Reason string
}

func (s *Session) process(m RawMessage) {
	// Message type
	var header msgHeader
	err := m.Decode(&header)
	if err != nil {
		log.Printf("Error %v decoding %v\n", err, m)
		return
	}
	mtype := header.Msg
	if mtype == "connect" {
		s.SendObject(struct {
			Msg     string
			Session string
		}{"connected", s.uuid})
	} else if mtype == "method" {
		var msg msgMethod
		err := m.Decode(&msg)
		if err != nil {
			s.SendObject(msgError{"error", "Malformed method invocation"})
			return
		}
		// This should not happen now, but eh
		s.SendObject(msgMethodAck{
			Msg:     "data",
			Methods: []string{msg.Id},
		})
		fn := s.server.GetFunction(msg.Method)
		if fn == nil {
			s.SendObject(msgResultError{
				Msg:   "result",
				Id:    msg.Id,
				Error: MeteorError{404, "Method not found"},
			})
			return
		}
		result, err := Call(fn, msg.Params)
		if err != nil {
			s.SendObject(msgResultError{
				Msg:   "result",
				Id:    msg.Id,
				Error: MeteorError{500, err.Error()},
			})
			return
		}
		s.SendObject(msgResult{
			Msg:    "result",
			Id:     msg.Id,
			Result: result,
		})
	}
}

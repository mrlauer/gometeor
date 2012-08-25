package meteor

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
)

// uuid returns a new type-4 uuid.
// Stolen directly from https://groups.google.com/forum/?fromgroups=#!msg/golang-nuts/owCogizIuZs/ZzmwkQGrlnEJ
func uuid() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		log.Fatal(err)
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

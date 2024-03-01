package str

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func Generate(length int) string {
	var b = make([]byte, length)

	s := rand.NewSource(time.Now().Unix() * rand.Int63())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		log.Fatalf("string generation error: %s", err.Error())
	}

	return fmt.Sprintf("%x", b)[:length]
}

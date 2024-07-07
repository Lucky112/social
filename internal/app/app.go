package app

import (
	"fmt"

	"github.com/Lucky112/social/internal/transport"
)

func Run() {
	s := transport.NewServer()
	fmt.Println(s.Start(11000))
}

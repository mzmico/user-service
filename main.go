package main

import (
	"github.com/mzmico/mz"
	e "github.com/mzmico/mz/http_service"
	_ "github.com/mzmico/user-service/impls"
)

func main() {

	s := e.Default(
		mz.WithAddress(":8020"))

	s.Run()

}

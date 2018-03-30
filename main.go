package main

import (
	"fmt"

	e "github.com/mzmico/mz/rpc_service"
	_ "github.com/mzmico/user-service/impls"
)

func main() {

	s := e.Default()

	err := s.Run()

	if err != nil {
		fmt.Println(err)
	}

}

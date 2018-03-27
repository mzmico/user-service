package main

import (
	"github.com/mzmico/mz"
	"github.com/mzmico/toolkit/state"
)

import (
	_ "github.com/mzmico/user-service/impls"
)

func main() {

	s := mz.NewHttpService(
		mz.WithAddress(":8020"),
	)

	e := s.Engine()
	e.GET("/", state.GinHandler(func(state *state.HttpState) {

		state.JSON(state.Session())
	}))

	s.Run()

}

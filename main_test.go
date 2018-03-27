package main

import (
	"fmt"
	"testing"

	"google.golang.org/grpc"
)

func Test_main(t *testing.T) {

	_, err := grpc.Dial("127.0.0.1:2680", grpc.WithInsecure())

	if err != nil {
		fmt.Println(err)
		return
	}

}

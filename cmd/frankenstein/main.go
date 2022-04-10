package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/oklog/oklog/pkg/group"
	"github.com/xcaballero/frankenstein/pkg/actor"
)

func main() {
	fmt.Println("Hello world, this is Frankenstein")
	sm := actor.NewStateMachine()
	var api = actor.NewAPI()
	var g group.Group
	{
		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			return sm.Run(ctx)
		}, func(error) {
			cancel()
		})
	}
	{
		ln, _ := net.Listen("tcp", ":8080")
		g.Add(func() error {
			return http.Serve(ln, api)
		}, func(error) {
			ln.Close()
		})
	}
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			return cronJobs(cancel, sm)
		}, func(error) {
			close(cancel)
		})
	}
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			return signalCacher(cancel)
		}, func(error) {
			close(cancel)
		})
	}
	log.Fatal(g.Run())
}

func signalCacher(cancel chan struct{}) error {
	defer close(cancel)
	return nil
}

func cronJobs(cancel chan struct{}, sm actor.StateMachine) error {
	return nil
}

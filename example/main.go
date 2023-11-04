package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dashotv/minion"
)

type TestPayload struct {
	Name  string
	Value int
}

func main() {
	min := setupMinion()
	min.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-time.After(3 * time.Second)
		min.Enqueue("test")
	}()
	go func() {
		<-time.After(5 * time.Second)
		min.EnqueueWithPayload("test", &TestPayload{"test", 1})
	}()

	select {
	case <-c:
		fmt.Println("interrupt")
	case <-time.After(60 * time.Second):
		fmt.Println("done")
	}
}

func setupMinion() *minion.Minion {
	c := 0
	m := minion.New(1)
	m.Register("test", func(payload any) error {
		if payload == nil {
			fmt.Println("test")
			return nil
		}

		fmt.Println("test", payload.(*TestPayload).Name, payload.(*TestPayload).Value)
		return nil
	})
	m.Schedule("* * * * * *", "seconds", func(payload any) error {
		c++
		fmt.Println("seconds:", c)
		return nil
	})
	return m
}

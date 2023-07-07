package main

import "github.com/lookuplogger/translationservice/pkg/server"

func main() {
	s := server.NewServer()
	s.Start()
}

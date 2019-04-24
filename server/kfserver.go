package main

import "server"

func main() {
	s := server.NewServer(":8888")

	for s.Running {
		s.PrintLogs()
	}
}

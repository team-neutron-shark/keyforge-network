package main

import kfnetwork "keyforge-network"

func main() {
	s := kfnetwork.NewServer(":8888")

	for s.Running {
		s.PrintLogs()
	}
}

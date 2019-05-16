package main

import kfnetwork "github.com/team-neutron-shark/keyforge-network"

func main() {
	s := kfnetwork.NewServer(":8888")

	for s.Running {
		s.PrintLogs()
	}
}

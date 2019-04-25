package main

import (
	kfnetwork "keyforge-network"
)

func main() {
	client := kfnetwork.NewClient()
	client.Connect(":8888")
	client.SendVersion()
}

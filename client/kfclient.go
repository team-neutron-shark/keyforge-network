package main

import (
	kfnetwork "keyforge-network"
)

func main() {
	client := kfnetwork.NewClient()
	client.Connect(":8888")
	client.SendVersionRequest()
	client.SendLoginRequest()
	client.SendCreateLobbyRequest()
}

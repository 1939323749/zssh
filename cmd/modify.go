package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"strconv"
)

// Modify the selected server
func modifyServer(selectedServer *Server) {
	prompt := promptui.Prompt{
		Label:     "New ID",
		Default:   selectedServer.ID,
		AllowEdit: true,
	}
	newID, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return
	}

	prompt = promptui.Prompt{
		Label:     "New IP",
		Default:   selectedServer.IP,
		AllowEdit: true,
	}
	newIP, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return
	}

	prompt = promptui.Prompt{
		Label:     "New Username",
		Default:   selectedServer.Username,
		AllowEdit: true,
	}
	newUsername, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return
	}

	prompt = promptui.Prompt{
		Label:     "New Port",
		Default:   fmt.Sprintf("%d", selectedServer.Port),
		AllowEdit: true,
	}
	newPortStr, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return
	}

	newPort, err := strconv.Atoi(newPortStr)
	if err != nil {
		fmt.Printf("Invalid port: %v\n", err)
		return
	}

	// Update the server
	selectedServer.ID = newID
	selectedServer.IP = newIP
	selectedServer.Username = newUsername
	selectedServer.Port = newPort

	// Save the modified servers
	saveServers()

	fmt.Println("Server modified successfully.")
}

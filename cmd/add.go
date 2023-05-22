package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"strconv"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH server",
	Long:  `Add a new SSH server to the list of managed servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt for the server details
		prompt := promptui.Prompt{
			Label: "Username",
		}
		username, err := prompt.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		prompt = promptui.Prompt{
			Label: "IP",
		}
		ip, _ := prompt.Run()

		prompt = promptui.Prompt{
			Label:   "Port",
			Default: "22",
			Validate: func(input string) error {
				_, err := strconv.Atoi(input)
				if err != nil {
					return errors.New("Invalid port number")
				}
				return nil
			},
		}
		portString, _ := prompt.Run()

		// Convert port from string to int
		port, err := strconv.Atoi(portString)
		if err != nil {
			fmt.Println("Invalid port number.")
			return
		}

		prompt = promptui.Prompt{
			Label: "ID",
		}
		id, _ := prompt.Run()

		if len(id) == 0 {
			id = generateID()
		}

		// Add the server to the list
		servers = append(servers, Server{
			ID:       id,
			Username: username,
			IP:       ip,
			Port:     port,
		})

		// Save the updated list of servers
		saveServers()
	},
}

func generateID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

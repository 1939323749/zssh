package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func removeServer(id string) {
	// Remove the selected server
	removeServerByID(id)

	// Save the updated server list to the config file
	saveServers()
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a SSH server",
	Long:  `Remove a SSH server from the list of managed servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		page := 0
		for {
			startIndex := page * itemsPerPage
			endIndex := (page + 1) * itemsPerPage
			if endIndex > len(servers) {
				endIndex = len(servers)
			}

			// Create a list of server labels for the current page
			serverLabels := make([]string, endIndex-startIndex)
			for i := startIndex; i < endIndex; i++ {
				server := servers[i]
				serverLabels[i-startIndex] = fmt.Sprintf("ID: %s, Username: %s, IP: %s, Port: %d", server.ID, server.Username, server.IP, server.Port)
			}

			buttons := make([]string, 0)
			if len(servers) > itemsPerPage {
				if startIndex > 0 {
					buttons = append(buttons, "<< Prev")
				}
				if endIndex < len(servers) {
					buttons = append(buttons, "Next >>")
				}
			}
			buttons = append(buttons, "Exit")
			serverLabels = append(serverLabels, buttons...)

			prompt := promptui.Select{
				Label: "Select a server to remove",
				Items: serverLabels,
			}

			selectedIndex, selectedItem, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			if selectedItem == "Exit" {
				// Exit the program
				return
			} else if selectedIndex >= 0 {
				if selectedItem == "<< Prev" {
					if page > 0 {
						page--
					}
				} else if selectedItem == "Next >>" {
					if endIndex < len(servers) {
						page++
					}
				} else {
					if len(servers) > itemsPerPage {
						if selectedItem == "Next >>" {
							selectedIndex--
						}
						if startIndex > 0 && selectedItem != "<< Prev" {
							selectedIndex--
						}
					}
					selectedIndex += startIndex
					selectedServer := &servers[selectedIndex]

					// Confirm the removal
					confirmPrompt := promptui.Prompt{
						Label:     fmt.Sprintf("Are you sure you want to remove server %s", selectedServer.ID),
						IsConfirm: true,
					}
					_, err = confirmPrompt.Run()
					if err != nil {
						// If the user does not confirm, continue with the loop
						continue
					}

					// Remove the selected server
					removeServer(selectedServer.ID)

					// Check if the current page is empty, and if so, go to the previous page
					if startIndex >= len(servers) && page > 0 {
						page--
					}

				}
			}
		}
	},
}

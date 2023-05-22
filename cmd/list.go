package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"strings"
)

const (
	itemsPerPage = 5
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed SSH servers",
	Long:  `List all SSH servers that are currently being managed.`,
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
				serverLabels = append(serverLabels, fmt.Sprintf("Page %d/%d", page+1, (len(servers)+itemsPerPage-1)/itemsPerPage))
			}
			buttons = append(buttons, "Exit")
			serverLabels = append(serverLabels, buttons...)

			prompt := promptui.Select{
				Label: "Select a server",
				Items: serverLabels,
			}

			selectedIndex, selectedItem, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			if strings.HasPrefix(selectedItem, "Page ") {
				// User selected a page number, display all servers
				serverLabels := make([]string, len(servers))
				for i, server := range servers {
					serverLabels[i] = fmt.Sprintf("ID: %s, Username: %s, IP: %s, Port: %d", server.ID, server.Username, server.IP, server.Port)
				}

				prompt := promptui.Select{
					Label: "Select a server",
					Items: append(serverLabels, "Back"),
				}

				selectedIndex, _, err := prompt.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}

				if selectedIndex < len(servers) {
					selectedServer := &servers[selectedIndex]
					// Let the user select an action
					prompt := promptui.Select{
						Label: "Select an action",
						Items: []string{"Back", "Modify", "Connect", "Remove"},
					}

					_, selectedAction, err := prompt.Run()
					if err != nil {
						fmt.Printf("Prompt failed %v\n", err)
						return
					}

					switch selectedAction {
					case "Back":
						// Go back to the page view
						continue

					case "Modify":
						// Modify the selected server
						modifyServer(selectedServer)

					case "Connect":
						// Connect to the selected server
						connectToServer(selectedServer)

					case "Remove":
						// Remove the selected server
						removeServer(selectedServer)
						// Remove the server from the list
						servers = append(servers[:selectedIndex], servers[selectedIndex+1:]...)
					}
				}
				continue
			} else if selectedItem == "Exit" {
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

					// Let the user select an action
					prompt := promptui.Select{
						Label: "Select an action",
						Items: []string{"Exit", "Modify", "Connect", "Remove"},
					}

					_, selectedAction, err := prompt.Run()
					if err != nil {
						fmt.Printf("Prompt failed %v\n", err)
						return
					}

					switch selectedAction {
					case "Exit":
						// Exit the program
						return

					case "Modify":
						// Modify the selected server
						modifyServer(selectedServer)

					case "Connect":
						// Connect to the selected server
						connectToServer(selectedServer)

					case "Remove":
						// Connect to the selected server
						removeServer(selectedServer)
						servers = append(servers[:selectedIndex], servers[selectedIndex+1:]...)
					}
				}
			}
		}
	},
}

// Register the command
func init() {
	rootCmd.AddCommand(listCmd)
}

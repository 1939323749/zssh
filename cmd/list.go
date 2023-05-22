package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed SSH servers",
	Long:  `List all SSH servers that are currently being managed.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a slice of server strings
		serverStrings := make([]string, len(servers))
		for i, server := range servers {
			serverStrings[i] = fmt.Sprintf("ID: %s, Username: %s, IP: %s, Port: %d", server.ID, server.Username, server.IP, server.Port)
		}

		// Create a new promptui.Select
		prompt := promptui.Select{
			Label: "Select Server",
			Items: serverStrings,
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You have selected: %s\n", result)
	},
}

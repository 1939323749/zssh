package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a SSH server",
	Long:  `Remove a SSH server from the list of managed servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide the ID of the server to remove.")
			return
		}

		// Remove the server
		removeServerByID(args[0])

		// Save the updated list of servers
		saveServers()
	},
}

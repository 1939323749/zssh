package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"os"
)

// Connect to the server
func connectToServer(server *Server) {
	// Prompt for password
	prompt := promptui.Prompt{
		Label: "Password",
		Mask:  '*',
	}
	password, err := prompt.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Connect to the server
	config := &ssh.ClientConfig{
		User: server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.IP, server.Port), config)
	if err != nil {
		fmt.Println("Failed to dial: ", err)
		return
	}

	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: ", err)
		return
	}

	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	defer session.Close()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		fmt.Println("request for pseudo terminal failed: ", err)
		return
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		fmt.Println("failed to start shell: ", err)
		return
	}

	session.Wait()
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a SSH server",
	Long:  `Connect to a SSH server using its ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide the ID of the server to connect to.")
			return
		}

		// Find the server
		server := findServerByID(args[0])
		if server == nil {
			fmt.Printf("No server found with ID %s.\n", args[0])
			return
		}

		connectToServer(server)
	},
}

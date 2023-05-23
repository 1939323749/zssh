package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
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
		fmt.Printf("Prompt failed %v\n", err)
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

	// Start a new session
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: ", err)
		return
	}

	defer session.Close()

	// Get terminal size
	fd := int(os.Stdin.Fd())
	width, height, err := terminal.GetSize(fd)
	if err != nil {
		fmt.Println("Failed to get terminal size: ", err)
		return
	}

	// Request a pseudo terminal on the server with terminal size
	err = session.RequestPty("xterm-256color", height, width, ssh.TerminalModes{})
	if err != nil {
		fmt.Println("Request for pseudo terminal failed: ", err)
		return
	}

	// Connect the session to Stdin, Stdout, and Stderr
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// Set terminal in raw mode
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		fmt.Println("Failed to set terminal in raw mode: ", err)
		return
	}

	// Restore terminal state at the end
	defer terminal.Restore(fd, oldState)

	// Start the shell on the remote server with terminal raw mode
	err = session.Shell()
	if err != nil {
		fmt.Println("Failed to start shell: ", err)
		return
	}

	// Wait for the session to finish
	err = session.Wait()
	if err != nil && err != io.EOF {
		fmt.Println("Failed to wait for session: ", err)
		return
	}
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

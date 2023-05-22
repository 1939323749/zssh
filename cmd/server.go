package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Server struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
}

const (
	defaultSSHConfigFile = ".sshManager.json"
	defaultPort          = 22
)

var (
	cfgFile string
	servers = make([]Server, 0)
)

func initConfig() {
	// Load the servers from the config file
	loadServers()
}

func saveServers() {
	// Convert the servers slice to JSON
	data, err := json.Marshal(servers)
	if err != nil {
		panic(err)
	}

	// Write the JSON data to the config file
	err = ioutil.WriteFile(defaultSSHConfigFile, data, 0644)
	if err != nil {
		panic(err)
	}
}

func loadServers() {
	// Check if the config file exists
	if _, err := os.Stat(defaultSSHConfigFile); os.IsNotExist(err) {
		// If not, create an empty one
		f, err := os.Create(defaultSSHConfigFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		// Write an empty JSON array to the new file
		if _, err := f.WriteString("[]"); err != nil {
			panic(err)
		}
	}

	// Read the data from the config file
	data, err := ioutil.ReadFile(defaultSSHConfigFile)
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON data into the servers slice
	err = json.Unmarshal(data, &servers)
	if err != nil {
		panic(err)
	}
}

func findServerByID(id string) *Server {
	for i, s := range servers {
		if s.ID == id {
			return &servers[i]
		}
	}
	return nil
}

func removeServerByID(id string) {
	for i, s := range servers {
		if s.ID == id {
			servers = append(servers[:i], servers[i+1:]...)
			return
		}
	}
}

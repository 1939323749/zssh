package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type Server struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	IP       string `json:"ip"`
}

type Config struct {
	Servers []Server `json:"servers"`
}

func ReadConfig() (Config, error) {
	file, err := os.Open("~/.sshManager.json")
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	config := Config{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func WriteConfig(config Config) error {
	file, err := os.Create("~/.sshManager.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	err = json.NewEncoder(file).Encode(config)
	if err != nil {
		return err
	}

	return nil
}

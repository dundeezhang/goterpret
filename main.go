package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	PromptColor string
	HomeDir     string
}

func loadConfig() Config {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return Config{}
	}

	configPath := filepath.Join(usr.HomeDir, ".gorc")
	file, err := os.Open(configPath)
	if err != nil {
		return Config{PromptColor: "\033[32m", HomeDir: usr.HomeDir}
	}
	defer file.Close()

	config := Config{PromptColor: "\033[32m", HomeDir: usr.HomeDir}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "prompt_color":
			config.PromptColor = value
		case "home_dir":
			config.HomeDir = value
		}
	}

	return config
}

func main() {
	config := loadConfig()
	os.Chdir(config.HomeDir)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s$ \033[0m", config.PromptColor)
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

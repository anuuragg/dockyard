package main

import (
	"fmt" 
	"os"
	"os/exec"
	"strings"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dockyard <command>")
		return
	}

	command := os.Args[1]

	if command == "deploy"{
		if len(os.Args) < 3 {
			fmt.Println("Usage: dockyard <command>")
			return
		}

		path := os.Args[2]

		appName := filepath.Base(path)
		fmt.Println("App name:", appName)

		dockerfile := path + "/Dockerfile"

		_, err := os.Stat(dockerfile)

		if err != nil {
			fmt.Println("Dockerfile not found")
			return
		}

		fmt.Println("Building Docker image...")

		cmd := exec.Command(
			"docker",
			"build",
			"-t",
			appName,
			path,
		)

		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Println("Docker build failed:")
			fmt.Println(string(output))
			return
		}

		fmt.Println(string(output))
		fmt.Println("Build successful!")

		fmt.Println("starting container...")

		cmd = exec.Command(
			"docker",
			"run",
			"-d",
			"-p",
			"0:8000",
			appName,
		)

		output, err = cmd.CombinedOutput()

		if err != nil {
			fmt.Println("Docker run failed:")
			fmt.Println(string(output))
			return
		}

		containerID := string(output)

		fmt.Println("Container started!")
		fmt.Println("Container ID:", containerID)

		cmd = exec.Command(
			"docker",
			"port",
			containerID,
			"8000",
		)

		output, err = cmd.CombinedOutput()

		if err != nil {
			fmt.Println("Port couldn't be found")
			fmt.Println(string(output))
			return
		}

		portOutput := string(output)

		lines := strings.Split(portOutput, "\n")

		port := strings.Split(lines[0], ":")[1]

		fmt.Println("Assignd port:", port)

	} else {
		fmt.Println("Unknown command: ", command)
	}
}

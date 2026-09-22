package docker

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"os"
	"net/http"
	"time"
)

func Build(imageName string, path string) error {
	cmd := exec.Command(
		"docker",
		"build",
		"-t",
		imageName,
		path,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf(
			"docker build failed: %s",
			string(output),
		)
	}

	fmt.Println(string(output))

	return nil
}

func Run(imageName string) (string, error) {
	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"--restart",
		"unless-stopped",
		"-p",
		"0:8000",
		imageName,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf(
			"docker run failed: %s",
			string(output),
		)
	}

	containerID := strings.TrimSpace(string(output))

	return containerID, nil
}

func GetPort(containerID string) (int, error) {
	cmd := exec.Command(
		"docker",
		"port",
		containerID,
		"8000",
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return 0, fmt.Errorf(
			"couldn't get container port: %s",
			string(output),
		)
	}

	portOutput := strings.TrimSpace(string(output))

	lines := strings.Split(portOutput, "\n")

	port := strings.Split(lines[0], ":")[1]

	portNumber, err := strconv.Atoi(port)

	if err != nil {
		return 0, fmt.Errorf("invalid port: %s", port)
	}

	return portNumber, nil
}


func Stop(containerID string) error {
	cmd := exec.Command(
		"docker",
		"stop",
		containerID,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf(
			"docker stop failed: %s",
			string(output),
		)
	}

	return nil
}

func Logs(containerID string) error {
	cmd := exec.Command(
		"docker",
		"logs",
		containerID,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func HealthCheck(port int) error {
	url := fmt.Sprintf("http://localhost:%d", port)

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	for i := 0; i < 5; i++ {
		resp, err := client.Get(url)

		if err == nil {
			resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				return nil
			}
		}

		time.Sleep(time.Second)
	}

	return fmt.Errorf("health check failed after 5 seconds")
}
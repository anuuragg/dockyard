package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/anuuragg/dockyard/models"
)

func GenerateCaddyConfig(apps []models.App) string {
	var config strings.Builder

	for _, app := range apps {
		block := fmt.Sprintf(`%s.yourdomain.local {
	reverse_proxy localhost:%d
}

`, app.Name, app.Port)

		config.WriteString(block)
	}

	return config.String()
}

func WriteCaddyConfig(config string) error {
	return os.WriteFile(
		"/etc/caddy/Caddyfile",
		[]byte(config),
		0644,
	)
}

func ReloadCaddy() error {
	cmd := exec.Command(
		"sudo",
		"systemctl",
		"reload",
		"caddy",
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf(
			"failed to reload Caddy: %s",
			string(output),
		)
	}

	return nil
}

func AddHostEntry(appName string) error {
	file, err := os.OpenFile(
		"/etc/hosts",
		os.O_APPEND|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return err
	}

	defer file.Close()

	entry := fmt.Sprintf(
		"\n127.0.0.1 %s.yourdomain.local\n",
		appName,
	)

	_, err = file.WriteString(entry)

	return err
}
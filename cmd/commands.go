package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anuuragg/dockyard/database"
	"github.com/anuuragg/dockyard/docker"
	"github.com/anuuragg/dockyard/models"
	"github.com/anuuragg/dockyard/proxy"
)

func Deploy(db *sql.DB, path string) error {
	appName := filepath.Base(path)

	fmt.Println("App name:", appName)

	dockerfile := filepath.Join(path, "Dockerfile")

	_, err := os.Stat(dockerfile)

	if err != nil {
		return fmt.Errorf("Dockerfile not found")
	}

	fmt.Println("Building Docker image...")

	err = docker.Build(appName, path)

	if err != nil {
		return err
	}

	fmt.Println("Build successful!")

	fmt.Println("Starting container...")

	containerID, err := docker.Run(appName)

	if err != nil {
		return err
	}

	fmt.Println("Container started!")
	fmt.Println("Container ID:", containerID)

	port, err := docker.GetPort(containerID)

	if err != nil {
		return err
	}

	fmt.Println("Assigned port:", port)

	app := models.App{
		Name:        appName,
		ContainerID: containerID,
		Port:         port,
	}

	err = database.InsertApp(db, app)

	if err != nil {
		return fmt.Errorf("failed to save app: %w", err)
	}

	fmt.Println("Data inserted into database!")

	err = proxy.AddHostEntry(appName)

	if err != nil {
		return fmt.Errorf("failed to add host entry: %w", err)
	}

	fmt.Println("Host entry added!")

	apps, err := database.GetApps(db)

	if err != nil {
		return fmt.Errorf("failed to get apps: %w", err)
	}

	config := proxy.GenerateCaddyConfig(apps)

	err = proxy.WriteCaddyConfig(config)

	if err != nil {
		return fmt.Errorf("failed to write Caddy config: %w", err)
	}

	fmt.Println("Caddy config updated!")

	err = proxy.ReloadCaddy()

	if err != nil {
		return err
	}

	fmt.Println("Caddy reloaded!")

	return nil
}
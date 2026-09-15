package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type App struct {
	Name        string
	ContainerID string
	Port        int
}



 // open sqlite database
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dockyard <command>")
		return
	}

	command := os.Args[1]

	db, err := sql.Open("sqlite", "dockyard.db")

	if err != nil {
		fmt.Println("Failed to open database:", err)
		return
	}

	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS apps (
			name TEXT PRIMARY KEY,
			container_id TEXT NOT NULL,
			port INTEGER NOT NULL
		)
	`)

	if err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}

	fmt.Println("Database ready!")



	// listen for deploy command
	if command == "deploy" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: dockyard deploy <path>")
			return
		}

		path := os.Args[2]



		// get app name from path
		appName := filepath.Base(path)
		fmt.Println("App name:", appName)

		dockerfile := path + "/Dockerfile"



		// check if Dockerfile exists
		_, err := os.Stat(dockerfile)

		if err != nil {
			fmt.Println("Dockerfile not found")
			return
		}

		fmt.Println("Building Docker image...")



		// build Docker image
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

		fmt.Println("Starting container...")



		// start Docker container and assign a random host port
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

		containerID := strings.TrimSpace(string(output))

		fmt.Println("Container started!")
		fmt.Println("Container ID:", containerID)



		// get the host port assigned to the container
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

		portOutput := strings.TrimSpace(string(output))

		lines := strings.Split(portOutput, "\n")

		port := strings.Split(lines[0], ":")[1]

		fmt.Println("Assigned port:", port)

		portNumber, err := strconv.Atoi(port)

		if err != nil {
			fmt.Println("Invalid port:", port)
			return
		}



		// save deployment information to sqlite
		_, err = db.Exec(
			"INSERT INTO apps (name, container_id, port) VALUES (?, ?, ?)",
			appName,
			containerID,
			portNumber,
		)

		if err != nil {
			fmt.Println("Data insertion failed")
			fmt.Println("DB error:", err)
			return
		}

		fmt.Println("Data inserted into the Database!")



		// add app domain to /etc/hosts
		err = addHostEntry(appName)

		if err != nil {
			fmt.Println("Failed to add host entry:", err)
			return
		}

		fmt.Println("Host entry added!")



		// get all deployed apps from sqlite
		apps, err := getApps(db)

		if err != nil {
			fmt.Println("Failed to get apps:", err)
			return
		}



		// generate Caddy configuration
		config := generateCaddyConfig(apps)

		err = writeCaddyConfig(config)

		if err != nil {
			fmt.Println("Failed to write Caddy config:", err)
			return
		}

		fmt.Println("Caddy config updated!")



		// reload Caddy
		cmd = exec.Command(
			"sudo",
			"systemctl",
			"reload",
			"caddy",
		)

		err = cmd.Run()

		if err != nil {
			fmt.Println("Failed to reload Caddy:", err)
			return
		}

		fmt.Println("Caddy reloaded!")

	} else {
		fmt.Println("Unknown command:", command)
	}
}



 // get all apps from database
func getApps(db *sql.DB) ([]App, error) {
	rows, err := db.Query(
		"SELECT name, container_id, port FROM apps",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var apps []App

	for rows.Next() {
		var app App

		err := rows.Scan(
			&app.Name,
			&app.ContainerID,
			&app.Port,
		)

		if err != nil {
			return nil, err
		}

		apps = append(apps, app)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return apps, nil
}



 // generate Caddy config for all apps
func generateCaddyConfig(apps []App) string {
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



 // write generated config to Caddyfile
func writeCaddyConfig(config string) error {
	err := os.WriteFile(
		"/etc/caddy/Caddyfile",
		[]byte(config),
		0644,
	)

	return err
}



 // add app domain to /etc/hosts
func addHostEntry(appName string) error {
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
package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/anuuragg/dockyard/cmd"
	"github.com/anuuragg/dockyard/database"

	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dockyard <command>")
		return
	}

	db, err := sql.Open("sqlite", "dockyard.db")

	if err != nil {
		fmt.Println("Failed to open database:", err)
		return
	}

	defer db.Close()

	err = database.Init(db)

	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		return
	}

	command := os.Args[1]

	switch command {

	case "deploy":

		if len(os.Args) < 3 {
			fmt.Println("Usage: dockyard deploy <path>")
			return
		}

		err := cmd.Deploy(db, os.Args[2])

		if err != nil {
			fmt.Println("Deploy failed:", err)
		}

	default:

		fmt.Println("Unknown command:", command)
	}
}
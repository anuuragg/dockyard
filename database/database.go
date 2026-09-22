package database

import (
	"database/sql"

	"github.com/anuuragg/dockyard/models"
)

func Init(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS apps (
			name TEXT PRIMARY KEY,
			container_id TEXT NOT NULL,
			port INTEGER NOT NULL
		)
	`)

	return err
}

func InsertApp(db *sql.DB, app models.App) error {
	_, err := db.Exec(
		"INSERT INTO apps (name, container_id, port) VALUES (?, ?, ?)",
		app.Name,
		app.ContainerID,
		app.Port,
	)

	return err
}

func GetApps(db *sql.DB) ([]models.App, error) {
	rows, err := db.Query(
		"SELECT name, container_id, port FROM apps",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var apps []models.App

	for rows.Next() {
		var app models.App

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


func GetApp(db *sql.DB, name string) (models.App, error) {
	var app models.App

	err := db.QueryRow(
		"SELECT name, container_id, port FROM apps WHERE name = ?",
		name,
	).Scan(
		&app.Name,
		&app.ContainerID,
		&app.Port,
	)

	return app, err
}
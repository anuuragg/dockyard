# Dockyard

A lightweight local PaaS built with Go that deploys Dockerized applications,
assigns them available ports, and exposes them through local subdomains using
Caddy.

Given a directory containing a Dockerfile, Dockyard builds the image, starts
the container, performs a health check, stores the deployment details in
SQLite, and configures Caddy to route traffic to the container.

For example:

```bash
sudo dockyard deploy ./path
```

creates a local deployment accessible at:

```text
http://path.yourdomain.local
```

## Folder structure

```text
dockyard/
├── cmd/
│   └── commands.go        CLI commands and deployment workflow
├── database/
│   └── database.go        SQLite operations
├── docker/
│   └── docker.go          Docker build, run, logs, stop, and health check
├── models/
│   └── app.go             Application data structure
├── proxy/
│   └── caddy.go           Caddy and /etc/hosts configuration
├── path/
│   └── Dockerfile         Test application
├── main.go                CLI entry point
├── install.sh             Installation script
├── go.mod                 Go module definition
├── go.sum                 Go dependencies
├── .gitignore
└── README.md
```

## How it works

The deployment flow is:

```text
Dockerfile
    ↓
Docker build
    ↓
Docker run
    ↓
Dynamic host port
    ↓
Health check
    ↓
SQLite
    ↓
/etc/hosts
    ↓
Caddy configuration
    ↓
Local subdomain
```

Docker containers are started with:

```bash
docker run -d --restart unless-stopped -p 0:8000 <image>
```

Using `0:8000` allows Docker to automatically select an available host port.

The assigned port and container ID are stored in SQLite:

```text
name
container_id
port
```

Caddy is then configured to reverse proxy the local subdomain to that port:

```text
path.yourdomain.local {
    reverse_proxy localhost:32776
}
```

## Commands

Deploy an application:

```bash
sudo dockyard deploy ./path
```

List deployed applications:

```bash
sudo dockyard list
```

View application logs:

```bash
sudo dockyard logs path
```

Stop an application:

```bash
sudo dockyard stop path
```

## Installation

Clone the repository:

```bash
git clone <repository-url>
cd dockyard
```

Run the installation script:

```bash
chmod +x install.sh
sudo ./install.sh
```

The installer checks for Docker and Caddy, installs them if required, builds
Dockyard, and installs the binary to `/usr/local/bin/dockyard`.

After installation:

```bash
sudo dockyard deploy ./path
```

## Tech stack

* Go
* Docker
* Caddy
* SQLite
* Linux

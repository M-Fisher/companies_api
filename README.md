# Companies Service
A single service for managing Companies

## 📁 Project structure

### `/cmd`

Starting point for the app

### `/app`

Application code

### `/app/api`

API configuration.
You can add new endpoints to the app/api/enpoints
Don't forget to register them in app/api/api.go

### `/dev`

Everything related to environment. Config files, Dockerfiles and docker-compose file for local running.

### `/migrations`
DB migrations files.

## 🚜 Running

### Docker-compose

For development purpose project can be started using [docker-compose](https://docs.docker.com/compose/)

```bash
make dev-up-app
```

Environment variables can be changed in _dev/compose.dev.env_.

### Local run

Optionally you can start environment using docker-compose, and start app in terminal:

```bash
make dev-local-env
make local-run
```
Environment variables can be changed in _dev/local.env_.

### Example JWT Token
For testing purposes you can use following authorization token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0

Also region checking is disabled for Dev environment. To enable checking - set DEVELOPMENT_MODE env to false 


### Reboot
docker-compose environment can be restarted using `make dev-restart`.

## 🛠 Building

Project can be build as binary file (`make build`) or as Docker image: `make @build`

## 👨🏼‍💻 Development

`make setup` - install all required dependencies for the project (e.g. golangci-lint)

`make lint` - linters local run (golangci-lint should be installed)

`make test` - running tests locally

## 🧪 Testing

You can run tests using

```bash
make test
```

## 📌 External dependencies
Infrastructure:
- PostgreSQL
- Kafka
  
Resoures:
- https://ipapi.co/


# Fipe Project

This project is a web application that provides an API to query vehicle information from the FIPE (Fundação Instituto de Pesquisas Econômicas) table. It also includes a frontend to interact with the API.

## Project Overview

The backend is written in Go and uses a MongoDB database to store the vehicle information. The API provides endpoints to:

- Get reference tables
- Get vehicle brands
- Get vehicle models
- Get vehicle years and prices
- Get a dashboard comparing vehicle data between two periods

The frontend is built with HTML, CSS, and JavaScript, and it allows users to interact with the API to query vehicle information.

## Project Structure

- **`.air.toml`**: Configuration file for `air`, a live-reloading tool for Go applications.
- **`.github/`**: Contains GitHub Actions workflows.
- **`Dockerfile`**: Defines the Docker container for the Go application.
- **`docker-compose.yaml`**: Configures the services for the project, including the Go application, a MongoDB database, and a mongo-express instance.
- **`docs/`**: Contains additional documentation.
- **`frontend/`**: Contains the frontend files (HTML, CSS, and JavaScript).
- **`go.mod`** and **`go.sum`**: Manage the project's Go dependencies.
- **`internal/`**: Contains the internal Go source code.
  - **`database/`**: Handles the connection to the MongoDB database.
  - **`handlers/`**: Contains the logic for handling API requests.
  - **`models/`**: Defines the data structures used in the application.
  - **`routes/`**: Defines the API routes.
  - **`utils/`**: Contains utility functions.
- **`main.go`**: The entry point of the Go application.

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/your-username/fipe-project.git
   cd fipe-project
   ```

2. **Run the application:**

   ```bash
   docker-compose up -d
   ```

   This will start the following services:
   - `go_app`: The Go application, accessible at `http://localhost:8080`.
   - `mongo`: The MongoDB database, accessible at `mongodb://localhost:27017`.
   - `mongo-express`: A web-based MongoDB admin interface, accessible at `http://localhost:8081`.

## API Endpoints

The API is available under the `/api` prefix.

- `GET /api/tabelas`: Get reference tables.
- `GET /api/marcas?tabela=<tabela_id>`: Get vehicle brands for a given reference table.
- `GET /api/modelos/{marca}?tabela=<tabela_id>`: Get vehicle models for a given brand and reference table.
- `GET /api/veiculos?modelo=<modelo_id>&tabela=<tabela_id>`: Get vehicle years and prices for a given model and reference table.
- `GET /api/dashboard?tabela1=<tabela1_id>&tabela2=<tabela2_id>&marca=<marca_id>`: Get a dashboard comparing vehicle data between two periods for a specific brand.
- `GET /api/0km?tabela=<tabela_id>`: Get all new vehicles for a given reference table.

## Frontend

The frontend is served from the `frontend/` directory and is accessible at `http://localhost:8080`. It provides a user interface to interact with the API.

- **`index.html`**: The main page for querying vehicle information.
- **`dashboard.html`**: The page for the dashboard feature.
- **`comparar.html`**: A page for comparing vehicles.

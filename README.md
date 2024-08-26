# todoapp

A simple Backend Todo application built with Go, using Gin and GORM. Frontend with Svelte.

## Features

- CRUD for Todo items
- MySQL database
- Dockerized deployment
- Github CI
- Calendar view
- Todo list sorting, filtering, search
- User management by admin
- Assignee of task
- Websocket notification on task updates


## Quick Start

### Prerequisites

- Docker
- Docker Compose
- Git
- Golang
- Svelte

## Steps to Start

1. Clone the repository
2. Create `.env` file in the project root

```
DB_USER=todouser
DB_PASS=todopassword
DB_NAME=tododb
DB_ROOT_PASSWORD=rootpassword

BACKEND_PORT=8080
BACKEND_URL=http://localhost:8080

FRONTEND_PORT=3000
FRONTEND_URL=http://localhost:3000

PORT=8080
JWT_KEY=key
```

3. Start/Stop the application using Docker Compose:

To Start
```
cd frontend
npm install
npm install -g rollup
npm run build

docker-compose up

```

Default the backend will crate an admin account now:
username: admin
password: admin123

To Stop
```
docker-compose down -v

```

## API Endpoints

API available at `http://localhost:8080`

- `GET /todos`: List all todos
- `GET /todos/:id`: Get a specific todo
- `POST /todos`: Create a new todo
- `PUT /todos/:id`: Update an existing todo
- `DELETE /todos/:id`: Delete a todo

### Example: Creating a Todo

POST to `/todos`:

```json
{
  "name": "project",
  "description": "Finish the todo application",
  "due_date": "2023-12-31T23:59:59Z",
  "status": "pending"
}
```
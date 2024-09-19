# Blogging Platform API

TODO api made using Go, Fiber, Postgres

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`DB_HOST`

`DB_PORT`

`DB_USER`

`DB_NAME`

`DB_PASSWORD`

`JWT_SECRET`

## Run Locally

Clone the project

```bash
  git clone https://github.com/hamidhasibul/todo-api-go.git
```

Go to the project directory

```bash
  cd <project-directory>
```

Start the server

```bash
  go run main.go
```

project will run on http://localhost:3000

## Endpoints

### User Registration

```
POST /register
{
  "name": "John Doe",
  "email": "john@doe.com"
  "password": "password"
}
```

### User Login

```
POST /login
{
  "email": "john@doe.com",
  "password": "password"
}
```

### Create a TODO

```
POST /todos
Authorization: Bearer <token>
{
  "title": "Buy groceries",
  "description": "Buy milk, eggs, and bread"
}
```

### Update a TODO

```
PUT /todos/:todoId
Authorization: Bearer <token>
{
  "title": "Buy groceries",
  "description": "Buy milk, eggs, and bread"
}
```

### Delete a TODO

```
DELETE /todos/:todoId
Authorization: Bearer <token>
```

### Get TODOs

```
GET /todos?page=1&limit=10
Authorization: Bearer <token>
```

## URL

- [roadmap.sh](https://roadmap.sh/projects/todo-list-api)

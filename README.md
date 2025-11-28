# Project Chirpy

Chirpy is a simple social network similar to X (Twitter).

## Motivation

This project was built as part of a guided course to learn Go web development.  
I implemented features like authentication, routing, and data persistence to practice building a real-world style backend.

## Features

- User registration and login
- Create, edit, and delete short text posts ("chirps")
- View a feed of all chirps
- Follow other users and see only their chirps in your feed
- Filter chirps by user

## Prerequisites

- Install Go: https://go.dev/dl/

Check your Go installation:

```bash
go version
```

## Configuration

Chirpy uses PostgreSQL and environment variables for configuration.

### Environment Variables

Create a `.env` file in the project root (or set these in your shell):

```env
DB_URL=postgres://chirpy_user:YOUR_PASSWORD@localhost:5432/chirpy?sslmode=disable
SECRET=some-long-random-secret
PLATFORM=local
POLKA_KEY=your-polka-api-key
```

SECRET is required – the server will exit if it’s missing.

### Database Setup

1. Install PostgreSQL
   https://www.postgresql.org/download/

2. Create a database and user (example):

```bash
createdb chirpy
createuser chirpy_user --pwprompt
```

3. Set the DB_URL environment variables before running the app(if not using .env):

```bash
export DB_URL="postgres://chirpy_user:YOUR_PASSWORD@localhost:5432/chirpy?sslmode=disable"
```

4. Run database migrations (if your project has a migration command, e.g.):

```bash
go run . migrate
```

## Tech Stack

- Go
- PostgreSQL
- `database/sql` with `lib/pq`
- `godotenv` for environment variable loading
- Standard library `net/http` for HTTP server and routing

## Installation

1. Clone the repo:

```bash
git clone https://github.com/Ikit24/Chirpy.git
cd Chirpy
```

2. Install dependencies:

```bash
go mod tidy
go mod download
```

3. Run the application:

```bash
go run .
```

## Usage

Once the server is running, open:

http://localhost:8080

in your browser or use it as the base URL for API requests.

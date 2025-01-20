# Movie Festival API

## Project Overview
This project is a simple application using Golang, Gin, and Supabase that allows users to manage movie data.

## Prerequisites
Before running the project, ensure that you have set up a Supabase account and created a new project. You will need to initiate the `movies` table using the following SQL command:

```sql
CREATE TABLE movies (
  id VARCHAR(255) PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  image VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  duration INTEGER NOT NULL,
  genres VARCHAR(255) NOT NULL,
  artists VARCHAR(255) NOT NULL,
  url VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  release_date DATE NULL,
  rating DECIMAL(3, 2) NULL,
  view_count INTEGER NULL
);
```

## Running the Project
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd movie-festival-api
   ```

2. Install the necessary dependencies:
   ```bash
   go mod tidy
   ```

3. Copy the `.env.example` file to `.env` and adjust the `POSTGRES_CONNECTION_STRING` variable:
   ```bash
   cp .env.example .env
   # Edit the .env file to set your PostgreSQL connection string
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

5. Access the API at `http://localhost:8080`.

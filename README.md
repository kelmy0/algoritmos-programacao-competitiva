# Competitive programming algorithms

This project is divided into two main parts.

- `/frontend`: Static SPA application built with Svelte 5 and Tailwind CSS (under development).
- `/backend`: API and server logic (under development).

## 🚀 Getting Started

First, clone the repository to your local machine:

```sh
git clone https://github.com/kelmy0/algoritmos-programacao-competitiva.git
cd algoritmos-programacao-competitiva
```

## How to run the Frontend
``` sh
cd frontend
npm install
npm run dev
```

# How to run the Backend

### Prerequisites
Before running the backend, make sure you have installed:
- [Go](https://go.dev/doc/install) (version 1.26 or higher)
- [PostgreSQL](https://www.postgresql.org/download/) running on your machine

### 1. Setup the Database
1. Open your PostgreSQL client and create a new database:
   ```sql
   CREATE DATABASE "algoritmos-programacao-competitiva";
   ```

### 2. Environment Variables
1. Navigate to the backend folder:
```sh
cd backend
```

2. Duplicate the .env.example file and rename it to .env:
```sh
cp .env.example .env
```

3. Open the .env file and adjust the DATABASE_URL with your local PostgreSQL credentials (user and password).

### 3. Install Dependencies
Inside the backend directory, run the following command to download the required Go packages:
```sh
go mod tidy
```

### 4. Seed the Database (Optional)
```sh
go run main.go --seed
```

_(This will automatically run database migrations and insert the default seeds)._


### 5. Run the Server
```sh
go run main.go
```
The server will be running at http://localhost:8080 or in another port if you changed it.
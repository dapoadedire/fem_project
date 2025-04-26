# Frontend Masters - Complete Go for Professional Developers

[Frontend Masters - Complete Go for Professional Developers](https://frontendmasters.com/courses/complete-go)

[x] Done

## How To Run

### Prerequisites

- Go installed on your machine
- Docker and Docker Compose installed

### Steps to Run the Project

1. **Start the PostgreSQL Database**

   Start the database services using Docker Compose:

   ```
   docker-compose up -d
   ```

   This will start:

   - Main PostgreSQL database (workoutDB) on port 5432
   - Test PostgreSQL database (workoutDB_test) on port 5433

2. **Run the Go Application**

   After the databases are running, start the application:

   ```
   go run main.go
   ```

   By default, the server runs on port 8080. You can specify a different port:

   ```
   go run main.go -port 3000
   ```

3. **API Endpoints**

   The project includes API documentation in the `fem_project_api_docs` directory with Bruno files for:

   - Health Check
   - User Registration
   - Token Creation
   - CRUD operations for Workouts

### Configuration

- Database connection: `localhost:5432`
- Database credentials:
  - Username: postgres
  - Password: postgres
  - Database name: postgres
- The application automatically runs migrations at startup

### For Testing

The test database runs on port 5433 and can be used for running test cases.

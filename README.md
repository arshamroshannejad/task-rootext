# Task-Rootext API

Task-Rootext is a RESTful API built with Go, designed to manage posts, user authentication, and voting on posts. The API supports features like pagination, sorting, and filtering for posts, as well as user registration, login, and JWT-based authentication. It uses PostgreSQL as the primary database, Redis for caching, and Docker for containerization.

## Features

- **User Authentication**: Register, login, and logout with JWT-based authentication.
- **Post Management**: Create, read, update, and delete posts.
- **Voting System**: Users can upvote or downvote posts.
- **Pagination & Sorting**: Fetch posts with pagination, sorting, and filtering options.
- **Dockerized**: Easy to set up and run using Docker Compose.

## Technologies Used

- **Go**: The primary programming language.
- **PostgreSQL**: The main database for storing posts and user data.
- **Redis**: Used for caching and session management.
- **Docker**: Containerization for easy deployment and development.
- **Chi**: A lightweight, idiomatic, and composable router for building Go HTTP services.
- **Swagger**: API documentation.

---

## Installation and Setup

### Prerequisites

- Docker and Docker Compose installed on your machine.
- Go (if you want to run the project locally without Docker).

### Steps to Run the Project

#### 1. Clone the Repository

```bash
git clone https://github.com/arshamroshannejad/task-rootext.git
cd task-rootext
```

#### 2. Create `.env` File

Copy the `.env-sample` file to `.env` and update the environment variables as needed:

```bash
cp .env-sample .env
```

Example `.env` file:

```env
POSTGRES_USER=your_db_user
POSTGRES_PASSWORD=your_db_password
POSTGRES_DB=your_db_name
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
REDIS_PASSWORD=your_redis_password
```

#### 3. Run the Project with Docker Compose

```bash
make up
```

This will start the following services:

- **PostgreSQL**: Database service.
- **Redis**: Caching service.
- **Migrate**: Runs database migrations.
- **Server**: The main API server.

#### 4. Access the API

The API will be running on `http://localhost:8000`. You can interact with it using tools like Postman or cURL.

#### 5. View Logs

To view logs for specific services, use the following commands:

```bash
make log-server      # API server logs
make log-postgres    # PostgreSQL logs
make log-redis       # Redis logs
make log-migrate     # Database migration logs
```

#### 6. Stop the Project

To stop the running services, use:

```bash
make down
```

---

## Configuration

### Environment Variables

The project uses environment variables for configuration. You can modify the `.env` file to change database credentials, Redis settings, and other configurations.

### Database and Redis Configuration

If you want to change the database or Redis properties, update the `.env` file and the `config/config.yaml` file accordingly.

---

## API Documentation

The API is documented using Swagger. Once the server is running, you can access the Swagger UI at:

```
http://localhost:8000/docs/index.html
```

---


## Project Structure

```
task-rootext/
├── config/               # Configuration files
├── internal/             # Internal packages (domain, entities, helpers)
├── migrations/           # Database migration files
├── .env-sample           # Sample environment variables
├── docker-compose.yml    # Docker Compose configuration
├── Dockerfile            # Dockerfile for the API server
├── go.mod                # Go module file
├── go.sum                # Go dependencies checksum
└── README.md             # Project documentation
```

---

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

---

## Contact

For any questions or feedback, please contact **Arsham Roshannejad**.

Enjoy using Task-Rootext! 🚀


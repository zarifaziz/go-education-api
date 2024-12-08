# Education API

A RESTful API service built with Go (Golang) for managing educational resources.

## 🚀 Features

- RESTful API endpoints
- MongoDB database integration
- Gin web framework
- Clean architecture pattern
- CLI commands using Cobra

## 📋 Prerequisites

- Go 1.23.4 or higher
- MongoDB
- Git

## 🛠️ Installation

1. Clone the repository

```bash
git clone https://github.com/zarifaziz/go-education-api.git
cd education-api
```

2. Install dependencies

```bash
go mod download
```

3. Set up environment variables (create a `.env` file in the root directory)

```bash
MONGODB_URI=mongodb://localhost:27017
PORT=8080
```

## 📁 Project Structure

```
education-api/
├── api/
│   └── models/      # Data models and structures
├── cmd/            # CLI commands
│   ├── root.go     # Root command setup
│   ├── serve.go    # Server command
│   └── version.go  # Version command
├── server/         # Server implementation
│   ├── server.go   # HTTP server and routes
│   └── db.go       # Database connection
├── main.go         # Application entry point
├── go.mod          # Go module file
├── go.sum          # Go module checksum
└── README.md       # Project documentation
```

## 🏃‍♂️ Running the Application

Start MongoDB:

```bash
mongod --dbpath ~/data/db
```

Available CLI commands:

```bash
# Start the API server
go run main.go serve

# Check version
go run main.go version

# Show available commands
go run main.go --help
```

## 🔄 API Endpoints

```
GET    /api/v1/courses           # List all courses
POST   /api/v1/courses           # Create a course
GET    /api/v1/students/:id      # Get student details
POST   /api/v1/students          # Create a student
POST   /api/v1/students/:id/enroll/:courseId  # Enroll student in course
```

## 🧪 Testing

Run the test suite:

```bash
go test ./...
```

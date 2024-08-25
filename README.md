Simple Golang API
This is a simple RESTful API built with Go (Golang). The API provides basic CRUD operations for managing resources, with a focus on clear and maintainable code.

Table of Contents
Overview
Features
Architecture
Getting Started
Endpoints
Contributing
License
Overview
This project is a straightforward example of a REST API implemented in Go. It demonstrates the use of the standard library along with some popular Go packages for handling HTTP requests, routing, and data management.

Features
RESTful API with standard CRUD operations.
Modular code structure for easy maintenance.
In-memory data storage for simplicity.
Clear separation of concerns with controllers, services, and repository layers.
Basic error handling and logging.
Architecture
The project follows a simple layered architecture to ensure separation of concerns and ease of maintenance:

Models: Define the data structures used in the application.
Controllers: Handle HTTP requests and responses, interacting with the service layer.
Services: Contain the business logic, orchestrating data flows between the controllers and repositories.
Repositories: Handle data storage and retrieval. For this simple API, an in-memory store is used, but this layer can be extended to interact with databases.
Routes: Define the API endpoints and map them to the corresponding controllers.

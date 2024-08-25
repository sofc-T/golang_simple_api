Task Management System API 
Manage your taks remotely
Postman Api documentation ---> https://documenter.getpostman.com/view/32082424/2sA3s3HAzd

the following variables should be defined from the .env file
db_mongo_uri:
db_mongo_name:
jwt_secret_key
	

Features:
Create A Task
See Tasks
Find a Task 
Delete a Task 
Update A Task

folder structure:
main.go: Entry point of the application.
controllers/task_controller.go: Handles incoming HTTP requests and invokes the appropriate service methods.
models/: Defines the data structures used in the application.
data/task_service.go: Contains business logic and data manipulation functions.
router/router.go: Sets up the routes and initializes the Gin router and Defines the routing configuration for the API.
docs/api_documentation.md: Contains API documentation and other related documentation.


api documentation
User Registration
Endpoint: POST /register
Description: Registers a new user.
Request Body:
json
Copy code
{
  "email": "string",
  "password": "string",
  "name": "string",
  "role": "string"
}
Responses:
201 Created: User successfully registered.
400 Bad Request: Invalid input data.
500 Internal Server Error: Server error during registration.
User Login
Endpoint: POST /login
Description: Logs in a user and returns a token.
Request Body:
json
Copy code
{
  "email": "string",
  "password": "string"
}
Responses:
200 OK: Successful login, returns JWT token.
401 Unauthorized: Invalid credentials.
500 Internal Server Error: Server error during login.
Get All Tasks
Endpoint: GET /tasks
Description: Retrieves all tasks. Requires authentication.
Headers:
Authorization: Bearer <token>
Responses:
200 OK: Returns a list of tasks.
401 Unauthorized: Missing or invalid token.
500 Internal Server Error: Server error during task retrieval.
Get Task by ID
Endpoint: GET /tasks/:id
Description: Retrieves a specific task by ID. Requires authentication.
Headers:
Authorization: Bearer <token>
Path Parameters:
id: The ID of the task to retrieve.
Responses:
200 OK: Returns the task details.
401 Unauthorized: Missing or invalid token.
404 Not Found: Task not found.
500 Internal Server Error: Server error during task retrieval.
Create Task (Admin Only)
Endpoint: POST /tasks
Description: Creates a new task. Requires admin privileges.
Headers:
Authorization: Bearer <token>
Request Body:
json
Copy code
{
  "title": "string",
  "description": "string",
  "due_date": "string"
}
Responses:
201 Created: Task successfully created.
401 Unauthorized: Missing or invalid token.
403 Forbidden: Insufficient privileges.
500 Internal Server Error: Server error during task creation.
Update Task (Admin Only)
Endpoint: PUT /tasks/:id
Description: Updates an existing task by ID. Requires admin privileges.
Headers:
Authorization: Bearer <token>
Path Parameters:
id: The ID of the task to update.
Request Body:
json
Copy code
{
  "title": "string",
  "description": "string",
  "due_date": "string"
}
Responses:
200 OK: Task successfully updated.
401 Unauthorized: Missing or invalid token.
403 Forbidden: Insufficient privileges.
404 Not Found: Task not found.
500 Internal Server Error: Server error during task update.
Delete Task (Admin Only)
Endpoint: DELETE /tasks/:id
Description: Deletes a task by ID. Requires admin privileges.
Headers:
Authorization: Bearer <token>
Path Parameters:
id: The ID of the task to delete.
Responses:
200 OK: Task successfully deleted.
401 Unauthorized: Missing or invalid token.
403 Forbidden: Insufficient privileges.
404 Not Found: Task not found.
500 Internal Server Error: Server error during task deletion.
Promote User (Admin Only)
Endpoint: POST /promote
Description: Promotes a user to a higher role. Requires admin privileges.
Headers:
Authorization: Bearer <token>
Request Body:
json
Copy code
{
  "email": "string",
  "new_role": "string"
}
Responses:
200 OK: User successfully promoted.
401 Unauthorized: Missing or invalid token.
403 Forbidden: Insufficient privileges.
404 Not Found: User not found.
500 Internal Server Error: Server error during user promotion.





Database and services documentation 

Variables
db: A pointer to the mongo.Client which holds the MongoDB client instance.
Functions
SetUpDataBase(uri string) (*mongo.Client, error)
Sets up the database connection using the provided URI.

Parameters:

uri (string): The MongoDB URI.
Returns:

(*mongo.Client, error): The MongoDB client instance and an error if any occurred.
Description:

Applies the URI options.
Connects to the MongoDB server.
Sets the global db variable to the connected client.
FetchTasks() ([]models.Task, error)
Fetches all tasks from the database.

Returns:

([]models.Task, error): A slice of tasks and an error if any occurred.
Description:

Retrieves all documents from the tasks collection in the task_manager database.
Decodes the documents into a slice of models.Task.
FindTask(id int) (models.Task, error)
Finds a specific task by its ID.

Parameters:

id (int): The ID of the task to find.
Returns:

(models.Task, error): The task with the specified ID and an error if any occurred.
Description:

Filters the tasks collection for a document with the specified ID.
Decodes the document into a models.Task.
UpdateTask(id int, title string) (models.Task, error)
Updates the title of a task with the specified ID.

Parameters:

id (int): The ID of the task to update.
title (string): The new title for the task.
Returns:

(models.Task, error): The updated task and an error if any occurred.
Description:

Filters the tasks collection for a document with the specified ID.
Updates the title of the task.
Returns the updated task.
DeleteTask(id int) error
Deletes a task with the specified ID.

Parameters:

id (int): The ID of the task to delete.
Returns:

error: An error if any occurred.
Description:

Filters the tasks collection for a document with the specified ID.
Deletes the document.
InsertTask(id int, title string) (models.Task, error)
Inserts a new task with the specified ID and title.

Parameters:

id (int): The ID of the new task.
title (string): The title of the new task.
Returns:

(models.Task, error): The newly inserted task and an error if any occurred.
Description:

Creates a new models.Task with the specified ID and title.
Inserts the task into the tasks collection.
Example Usage
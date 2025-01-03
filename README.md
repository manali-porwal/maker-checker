## Maker-Checker Approval Process Service
This is a simplified service for implementing a maker-checker approval process using Golang Clean architecture, Gin Framework, JWT Authentication, and PostgreSQL. The service enables users to create and approve messages before sending them, with role-based access control (RBAC) distinguishing between makers (who create the message) and checkers (who approve or reject it).

## Logic:
A message will initially be in the `pending` state when created, if the required number of users approve the message (this value is configurable) then message status will be moved to `approved`, if any user rejects the message while it is still in `pending` state then the message status is set to `rejected`.
When the message is `approved`, it is sent to the recipient and marked as `delivered`.

## Features
Maker can create a message.
Checker can approve or reject the message.
Uses JWT Authentication for secure API access.
Role-Based Access Control (RBAC) ensures that only users with the correct role can access certain functionality.
Uses PostgreSQL for storing messages, users, and approvals.
Built with Gin Web Framework for handling HTTP requests.

## Technologies
Golang: Main programming language.
Gin Framework: Web framework for routing and handling HTTP requests.
PostgreSQL: Database for storing messages, users, and approvals.
JWT: JSON Web Token for secure authentication.
Clean Architecture: Separation of concerns into entities, use cases, interfaces, and infrastructure layers.

## Setup and Installation
- Run `go mod download`, to download dependencies.
- Copy `.env.example` to `.env`, change configurations in `.env` file; database, environments, etc
- Install `goose` https://pressly.github.io/goose/installation/ in your local computer
- Set the environment variable `GOOSE_DBSTRING` in your _~/.bash___profile_ or _~/.zshrc_ with the local database connection string. For example

```
$ echo "GOOSE_DBSTRING=\"user=username password=password dbname=database_name sslmode=disable\"" >> ~/.zshrc
```
- Source the _~/.bash___profile_ or _~/.zshrc_

```
$ source ~/.zshrc
```
- Run the command `goose -dir ./database/migrations postgres $GOOSE_DBSTRING up` to migrate the local database to the latest version
- Run `go run ./cmd/api/.` to instantiate a local http server for development on port 8080

## API Documentation
## Create User Account
curl --location 'http://localhost:8080/api/v1/users' \
--header 'Content-Type: application/json' \
--data '{
    "user_name" : "admin",
    "password" : "welcome",
    "role": "maker"
}'

## Login User
curl --location 'http://localhost:8080/api/v1/users/login' \
--header 'Content-Type: application/json' \
--data '{
   "user_name" : "admin",
    "password" : "welcome"
}'

## Create message
curl --location 'http://localhost:8080/api/v1/messages' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT token>' \
--data-raw '{
    "content" : "Hello!!",
    "recipient" : "abc@example.com"
}'

## Approve message
curl --location --request PATCH 'http://localhost:8080/api/v1/messages/{message_id}/approval' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT token>' \
--data '{
    "approved" : true,
    "comment" : "looks good again!"
}'

## Reject message
curl --location --request PATCH 'http://localhost:8080/api/v1/messages/{message_id}/approval' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT token>' \
--data '{
    "approved" : false,
    "comment" : "looks good again!"
}'

## Get message
curl --location 'http://localhost:8080/api/v1/messages/{message_id}'

## Environment Variables: 
Check `.env.development.example` file for reference
You can set the following environment variables in the .env file to customize the application:

DATABASE_WRITE_DSN, DATABASE_READ_DSN: The connection string for PostgreSQL (example: host=localhost user=postgres password=secret dbname=maker_checker port=5432 sslmode=disable)
JWT_SECRET: The secret key for JWT (example: your_secret_key)
JWT_EXPIRY: This value will be in hours used for setting token expiry time (example: 2)
NUM_REQUIRED_APPROVALS: This will be the number of approvals required to change the state of a message from `pending` -> `approved`

## Testing the Application
You can test the application using any HTTP client like Postman or curl.

Steps:
Create User: Create user accounts with correct roles.
Login: First, log in and get the JWT token.
Create a Message: Use the token to create a message as a maker or admin.
Approve/Reject the Message: As a checker or admin, approve or reject the message.
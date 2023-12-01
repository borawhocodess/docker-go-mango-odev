# docker-go-mango-odev

## How 2 Run
**Clone the Repository**:
```bash
git clone https://github.com/borawhocodess/docker-go-mango-odev
```

**Start the Application**:  Navigate to the project directory and run:
```bash
docker-compose up --build
```

This command will start the Go API server and the MongoDB database.

## Features
- **Add Users**: Users can be added to the MongoDB database with unique usernames and IDs.
- **Get Users**: Retrieve a list of all users stored in the database.
- **Delete Users**: Remove users from the database based on their username.

## API Endpoints
- `POST /addUserToDB`: Adds a new user to the database.
- `GET /getUsersFromDB`: Retrieves all users from the database.
- `DELETE /deleteUser/:username`: Deletes a user from the database based on the username.

## Testing
**Add User**:
  ```bash
  curl -X POST http://localhost/addUserToDB -H "Content-Type: application/json" -d '{"username": "testuser", "id": 1}'
  ```
**Get User**:
  ```bash
  curl -X GET http://localhost/getUsersFromDB
  ```
**Delete User**:
  ```bash
  curl -X DELETE http://localhost/deleteUser/testuser
  ```



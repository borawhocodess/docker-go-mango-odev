package main

import (
	"net/http"
	"os"
	"log"
    "context"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
}

var client *mongo.Client

func main() {

	// Set MongoDB client options
    clientOptions := options.Client().ApplyURI("mongodb://root:example@mongo:27017/")

    // Correct way to initialize the global client variable
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

    log.Println("\n\n\n\nConnected to MongoDB!\n\n\n\n")


	r := gin.Default()
	baseURL := os.Getenv("baseURL")

	// r.GET(baseURL+"/", handleUserEvent)
	r.GET(baseURL+"/ping", handlePing)
	r.GET(baseURL+"/pong", handlePong)

	r.GET(baseURL+"/getUsers", getUsers)
	r.GET(baseURL+"/getUsersFromDB", getUsersFromDB)

	r.POST(baseURL+"/addUserToDB", addUserToDB)

	r.DELETE(baseURL+"/deleteUser/:username", deleteUserFromDB)

	// r.PUT(baseURL+"/updateUser", updateUserInDB)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func handlePong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ping",
	})
}

func getUsers(c *gin.Context) {
	var users []User

	u1 := User{
		Username: "userA",
		Id:       1,
	}

	u2 := User{
		Username: "userB",
		Id:       2,
	}

	users = append(users, u1, u2)

	c.JSON(http.StatusOK, users)
}

func addUserToDB(c *gin.Context) {
    var newUser User
    if err := c.BindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    collection := client.Database("mydatabase").Collection("users")

    // Check if a user with the same username or ID already exists
    var existingUser User
    err := collection.FindOne(context.TODO(), bson.M{
        "$or": []bson.M{
            {"username": newUser.Username},
            {"id": newUser.Id},
        },
    }).Decode(&existingUser)

    if err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "User with the same username or ID already exists"})
        return
    } else if err != mongo.ErrNoDocuments {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Insert the new user
    _, err = collection.InsertOne(context.TODO(), newUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    log.Printf("User added successfully: Username: %s, ID: %d\n", newUser.Username, newUser.Id)
    c.JSON(http.StatusOK, gin.H{"message": "User added successfully"})
}


func getUsersFromDB(c *gin.Context) {
    var users []User

    collection := client.Database("mydatabase").Collection("users")

    // Finding multiple documents returns a cursor
    cursor, err := collection.Find(context.TODO(), bson.D{{}})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	log.Printf("\n\n\n\n")

    // Iterate through the returned cursor.
    for cursor.Next(context.TODO()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)

		// Print each user to the terminal
		log.Printf("Read user from DB: Username: %s, ID: %d\n", user.Username, user.Id)
	}

	log.Printf("\n\n\n\n")

    if err := cursor.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Close the cursor once finished
    cursor.Close(context.TODO())

    c.JSON(http.StatusOK, users)
}

func deleteUserFromDB(c *gin.Context) {
    // Get the username from the URL parameter
    username := c.Param("username")

    collection := client.Database("mydatabase").Collection("users")

    // Create a filter to match the user by username
    filter := bson.M{"username": username}

    // Delete the user from the database
    _, err := collection.DeleteOne(context.TODO(), filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Print to the terminal
    log.Printf("User deleted successfully: Username: %s\n", username)

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

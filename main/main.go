package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NikolaSmrcek/TableCalc/routes"

	"gopkg.in/redis.v5"
)

// redisConnectionAddress used as "configuration" for redis connection.
var redisConnectionAddress string = "localhost:6379";

//Redis connection object pointer.
var redisClient *redis.Client;

// Setting up redis connection object and checking if it's valid connection.
// Returns: redis connection object pointer if valid connection or nil if invalid.
func connectToRedis()(*redis.Client){
	client := redis.NewClient(&redis.Options{
		Addr:     redisConnectionAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	});
	pong, err := client.Ping().Result()
	if err != nil{
		fmt.Println("Redis connection error: ", err);
		fmt.Println("Closing server...");
		return nil;
	} else{
		fmt.Println("Redis connection established.", pong);
		return client;
	}

}

func main() {
	fmt.Println("Starting server...");
	redisClient = connectToRedis()
	// Testing if the connection is alive
	if redisClient == nil {
		return;
	}

	routes.RedisClient = redisClient;
	router := routes.NewRouter();

	log.Fatal(http.ListenAndServe(":8080", router))

}

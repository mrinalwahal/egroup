package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/mrinalwahal/egroup/handler"
)

func main() {

	//	We can also read this port value from command line arguments.
	//	But for ease, let's hardcode it.
	const port = ":8080"

	//	Read the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to read .env file: %s", err)
	}

	//	Initialize a new router
	router := http.NewServeMux()

	//	Register our handler function on a random route
	router.HandleFunc("/api", handler.Query)

	//	Start the server
	log.Println("Server is now listening on ", port)
	log.Fatal(http.ListenAndServe(port, router))
}

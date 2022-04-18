package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/machinebox/graphql"
)

func Query(w http.ResponseWriter, r *http.Request) {

	//	It's good practice to protect your endpoints with some credentials.
	//	To enable this check, simply uncomment the following code.
	/*
			if r.Header.Get("webhook-secret") != os.Getenv("WEBHOOK_SECRET") {
		   		http.Error(w, "invalid webhook secret", http.StatusUnauthorized)
		   		return
		   	}
	*/

	//	Unmarshal the request payload into structure
	var payload Request
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {

		//	If no request body has been supplied in the request,
		//	we may encounter an EOF error.
		//	Therefore, we must ignore it.
		if !errors.Is(err, io.EOF) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	//	Initialize GraphQL Client
	var gqlClient = graphql.NewClient(os.Getenv("GRAPHQL_BACKEND_URL"))

	// Prepare the GraphQL request
	gqlReq := graphql.NewRequest(`
	query last_projects($n: Int = 5) {
		projects(last:$n) {
		  nodes {
			name
			description
			forksCount
		  }
		}
	  }`)

	//	If the GraphQL endpoint is protected and requires an authorization header,
	//	set that header to access the endpoint
	//	gqlReq.Header.Set("Authorization:", "Bearer " + os.Getenv("ACCESS_TOKEN"))

	// Set the variables
	if payload.N > 0 {
		gqlReq.Var("n", payload.N)
	}

	var resp GQLResponse

	//	Run the query
	if err := gqlClient.Run(context.Background(), gqlReq, &resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var names []string
	var forkSum int

	//	Calculate the response
	for _, item := range resp.Projects.Nodes {
		names = append(names, item.Name)
		forkSum += item.ForksCount
	}

	//	Perpare the response JSON
	response, err := json.Marshal(Response{
		Names: strings.Join(names, ", "),
		Sum:   forkSum,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//	Send the JSON in response
	w.WriteHeader(200)
	w.Write(response)
}

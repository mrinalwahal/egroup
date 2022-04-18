package handler

type (

	//	HTTP Request payload received from client
	Request struct {
		N int `json:"n,omitempty"`
	}

	//	GitLab GraphQL response returned
	GQLResponse struct {
		Projects struct {
			Nodes []Node `json:"nodes,omitempty"`
		} `json:"projects,omitempty"`
	}

	//	Structure of every node
	Node struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		ForksCount  int    `json:"forksCount,omitempty"`
	}

	//	Response structure to be returned by our handler function
	Response struct {
		Names string `json:"names"`
		Sum   int    `json:"sum"`
	}
)

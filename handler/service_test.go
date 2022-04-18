package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
)

func Test_Query(t *testing.T) {

	//	Read the .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Failed to read .env file: %s", err)
	}

	const url = "http://localhost:8080/api"
	body, _ := json.Marshal(Request{
		N: 10,
	})

	type args struct {
		w httptest.ResponseRecorder
		r *http.Request
	}

	tests := []struct {
		name       string
		args       args
		statusCode int
		result     Response
	}{
		//	First mandatory test without any request body.
		//	Which means, the default N value is used in GraphQL request.
		{
			name: "request w/o payload",
			args: args{
				w: *httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, url, nil),
			},
			statusCode: 200,
			result: Response{
				Sum: 13,

				//	You can also add name strings over here
				//	to have them matched.
				//	Names: "",
			},
		},

		//	Second mandatory test with a custom request body.
		//	Which means, the supplied N value is used in GraphQL request.
		{
			name: "request w/ payload",
			args: args{
				w: *httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(body)),
			},
			statusCode: 200,
			result: Response{
				Sum: 19,

				//	You can also add name strings over here
				//	to have them matched.
				//	Names: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Query(&tt.args.w, tt.args.r)
		})

		//	Read the HTTP response recorded
		res := tt.args.w.Result()
		defer res.Body.Close()

		//	Match the response status code
		if res.StatusCode != tt.statusCode {
			t.Errorf("statusCode => expected %v; got %v", tt.statusCode, res.StatusCode)
		}

		//	Read the HTTP response body
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		//	Unmarshal the response
		var result Response
		json.Unmarshal(data, &result)

		//	Match the forkCountSum
		if result.Sum != tt.result.Sum {
			t.Errorf("sum => expected %v; got %v", tt.result, result.Sum)
		}
	}
}

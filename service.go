package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	api "github.com/mesg-foundation/core/api/service"
	"google.golang.org/grpc"
)

// Service for MESG
type Service struct {
	client   api.ServiceClient
	endpoint string
	token    string
}

// OnRequestEvent holds info about request
type OnRequestEvent struct {
	Date time.Time `json:"date"`
	ID   string    `json:"id"`
	Body string    `json:"body"`
}

// Request to service
type Request struct {
	ID   int    `json:"id,omitempty"`
	URL  string `json:"url"`
	Body string `json:"body"`
}

// Response from server struct
type Response struct {
	ID      int              `json:"id,omitempty"`
	Success *SuccessResponse `json:"success,omitempty"`
	Error   *ErrorResponse   `json:"error,omitempty"`
}

// SuccessResponse contains custom success types
type SuccessResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

// ErrorResponse contains custom error type
type ErrorResponse struct {
	Message string `json:"message"`
}

// NewService returns new MSG service
func NewService(endpoint string, token string) (*Service, error) {
	connection, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Service{
		client:   api.NewServiceClient(connection),
		endpoint: endpoint,
		token:    token,
	}, nil
}

// OnRequest method for handling event
func (s *Service) OnRequest(data string) error {
	id, err := GenerateID()
	if err != nil {
		return err
	}
	buff, err := json.Marshal(OnRequestEvent{
		Date: time.Now(),
		ID:   id,
		Body: data,
	})

	_, err = s.client.EmitEvent(context.Background(), &api.EmitEventRequest{
		Token:     s.token,
		EventKey:  "onRequest",
		EventData: string(buff),
	})
	return err
}

// ListenTask handles incoming MESG task
func (s *Service) ListenTask() error {
	stream, err := s.client.ListenTask(context.Background(), &api.ListenTaskRequest{
		Token: os.Getenv("MESG_TOKEN"),
	})
	if err != nil {
		return err
	}

	for {
		var err error
		res, err := stream.Recv()
		if err != nil {
			log.Fatalln(err)
			continue
		}

		var outputKey string
		var responseBuffer []byte

		switch res.TaskKey {
		case "execute":
			outputKey = "response"
			req := new(Request)
			err = json.Unmarshal([]byte(res.InputData), req)
			if err != nil {
				break
			}

			responseBuffer, err = json.Marshal(s.handleTask(req, nil))
			break

		case "batchExecute":
			outputKey = "batchResponse"
			var requests []Request
			var responses []Response
			wg := new(sync.WaitGroup)

			if err := json.Unmarshal([]byte(res.InputData), &requests); err != nil {
				break
			}

			for _, r := range requests {
				wg.Add(1)
				go func(wg *sync.WaitGroup, r Request) {
					response := s.handleTask(&r, wg)
					response.ID = r.ID
					responses = append(responses, *response)
				}(wg, r)
			}

			wg.Wait()
			responseBuffer, err = json.Marshal(responses)
			break
		}

		if err != nil {
			log.Fatalln(err)
			responseBuffer, _ = json.Marshal(makeResponse(0, "", err))
		}

		_, err = s.client.SubmitResult(context.Background(), &api.SubmitResultRequest{
			ExecutionID: res.ExecutionID,
			OutputKey:   outputKey,
			OutputData:  string(responseBuffer),
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (s *Service) handleTask(req *Request, wg *sync.WaitGroup) *Response {
	if wg != nil {
		defer wg.Done()
	}

	// validate only URL, body can be empty
	if len(req.URL) == 0 {
		return makeResponse(0, "", errors.New("Missing URL parametar"))
	}

	resp, err := http.Post(req.URL, "application/json", strings.NewReader(req.Body))
	if err != nil {
		return makeResponse(0, "", err)
	}

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return makeResponse(0, "", err)
	}
	return makeResponse(resp.StatusCode, string(buff), nil)
}

func makeResponse(statusCode int, body string, err error) *Response {
	if err != nil {
		return &Response{
			Error: &ErrorResponse{
				Message: err.Error(),
			},
		}
	}
	return &Response{
		Success: &SuccessResponse{
			StatusCode: statusCode,
			Body:       body,
		},
	}
}

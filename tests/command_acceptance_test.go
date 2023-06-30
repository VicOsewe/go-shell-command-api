package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/VicOsewe/go-shell-command-api/configs"
	"github.com/VicOsewe/go-shell-command-api/presentation"
	"github.com/imroc/req"
)

var srv *http.Server
var baseURL string
var serverErr error

func randomPort() int {
	rand.Seed(time.Now().Unix())
	min := 32768
	max := 60999
	port := rand.Intn(max-min+1) + min
	return port
}

func startTestServer() (*http.Server, string, error) {
	port := randomPort()
	srv := presentation.PrepareServer(port)
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if srv == nil {
		return nil, "", fmt.Errorf("nil test server")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, "", fmt.Errorf(
			"unable to listen on port %d: %w",
			port,
			err,
		)
	}
	if listener == nil {
		return nil, "", fmt.Errorf("nil test server listener")
	}

	log.Printf("LISTENING on port %d", port)

	go func() {
		err := srv.Serve(listener)
		if err != nil {
			log.Printf("serve error: %s", err)
		}
	}()

	return srv, baseURL, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	srv, baseURL, serverErr = startTestServer()
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
	}

	code := m.Run()

	defer func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("test server shutdown error: %s", err)
		}
	}()
	os.Exit(code)
}

func TestCMDHandler(t *testing.T) {
	client := http.Client{
		Timeout: time.Minute * 10,
	}

	var input struct {
		Command string `json:"command"`
	}

	input.Command = "pwd"

	inputBytes, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("can't marshal message to JSON: %s", err)
	}

	headers := req.Header{
		"Content-Type": "application/json",
	}

	badInputBytes, err := json.Marshal("bad-input")
	if err != nil {
		t.Fatalf("can't marshal message to JSON: %s", err)
	}

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
		headers    map[string]string
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "happy case - json body",
			args: args{
				url:        fmt.Sprintf("%s/api/cmd", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(inputBytes),
				headers:    headers,
			},
			wantStatus: 200,
		},
		{
			name: "happy case - query params",
			args: args{
				url:        fmt.Sprintf("%s/api/cmd?command=ls", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(inputBytes),
				headers:    headers,
			},
			wantStatus: 200,
		},
		{
			name: "sad case - bad route",
			args: args{
				url:        fmt.Sprintf("%s/not-cmd", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(inputBytes),
				headers:    headers,
			},
			wantStatus: 404,
		},
		{
			name: "sad case - empty command",
			args: args{
				url:        fmt.Sprintf("%s/api/cmd", baseURL),
				httpMethod: http.MethodPost,
				headers:    headers,
			},
			wantStatus: 400,
		},
		{
			name: "sad case - bad request",
			args: args{
				url:        fmt.Sprintf("%s/api/cmd", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(badInputBytes),
				headers:    headers,
			},
			wantStatus: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			r.SetBasicAuth(
				configs.MustGetEnvVar("AUTH_USERNAME"),
				configs.MustGetEnvVar("AUTH_PASSWORD"),
			)

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("cannot read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf(
					"expected status %d, got %d and response %s",
					tt.wantStatus,
					resp.StatusCode,
					string(data),
				)
				return
			}
		})
	}
}

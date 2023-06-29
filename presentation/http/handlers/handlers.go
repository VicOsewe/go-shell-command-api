package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VicOsewe/go-shell-command-api/configs"

	"github.com/VicOsewe/go-shell-command-api/internal/app/entities"
	"github.com/VicOsewe/go-shell-command-api/internal/app/usecases"
)

// Presentation represents the presentation layer contract
type Presentation interface {
}

// RestFulAPIs set up RESTFUL APIs with all necessary dependencies
type RestFulAPIs struct {
	auth struct {
		username string
		password string
	}
}

func NewRestFulAPIs() *RestFulAPIs {
	rst := &RestFulAPIs{}
	rst.auth.password = configs.MustGetEnvVar("AUTH_PASSWORD")
	rst.auth.username = configs.MustGetEnvVar("AUTH_USERNAME")
	rst.checkPreconditions()
	return rst
}

func (rst *RestFulAPIs) checkPreconditions() {
	if rst.auth.username == "" {
		log.Panicf("error, basic auth password must be provided")
	}
	if rst.auth.password == "" {
		log.Panicf("error, basic auth password must be provided")
	}
}

func (rst *RestFulAPIs) CMDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := entities.CommandValue{}
		DecodeJSONTargetStruct(w, r, &input)
		if input.Command == "" {
			http.Error(w, "Missing 'command' query parameter", http.StatusBadRequest)
			return
		}
		output, err := usecases.ExecuteCommand(input.Command)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", output)
	}
}

package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/VicOsewe/go-shell-command-api/configs"

	"github.com/VicOsewe/go-shell-command-api/entities"
	"github.com/VicOsewe/go-shell-command-api/usecases"
)

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
		command := r.FormValue("command")

		if command == "" {
			var body struct {
				Command string `json:"command"`
			}
			DecodeJSONTargetStruct(w, r, &body)
			command = body.Command
		}

		if command == "" {
			payload := entities.APIResponseMessage{
				Message:    "invalid request data, ensure `command` is provided",
				StatusCode: http.StatusBadRequest,
			}
			RespondWithJSON(w, http.StatusBadRequest, payload)
			return
		}

		output, err := usecases.ExecuteCommand(command)
		if err != nil {
			payload := entities.APIResponseMessage{
				Message:    "Command execution failed",
				StatusCode: http.StatusBadRequest,
			}
			RespondWithJSON(w, http.StatusBadRequest, payload)
			return
		}

		payload := entities.APIResponseMessage{
			Message:    "command retrieved successfully",
			StatusCode: http.StatusOK,
			Body:       strings.ReplaceAll(output, "\n", " "),
		}
		RespondWithJSON(w, http.StatusOK, payload)
	}
}

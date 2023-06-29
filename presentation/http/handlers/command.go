package handlers

import (
	"fmt"
	"net/http"

	"github.com/VicOsewe/go-shell-command-api/internal/app/entities"
	"github.com/VicOsewe/go-shell-command-api/internal/app/usecases"
)

func CMDHandler() http.HandlerFunc {
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

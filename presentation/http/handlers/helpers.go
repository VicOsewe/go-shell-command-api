package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func DecodeJSONTargetStruct(w http.ResponseWriter, r *http.Request, targetStruct interface{}) {
	err := json.NewDecoder(r.Body).Decode(targetStruct)
	if err != nil {
		WriteJSONResponse(w, ErrorMap(err), http.StatusBadRequest)
		return
	}
}

// ErrorMap turns the supplied error into a map with "error" as the key
func ErrorMap(err error) map[string]string {
	errMap := make(map[string]string)
	errMap["error"] = err.Error()
	return errMap
}

// WriteJSONResponse writes the content supplied via the `source` parameter to
// the supplied http ResponseWriter. The response is returned with the indicated
// status.

func WriteJSONResponse(w http.ResponseWriter, source interface{}, status int) {
	w.WriteHeader(status)
	content, errMap := json.Marshal(source)
	if errMap != nil {
		msg := fmt.Sprintf("error when marshalling %#v to JSON bytes: %#v", source, errMap)
		http.Error(w, msg, http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	_, errMap = w.Write(content)
	if errMap != nil {
		msg := fmt.Sprintf(
			"error when writing JSON %s to http.ResponseWriter: %#v", string(content), errMap)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

// RespondWithError writes an error response
func RespondWithError(w http.ResponseWriter, code int, err error) {
	errMap := ErrorMap(err)
	errBytes, err := json.Marshal(errMap)
	if err != nil {
		errBytes = []byte(fmt.Sprintf("error: %s", err))
	}
	RespondWithJSON(w, code, errBytes)
}

// RespondWithJSON writes a JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(marshalledPayload)
	if err != nil {
		log.Printf(
			"unable to write payload `%s` to the http.ResponseWriter: %s",
			string(marshalledPayload),
			err,
		)
	}
}

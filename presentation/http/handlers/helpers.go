package handlers

import (
	"encoding/json"
	"fmt"
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

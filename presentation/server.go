package presentation

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	rest "github.com/VicOsewe/go-shell-command-api/presentation/http/handlers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	serverTimeoutSeconds = 120
)

// Router sets up the ginContext router
func Router() (*mux.Router, error) {
	r := mux.NewRouter()
	h := InitHandlers()

	r.Path("/health").HandlerFunc(HealthStatusCheck)

	RESTRoutes := r.PathPrefix("/api").Subrouter()
	RESTRoutes.Use(h.BasicAuth())

	RESTRoutes.Path("/cmd").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(h.CMDHandler())
	return r, nil
}

// HealthStatusCheck checks if the server is working
func HealthStatusCheck(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(true)
	if err != nil {
		log.Fatal(err)
	}
}

// InitHandlers initializes all the handlers dependencies
func InitHandlers() *rest.RestFulAPIs {
	return rest.NewRestFulAPIs()
}

// PrepareServer prepares the http server
func PrepareServer(port int) *http.Server {
	r, err := Router()
	if err != nil {
		log.Fatalln("There's an error with the server,", err)
	}

	addr := fmt.Sprintf(":%d", port)
	h := handlers.CompressHandlerLevel(r, gzip.BestCompression)

	h = handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST"}),
	)(h)

	h = handlers.CombinedLoggingHandler(os.Stdout, h)

	return &http.Server{
		Handler:      h,
		Addr:         addr,
		WriteTimeout: serverTimeoutSeconds * time.Second,
		ReadTimeout:  serverTimeoutSeconds * time.Second,
	}
}

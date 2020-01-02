package tomolink

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Limit the size of incoming requests to something sensible
var readLimit int64 = 500

var (
	// Logrus structured logging setup
	tlLog *log.Entry

//	cfg   interface{}
)

// Router instantiates a new gorilla mux router and adds the various routes and
// handlers.  This is exported as it is also instantiated by the tests in
// tomolink_test.go
func Router() *mux.Router {
	router := mux.NewRouter()

	// You can find the handlers in handlers.go
	router.HandleFunc("/users/{UUID}", CreateEndpoint).
		Methods("GET", "PUT").
		Name("retrieve")
	router.HandleFunc("/create2", CreateEndpoint2).Methods("GET").Headers("Content-Type", "application/json").Name("Test2")

	return router
}

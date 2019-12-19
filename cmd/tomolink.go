package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Limit the size of incoming requests to something sensible
var readLimit int64 = 500

// Router instantiates a new gorilla mux router and adds the various routes and
// handlers.  This is exported as it is also instantiated by the tests in
// tomolink_test.go
func Router() *mux.Router {
	router := mux.NewRouter()

	// You can find the handlers in handlers.go
	//
	router.HandleFunc("/users/{UUID}", CreateEndpoint).
		Methods("GET", "PUT").
		Name("retrieve")
	router.HandleFunc("/create2", CreateEndpoint2).Methods("GET").Headers("Content-Type", "application/json").Name("Test2")

	return router
}

// Player is TODO
type Player struct {
	// TODO
	Rels string
}

// CreateEndpoint2 ++
func CreateEndpoint2(w http.ResponseWriter, r *http.Request) {
	var p map[string]string
	err := decodeJSONBody(w, r, &p)
	if err != nil {
		log.Fatal(err.Error())
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	t, err := json.Marshal(p)
	io.WriteString(w, string(t))
	log.Info(string(t))

}

// CreateEndpoint is WIP
func CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	// Check that the content header (if set) is application/json
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		log.Info(value)
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	//body, err := ioutil.ReadAll(io.LimitReader(r.Body, readLimit))
	if r.Body != nil {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, readLimit))
		//body, err := ioutil.ReadAll(r.Body)
		_ = body
		_ = err
		// https: //stackoverflow.com/questions/32710847/what-is-the-best-way-to-check-for-empty-request-body
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		if len(body) > 0 {
			io.WriteString(w, string(body))
			return
		}
	}
	io.WriteString(w, `{"body": "nope"}`)
}

func main() {
	router := Router()

	log.Info("Starting HTTP server")
	log.Fatal(http.ListenAndServe(":8080", router))
}

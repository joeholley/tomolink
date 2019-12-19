package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	testBody map[string]string
	router   *mux.Router
)

func TestMain(m *testing.M) {
	// Set up tests
	router = Router()
	testBody = map[string]string{"this": "yep"}

	flag.Parse()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestWalk(t *testing.T) {

	// Walk through every route added to the router in tomolink.go, calling the
	// handler for each route and testing for expected output.
	err := router.Walk(func(route *mux.Route, rtr *mux.Router, ancestors []*mux.Route) error {

		routeLog := log.WithFields(log.Fields{"name": route.GetName()})

		// Get URL for this route
		URL, err := route.URL("UUID", "jojo")
		if err != nil {
			routeLog.Fatal(err)
		}
		routeLog = routeLog.WithFields(log.Fields{"URL": URL})

		// Get HTTP methods (verbs) supported by this route
		methods, err := route.GetMethods()
		if err != nil {
			routeLog.Fatal(err)
		}

		// Fake a body
		requestBody, err := json.Marshal(testBody)
		if err != nil {
			routeLog.Fatal(err)
		}

		// make a call
		for _, method := range methods {
			routeLog = routeLog.WithFields(log.Fields{"method": method})
			routeLog.Info("Testing route handler")

			request, err := http.NewRequest(method, URL.String(), bytes.NewBuffer(requestBody))
			request.Header.Set("Content-Type", "application/json")
			if err != nil {
				routeLog.Fatal(err)
			}

			response := httptest.NewRecorder()
			rtr.ServeHTTP(response, request)
			assert.Equal(t, 200, response.Code, "HTTP OK response expected")
			//TODO: Moar testss
			routeLog.Info("Tests complete!")
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

//func TestCreateEndpoint(t *testing.T) {
//	ts := map[string]string{"this": "yep"}
//	requestBody, err := json.Marshal(ts)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	_ = requestBody
//	request, err := http.NewRequest("GET", "/create", bytes.NewBuffer(requestBody))
//	request.Header.Set("Content-Type", "application/json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	response := httptest.NewRecorder()
//	router.ServeHTTP(response, request)
//
//	responseBody, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Tests
//	assert.Equal(t, 200, response.Code, "HTTP OK response expected")
//	assert.Equal(t,
//		"application/json",
//		response.Header().Get("Content-Type"), // https://groups.google.com/g/golang-nuts/c/qRJOtFV1J2g
//		"Response header Content-Type value of 'application/json' expected")
//	assert.Equal(t, string(requestBody), string(responseBody), "Unexpected HTTP response body contents")
//}

// Read config, do a test based on maximum readLimit size?

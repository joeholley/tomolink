// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains the application initialization code for the Tomolink service.
// Actual code for the service itself can be found in ../internal/app/tomolink
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joeholley/tomolink/internal/app/tomolink"
	"github.com/joeholley/tomolink/internal/config"
	"github.com/joeholley/tomolink/internal/logging"
	"github.com/sirupsen/logrus"
)

var (
	// Fields to add to the structured logs
	tlLog = logrus.WithFields(logrus.Fields{
		"app":       "tomolink",
		"component": "app.main",
	})
)

func main() {
	// Read config file, file name is
	//   <string_argument>_defaults.cfg
	// so in this case it will be called
	//   tomolink_defaults.cfg
	ac := config.AppConfig{}
	err := ac.Read("tomolink")
	if err != nil {
		tlLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("Cannot load configuration")
	}

	// set up logrus structured logging
	format, err := ac.Cfg.StringOr("logging.format", "text")
	if err != nil {
		tlLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("Cannot load logging format configuration")
	}
	lvl, err := ac.Cfg.StringOr("logging.level", "info")
	if err != nil {
		tlLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("Cannot load logging level configuration")
	}
	logging.ConfigureLogging(format, lvl)

	// Instantiate router
	router := tomolink.Router(&ac)

	// Get port from config
	port, err := ac.Cfg.StringOr("http.port", "8080")
	if err != nil {
		tlLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Unable to read http port from config; defaulting to serving on port 8080")
	}

	// Start server.  Largely this is using the example code from https://github.com/gorilla/mux
	tlLog.Info("Starting HTTP server")
	//tlLog.Fatal(http.ListenAndServe(":"+port, router))
	srv := &http.Server{
		Addr: "0.0.0.0:" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			tlLog.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	deadline, err := ac.Cfg.IntOr("http.gracefulwait", 15)
	if err != nil {
		tlLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Unable to read http graceful wait duration from config; defaulting to 15 seconds")
	}
	wait := time.Duration(deadline) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	tlLog.Println("shutting down")
	os.Exit(0)
}

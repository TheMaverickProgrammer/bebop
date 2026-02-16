package main

import (
	"flag"
	"os"

	"github.com/disintegration/bebop/store"
)

// addUser creates a new user from username and password.
func addUser() {
	username := flag.Arg(1)
	if username == "" {
		help()
		os.Exit(2)
	}

	password := flag.Arg(2)
	if password == "" {
		help()
		os.Exit(2)
	}

	cfg, err := getConfig()
	if err != nil {
		logger.Fatalf("failed to load configuration: %s", err)
	}

	s, err := getStore(cfg)
	if err != nil {
		logger.Fatalf("failed to get data store: %s", err)
	}

	_, err = s.Users().GetByName(username)
	if err != nil {
		// Expect the username to be unreserved. Other errors are a problem.
		if err != store.ErrNotFound {
			logger.Fatalf("user search by username failed: %s", err)
		}
	}

	_, err = s.Users().NewLocal(username, password)

	if err != nil {
		logger.Fatalf("failed to create local user: %s", err)
	}
}

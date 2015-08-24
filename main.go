package main

import (
	"github.com/tyler-sommer/shotgun/config"

	"flag"

	"fmt"

	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

var databaseFile = flag.String("databaseFile", "shotgun.db", "Database storage file")
var privKeyFile = flag.String("privateKeyFile", "~/.ssh/id_rsa", "Private SSH key for authentication")

func main() {
	dbFile, err := tilde.Expand(*databaseFile)
	if err != nil {
		panic(err)
	}

	keyFile, err := tilde.Expand(*privKeyFile)
	if err != nil {
		panic(err)
	}

	conf, err := config.New(dbFile, keyFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(conf)
}

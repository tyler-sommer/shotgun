package main

import (
	"github.com/tyler-sommer/shotgun/config"
	"github.com/tyler-sommer/shotgun/model"
	"github.com/tyler-sommer/shotgun/data"
	"github.com/tyler-sommer/shotgun/executer"

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

	dbm, err := conf.NewDatabaseManager()
	if err != nil {
		panic(err)
	}

	repo, err := dbm.NewServerRepository()
	if err != nil {
		panic(err)
	}

	key := "somekey"

	s1, err := repo.Find(key)
	if err != nil && err != data.KeyNotFoundError {
		panic(err)
	}

	if err == data.KeyNotFoundError {
		fmt.Println("Creating a server")
		script := model.NewScript()
		script.Enabled = true
		script.RequiresSudo = true
		script.Commands = append(script.Commands, "service httpd stop")
		script.Commands = append(script.Commands, "service application stop")

		s1 = model.NewServer()
		s1.SetKey(key)
		s1.Host = "server01.localhost"
		s1.User = "app"
		s1.Scripts = append(s1.Scripts, script)

		err = repo.Save(s1)
		if err != nil {
			panic(err)
		}

		err = dbm.Commit()
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(s1)

	servers, err := repo.All()
	if err != nil {
		panic(err)
	}

	e := executer.New(conf, servers)
	e.Execute()

	for {

	}
}

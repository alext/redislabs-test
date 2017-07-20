package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/alext/redislabs-test/client"
	"github.com/kr/pretty"
)

func main() {
	email := flag.String("email", "", "")
	password := flag.String("password", "", "")
	flag.Parse()

	apiClient, err := client.New(*email, *password)
	if err != nil {
		log.Fatal(err)
	}

	db, err := apiClient.ProvisionDB("testdb-" + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		log.Fatal(err)
	}
	pretty.Println(db)
	err = ioutil.WriteFile("certs/"+db.Name+"_cert.pem", db.Cert, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("certs/"+db.Name+"_key.pem", db.Key, 0600)
	if err != nil {
		log.Fatal(err)
	}

	dbData, err := apiClient.ListDBs()
	if err != nil {
		log.Fatal(err)
	}
	pretty.Println(dbData)
}

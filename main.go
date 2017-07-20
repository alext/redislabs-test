package main

import (
	"flag"
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

	dbData, err := apiClient.ListDBs()
	if err != nil {
		log.Fatal(err)
	}
	pretty.Println(dbData)
}

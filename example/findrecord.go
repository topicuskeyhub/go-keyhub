package main

import (
	"encoding/json"
	"flag"
	"github.com/google/uuid"
	"github.com/topicuskeyhub/go-keyhub"
	"github.com/topicuskeyhub/go-keyhub/model"
	"log"
	"strconv"
)

func main() {

	// variables declaration
	var issuer string
	var clientid string
	var clientsecret string
	var record string
	var err error

	// flags declaration using flag package
	flag.StringVar(&issuer, "i", "https://test.topicus-keyhub.com", "Specify issuer")
	flag.StringVar(&clientid, "ci", "", "Specify client id")
	flag.StringVar(&clientsecret, "cs", "", "Specify client secret")
	flag.StringVar(&record, "r", "", "Record to find")

	flag.Parse()

	client, err := keyhub.NewClientDefault(issuer, clientid, clientsecret)

	if err != nil {
		log.Fatalf("ERROR %s", err)
	}

	additional := &model.VaultRecordAdditionalQueryParams{Secret: true, Audit: true}

	recordUUID, errUUID := uuid.Parse(record)
	recordID, errID := strconv.ParseInt(record, 10, 64)
	if errUUID != nil && errID != nil {
		log.Fatalln("-r is not a valid uuid or id")
	}

	var vaultRecord *model.VaultRecord

	if errUUID == nil {
		vaultRecord, err = client.Vaults.FindByUUIDForClient(recordUUID, additional)
	} else {
		vaultRecord, err = client.Vaults.FindByIDForClient(recordID, additional)
	}

	if err != nil {
		log.Fatalln("vaults.FindByXXXXForClient", err.Error())
	}

	out, err := json.MarshalIndent(vaultRecord, "", "  ")
	if err != nil {
		log.Fatalln("Json marshal error: ", err.Error())
	}
	log.Printf("vault record found: \n%s\n\n", out)

}

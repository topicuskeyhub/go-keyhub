package main

import (
	"encoding/json"
	"flag"
	"github.com/google/uuid"
	"github.com/topicuskeyhub/go-keyhub"
	"github.com/topicuskeyhub/go-keyhub/model"
	"log"
	"math/big"
)

func main() {

	// variables declaration
	var issuer string
	var clientid string
	var clientsecret string
	var record string

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

	recordUUID, err := uuid.Parse(record)

	var vaultRecord *model.VaultRecord

	if err == nil {
		vaultRecord, err = client.Vaults.FindByUUIDForClient(recordUUID, additional)
	} else {
		recordid := big.Int{}
		recordid.SetString(record, 10)
		vaultRecord, err = client.Vaults.FindByIDForClient(recordid.Int64(), additional)
	}

	if err != nil {
		log.Fatalln("vaults.FindByXXXXForClient", err.Error())
	}

	outJson("vault record", vaultRecord)

}

func outJson(name string, v any) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalln("outJson: ", err.Error())
	}
	log.Printf("%s: \n%s\n\n", name, out)
}

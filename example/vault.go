/* Licensed to the Apache Software Foundation (ASF) under one or more
   contributor license agreements.  See the NOTICE file distributed with
   this work for additional information regarding copyright ownership.
   The ASF licenses this file to You under the Apache License, Version 2.0
   (the "License"); you may not use this file except in compliance with
   the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License. */

package main

import (
	"flag"
	"log"

	"github.com/google/uuid"
	"github.com/topicuskeyhub/go-keyhub"
	"github.com/topicuskeyhub/go-keyhub/model"
	keyhubmodel "github.com/topicuskeyhub/go-keyhub/model"
)

func main() {

	// variables declaration
	var issuer string
	var clientid string
	var clientsecret string
	var groupadmin string
	var groupvault string

	// flags declaration using flag package
	flag.StringVar(&issuer, "i", "https://test.topicus-keyhub.com", "Specify issuer")
	flag.StringVar(&clientid, "ci", "", "Specify client id")
	flag.StringVar(&clientsecret, "cs", "", "Specify client secret")
	flag.StringVar(&groupadmin, "ga", "nil", "Specify UUID of existing account to become group admin")
	flag.StringVar(&groupvault, "gv", "nil", "Specify UUID of existing group vault to write in")

	flag.Parse()

	client, err := keyhub.NewClientDefault(issuer, clientid, clientsecret)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	//
	// Version Service
	//

	version, err := client.Version.Get()
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if len(version.ContractVersions) == 0 {
		log.Fatalln("No Contract versions")
	} else {
		log.Println("KeyHub version", version.KeyhubVersion, " has contract versions", version.ContractVersions)
	}

	//
	// Group Service
	//

	var foundGroupAdmin *keyhubmodel.Account
	if groupadmin != "nil" {
		a, _ := client.Accounts.GetByUUID(groupadmin)
		if a != nil {
			foundGroupAdmin = a
		}
	}
	if foundGroupAdmin == nil {
		log.Fatalln("ERROR", "No GroupAdmin found for UUID "+groupadmin)
	}

	//TODO activate when client can get permission to add vault records directly after creating a group

	// group, err := client.Groups.Create(keyhubmodel.NewGroup("Terraform "+uuid.NewString(), foundGroupAdmin))
	// if err != nil {
	// 	log.Fatalln("ERROR", err)
	// }
	// if group != nil {
	// 	log.Println("Created group. Result is", group.Name, ", UUID =", group.UUID)
	// }
	// if group.AdditionalObjects.Admins == nil ||
	// 	group.AdditionalObjects.Admins.Items == nil {
	// 	log.Fatalln("ERROR", "Group with UUID does not contain any group admins")
	// }
	// if group.AdditionalObjects.Admins.Items[0].UUID != groupadmin {
	// 	log.Fatalln("ERROR", "Group with UUID does not contain th correct group admin", groupadmin)
	// }

	// if groupvault == "nil" {
	// 	log.Fatalln("ERROR", "No Group specified to perform other Group & Vault operations on")
	// }

	groups, err := client.Groups.List()
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if len(groups) > 0 {
		log.Println("Get groups. First group is", groups[0].Name, ", UUID =", groups[0].UUID)
	}

	group, err := client.Groups.Get(groupvault)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if group != nil {
		log.Println("Get group. Result is", group.Name, ", UUID =", group.UUID)
	}

	//
	// Vault Service
	//

	vaultRecords, err := client.Vaults.GetRecords(group)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if len(vaultRecords) > 0 {
		log.Println("Get vaultRecords. First vaultRecord is", vaultRecords[0].Name, ", UUID =", vaultRecords[0].UUID)
	}

	password := "Banaan"
	secrets := &keyhubmodel.VaultRecordSecretAdditionalObject{
		Password: &password,
	}
	vaultRecord, err := client.Vaults.Create(group, keyhubmodel.NewVaultRecord("Random Password "+uuid.NewString(), secrets))
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if vaultRecord != nil {
		log.Println("Created vaultRecord. Result is", vaultRecord.Name, ", UUID =", vaultRecord.UUID)
	}
	if vaultRecord.AdditionalObjects == nil ||
		vaultRecord.AdditionalObjects.Secret == nil ||
		vaultRecord.AdditionalObjects.Secret.Password == nil {
		log.Fatalln("ERROR", "vaultRecord has no secrets")
	}

	vaultRecord, err = client.Vaults.GetByUUID(group, vaultRecord.UUID, &keyhubmodel.VaultRecordAdditionalQueryParams{Audit: true, Secret: true})
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if vaultRecord != nil {
		log.Println("Get vaultRecord. Result is", vaultRecord.Name, ", UUID =", vaultRecord.UUID)
	}
	if vaultRecord.AdditionalObjects == nil ||
		vaultRecord.AdditionalObjects.Audit == nil {
		log.Fatalln("ERROR", "vaultRecord has no audit")
	}
	if vaultRecord.AdditionalObjects == nil ||
		vaultRecord.AdditionalObjects.Secret == nil ||
		vaultRecord.AdditionalObjects.Secret.Password == nil {
		log.Fatalln("ERROR", "vaultRecord has no secrets")
	}

	vaultRecords, err = client.Vaults.List(group, &model.VaultRecordQueryParams{Name: vaultRecord.Name}, &model.VaultRecordAdditionalQueryParams{Secret: true})
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if len(vaultRecords) > 0 {
		log.Println("Get vaultRecords. First vaultRecord is", vaultRecords[0].Name, ", UUID =", vaultRecords[0].UUID)
	}

	password = "Banaan2"
	vaultRecord.AdditionalObjects.Secret.Password = &password
	vaultRecord, err = client.Vaults.Update(group, vaultRecord)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if vaultRecord != nil {
		log.Println("Updated vaultRecord. Result is", vaultRecord.Name, ", UUID =", vaultRecord.UUID)
	}

	err = client.Vaults.DeleteByUUID(group, vaultRecord.UUID)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	log.Println("Deleted vaultRecord", vaultRecord.UUID)
}

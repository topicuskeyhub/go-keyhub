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
		log.Fatalf("ERROR %s", err)
	}

	//
	// Version Service
	//

	version, err := client.Version.Get()
	if err != nil {
		log.Fatalf("ERROR %s", err)
	}
	if len(version.ContractVersions) == 0 {
		log.Fatalf("No Contract versions")
	} else {
		log.Printf("KeyHub version %s has contract versions %d", version.KeyhubVersion, version.ContractVersions)
	}

	//
	// Group Service
	//

	var foundGroupAdmin *keyhubmodel.Account
	if groupadmin != "nil" {
		groupadminUUID, err := uuid.Parse(groupadmin)
		if err != nil {
			log.Fatalf("ERROR Provided groupadmin value is not an UUID %q", groupadmin)
		}

		a, _ := client.Accounts.GetByUUID(groupadminUUID)
		if a != nil {
			foundGroupAdmin = a
		}

	}
	if foundGroupAdmin == nil {
		log.Fatalf("ERROR No GroupAdmin found for UUID %q", groupadmin)
	}

	//TODO activate when client can get permission to add vault records directly after creating a group

	// group, err := client.Groups.Create(keyhubmodel.NewGroup("Terraform "+uuid.NewString(), foundGroupAdmin))
	// if err != nil {
	// 	log.Fatalf("ERROR %s", err)
	// }
	// if group != nil {
	// log.Printf("Created group. Result is %q, UUID = %q", group.Name, group.UUID)
	// }
	// if group.AdditionalObjects.Admins == nil ||
	// 	group.AdditionalObjects.Admins.Items == nil {
	// 	log.Fatalf("ERROR Group with UUID does not contain any group admins")
	// }
	// if group.AdditionalObjects.Admins.Items[0].UUID != groupadmin {
	// 	log.Fatalf("ERROR Group with UUID does not contain the correct group admin %q", groupadmin)
	// }

	// if groupvault == "nil" {
	// 	log.Fatalf("ERROR No Group specified to perform other Group & Vault operations on")
	// }

	// groups, err := client.Groups.List()
	// if err != nil {
	// 	log.Fatalf("ERROR %s", err)
	// }
	// if len(groups) > 0 {
	// 	log.Printf("Get groups. First group is %q UUID = %q", groups[0].Name, groups[0].UUID)
	// }

	groupvaultUUID, err := uuid.Parse(groupvault)
	if err != nil {
		log.Fatalf("ERROR Provided groupvault value is not an UUID %q", groupvault)
	}

	group, err := client.Groups.GetByUUID(groupvaultUUID)
	if err != nil {
		log.Fatalf("ERROR %s", err)
	}
	if group != nil {
		log.Printf("Get group. Result is %q, UUID = %q", group.Name, group.UUID)
	}

	//
	// Vault Service
	//

	vaultRecords, err := client.Vaults.GetRecords(group)
	if err != nil {
		log.Fatalf("ERROR %s", err)
	}
	if len(vaultRecords) > 0 {
		log.Printf("Get vaultRecords. First vaultRecord is %q, UUID = %q", vaultRecords[0].Name, vaultRecords[0].UUID)
	}

	password := "Banaan"
	secrets := &keyhubmodel.VaultRecordSecretAdditionalObject{
		Password: &password,
	}
	vaultRecord, err := client.Vaults.Create(group, keyhubmodel.NewVaultRecord("Random Password "+uuid.NewString(), secrets))
	if err != nil {
		log.Fatalf("ERROR %s", err)
	}
	if vaultRecord != nil {
		log.Printf("Created vaultRecord. Result is %q, UUID = %q", vaultRecord.Name, vaultRecord.UUID)
	}
	if vaultRecord.AdditionalObjects == nil ||
		vaultRecord.AdditionalObjects.Secret == nil ||
		vaultRecord.AdditionalObjects.Secret.Password == nil {
		log.Fatalf("ERROR vaultRecord has no secrets")
	}

	vaultrecordUUID, _ := uuid.Parse(vaultRecord.UUID)
	vaultRecord, err = client.Vaults.GetByUUID(group, vaultrecordUUID, &keyhubmodel.VaultRecordAdditionalQueryParams{Audit: true, Secret: true})
	if err != nil {
		log.Fatalf("ERROR %s", err)
	}
	if vaultRecord != nil {
		log.Printf("Get vaultRecord. Result is %q, UUID = %q", vaultRecord.Name, vaultRecord.UUID)
	}
	if vaultRecord.AdditionalObjects == nil ||
		vaultRecord.AdditionalObjects.Audit == nil {
		log.Fatalf("ERROR vaultRecord has no audit")
	}
	if vaultRecord.AdditionalObjects == nil ||
		vaultRecord.AdditionalObjects.Secret == nil ||
		vaultRecord.AdditionalObjects.Secret.Password == nil {
		log.Fatalf("ERROR vaultRecord has no secrets")
	}

	vaultRecords, err = client.Vaults.List(group, &model.VaultRecordQueryParams{Name: vaultRecord.Name}, &model.VaultRecordAdditionalQueryParams{Secret: true})
	if err != nil {
		log.Fatalf("ERROR %s", err)
	}
	if len(vaultRecords) > 0 {
		log.Printf("Get vaultRecords. First vaultRecord is %q, UUID = %q", vaultRecords[0].Name, vaultRecords[0].UUID)
	}

	password = "Banaan2"
	vaultRecord.AdditionalObjects.Secret.Password = &password
	vaultRecord, err = client.Vaults.Update(group, vaultRecord)
	if err != nil {
		log.Fatalf("ERROR %s", err)
	}
	if vaultRecord != nil {
		log.Printf("Updated vaultRecord. Result %q, UUID = %q", vaultRecord.Name, vaultRecord.UUID)
	}

	vaultrecordUUID, _ = uuid.Parse(vaultRecord.UUID)
	err = client.Vaults.DeleteByUUID(group, vaultrecordUUID)
	if err != nil {
		log.Fatalf("ERROR %s", err)
	}
	log.Printf("Deleted vaultRecord %q", vaultRecord.UUID)
}

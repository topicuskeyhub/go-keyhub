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
	"net/http"

	"github.com/google/uuid"
	"github.com/topicuskeyhub/go-keyhub"
	keyhubmodel "github.com/topicuskeyhub/go-keyhub/model"
)

func main() {

	// variables declaration
	var issuer string
	var clientid string
	var clientsecret string
	var groupadmin string

	// flags declaration using flag package
	flag.StringVar(&issuer, "i", "https://test.topicus-keyhub.com", "Specify issuer")
	flag.StringVar(&clientid, "ci", "", "Specify client id")
	flag.StringVar(&clientsecret, "cs", "", "Specify client secret")
	flag.StringVar(&groupadmin, "ga", "nil", "Specify group admin by UUID")

	flag.Parse()

	client, err := keyhub.NewClient(http.DefaultClient, issuer, clientid, clientsecret)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	var foundGroupAdmin *keyhubmodel.Account
	if groupadmin != "nil" {
		a, _ := client.Accounts.Get(groupadmin)
		if a != nil {
			foundGroupAdmin = a
		}
	}

	if foundGroupAdmin == nil {
		log.Fatalln("ERROR", "No GroupAdmin found for UUID "+groupadmin)
	}

	group, err := client.Groups.Create(keyhubmodel.NewGroup("Terraform "+uuid.NewString(), foundGroupAdmin))
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if group != nil {
		log.Println("Created group. Result is", group.Name, ", UUID =", group.UUID)
	}

	groups, err := client.Groups.List()
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if len(groups) > 0 {
		log.Println("Get groups. First group is", groups[0].Name, ", UUID =", group.UUID)
	}
}

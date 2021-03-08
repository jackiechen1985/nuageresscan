// Copyright (C) 2021 Nokia-Sbell Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"io/ioutil"
)

/*
{
  "neutron": {
    "username": "neutron",
    "password": "0c75bcb291a94015",
    "ip_addr": "135.251.96.37",
    "port": 3306,
    "db_name": "neutron"
  },
  "vsd": [
    {
      "username": "csproot",
      "password": "csproot",
      "organization": "csp",
      "url": "https://135.251.96.136:8443",
      "net_partition": "OpenStack_Pike_beijing",
      "cms_id": "ebbdadd3-cc73-42ed-9a04-361d99e12aee",
      "az": "beijing"
    },
    {
      "username": "csproot",
      "password": "csproot",
      "organization": "csp",
      "url": "https://135.251.96.137:8443",
      "net_partition": "OpenStack_pike",
      "cms_id": "e4e6b7b1-202a-4d7a-87f2-412beb513d17",
      "az": "changsha"
    }
  ]
}
*/

type Neutron struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IPAddr   string `json:"ip_addr"`
	Port     uint16 `json:"port"`
	DBName   string `json:"db_name"`
}

type VSD struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
	URL          string `json:"url"`
	NetPartition string `json:"net_partition"`
	CMSID        string `json:"cms_id"`
	AZ           string `json:"az"`
}

type Config struct {
	Neu  Neutron `json:"neutron"`
	Vsds []VSD   `json:"vsd"`
}

func LoadConfig(configPath string) (*Config, error) {
	buf, err := ioutil.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetCMSID(config *Config, az string) string {
	for _, vsd := range config.Vsds {
		if vsd.AZ == az {
			return vsd.CMSID
		}
	}

	return ""
}

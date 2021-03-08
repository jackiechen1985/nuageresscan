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
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	ResTypeSubnet        string = "subnet"
	ResTypeRouter        string = "router"
	ResTypePort          string = "port"
	ResTypeDummyfip      string = "dummyfip"
	ResTypeSecuritygroup string = "securitygroup"
	ResTypeUnderlayacl   string = "underlayacl"
)

var globalConfig *Config

func printUsage() {
	s := fmt.Sprintf(
		`Usage:
  nuageresscan [config] [%s|%s|%s|%s|%s|%s]

Flags:
  -h, --help             help for program
  -v, --version          show program version
  -i, --info             set log level to info
}`, ResTypeSubnet, ResTypeRouter, ResTypePort, ResTypeDummyfip, ResTypeSecuritygroup, ResTypeUnderlayacl)
	fmt.Println(s)
}

func startJob(configPath string, resourceType string) error {
	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %s", err)
	}

	globalConfig = config

	neu := globalConfig.Neu
	err = OpenDB(neu.Username, neu.Password, neu.IPAddr, neu.Port, neu.DBName)
	if err != nil {
		return fmt.Errorf("failed to open database: %s", err)
	}

	switch resourceType {
	case ResTypeSubnet:
		scanResForSubnet()
	case ResTypeRouter:
		scanResForRouter()
	case ResTypePort:
		scanResForPort()
	case ResTypeDummyfip:
		scanResForDummyFip()
	case ResTypeSecuritygroup:
		scanResForSecurityGroup()
	case ResTypeUnderlayacl:
		scanResForUnderlayAcl()
	default:
		logrus.WithField("func", "startJob").
			Error("Unknown resource type:" + resourceType)
		printUsage()
	}
	return nil
}

func main() {
	logrus.SetLevel(logrus.WarnLevel)

	configPos := 1
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			printUsage()
			os.Exit(0)
		} else if os.Args[1] == "-v" || os.Args[1] == "--version" {
			fmt.Println("nuageresscan version", Version())
			os.Exit(0)
		} else if os.Args[1] == "-i" || os.Args[1] == "--info" {
			logrus.SetLevel(logrus.InfoLevel)
			configPos = 2
		}
	}

	if len(os.Args) > configPos+1 {
		err := startJob(os.Args[configPos], os.Args[configPos+1])
		if err != nil {
			logrus.WithField("func", "main").Error(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	printUsage()
	os.Exit(1)
}

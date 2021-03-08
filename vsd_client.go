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
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

const maxPageSize = 500

func StartSession(username string, password string, organization string, url string) (*vspk.Me, error) {
	session, me := vspk.NewSession(username, password, organization, url)
	err := session.Start()
	if err != nil {
		return nil, err
	}
	return me, nil
}

func FetchEnterpriseByName(me *vspk.Me, name string) (*vspk.Enterprise, error) {
	filter := fmt.Sprintf("name == '%s'", name)
	enterprises, err := me.Enterprises(&bambou.FetchingInfo{Filter: filter})
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	return enterprises[0], nil
}

func FetchAllL2DomainTemplates(enterprise *vspk.Enterprise) (vspk.L2DomainTemplatesList, error) {
	var allL2DomainTemplates vspk.L2DomainTemplatesList
	for page := 0; ; page++ {
		l2DomainTemplates, err := enterprise.L2DomainTemplates(&bambou.FetchingInfo{Page: page, PageSize: maxPageSize})
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		if l2DomainTemplates == nil {
			break
		}
		allL2DomainTemplates = append(allL2DomainTemplates, l2DomainTemplates...)
	}

	return allL2DomainTemplates, nil
}

func FetchAllL2Domains(enterprise *vspk.Enterprise) (vspk.L2DomainsList, error) {
	var allL2Domains vspk.L2DomainsList
	for page := 0; ; page++ {
		l2Domains, err := enterprise.L2Domains(&bambou.FetchingInfo{Page: page, PageSize: maxPageSize})
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		if l2Domains == nil {
			break
		}
		allL2Domains = append(allL2Domains, l2Domains...)
	}

	return allL2Domains, nil
}

func FetchAllDomains(enterprise *vspk.Enterprise) (vspk.DomainsList, error) {
	var allDomains vspk.DomainsList
	for page := 0; ; page++ {
		domains, err := enterprise.Domains(&bambou.FetchingInfo{Page: page, PageSize: maxPageSize})
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		if domains == nil {
			break
		}
		allDomains = append(allDomains, domains...)
	}

	return allDomains, nil
}

func FetchAllSubnets(domain *vspk.Domain) (vspk.SubnetsList, error) {
	var allSubnets vspk.SubnetsList
	for page := 0; ; page++ {
		subnets, err := domain.Subnets(&bambou.FetchingInfo{Page: page, PageSize: maxPageSize})
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		if subnets == nil {
			break
		}
		allSubnets = append(allSubnets, subnets...)
	}

	return allSubnets, nil
}

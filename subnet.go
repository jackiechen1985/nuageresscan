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
	"github.com/nuagenetworks/vspk-go/vspk"
	"github.com/sirupsen/logrus"
)

// Neutron resources
var neutronSubnets []Subnet
var neutronL2domMappings []NuageSubnetL2domMapping
var neutronSubnetMap map[string]*Subnet
var neutronL2domMappingSubnetIDMap map[string]*NuageSubnetL2domMapping
var neutronL2domMappingNuageSubnetIDMap map[string]*NuageSubnetL2domMapping
var neutronL2domMappingNuageL2domTmpltIDMap map[string]*NuageSubnetL2domMapping

// Nuage resources
var nuageL2DomainTemplates vspk.L2DomainTemplatesList
var nuageL2Domains vspk.L2DomainsList
var nuageSubnets vspk.SubnetsList
var nuageL2DomainTemplateMap map[string]*vspk.L2DomainTemplate
var nuageL2DomainMap map[string]*vspk.L2Domain
var nuageSubnetMap map[string]*vspk.Subnet

func dumpAllNeutronSubnetResources() error {
	logrus.WithField("func", "dumpAllNeutronSubnetResources").
		Info("SelectAllSubnets")
	err := SelectAllSubnets(&neutronSubnets)
	if err != nil {
		return err
	}
	neutronSubnetMap = make(map[string]*Subnet)
	for i := 0; i < len(neutronSubnets); i++ {
		subnet := &neutronSubnets[i]
		neutronSubnetMap[subnet.ID] = subnet
	}

	logrus.WithField("func", "dumpAllNeutronSubnetResources").
		Info("SelectAllNuageSubnetL2domMappings")
	err = SelectAllNuageSubnetL2domMappings(&neutronL2domMappings)
	if err != nil {
		return err
	}
	neutronL2domMappingSubnetIDMap = make(map[string]*NuageSubnetL2domMapping)
	neutronL2domMappingNuageSubnetIDMap = make(map[string]*NuageSubnetL2domMapping)
	neutronL2domMappingNuageL2domTmpltIDMap = make(map[string]*NuageSubnetL2domMapping)
	for i := 0; i < len(neutronL2domMappings); i++ {
		l2domMapping := &neutronL2domMappings[i]
		neutronL2domMappingSubnetIDMap[l2domMapping.SubnetID] = l2domMapping
		neutronL2domMappingNuageSubnetIDMap[l2domMapping.NuageSubnetID] = l2domMapping
		if l2domMapping.NuageL2domTmpltID.Valid {
			neutronL2domMappingNuageL2domTmpltIDMap[l2domMapping.NuageL2domTmpltID.String] = l2domMapping
		}
	}

	return nil
}

func dumpAllNuageL2DomainResources() error {
	for _, vsd := range globalConfig.Vsds {
		me, err := StartSession(vsd.Username, vsd.Password, vsd.Organization, vsd.URL)
		if err != nil {
			return err
		}

		enterprise, err := FetchEnterpriseByName(me, vsd.NetPartition)
		if err != nil {
			return err
		}

		logrus.WithField("func", "dumpAllNuageL2DomainResources").
			Info("FetchAllL2DomainTemplates from " + vsd.NetPartition)
		l2domTmplts, err := FetchAllL2DomainTemplates(enterprise)
		if err != nil {
			return err
		}
		nuageL2DomainTemplates = append(nuageL2DomainTemplates, l2domTmplts...)
		nuageL2DomainTemplateMap = make(map[string]*vspk.L2DomainTemplate)
		for _, l2domTmplt := range l2domTmplts {
			nuageL2DomainTemplateMap[l2domTmplt.ID] = l2domTmplt
		}

		logrus.WithField("func", "dumpAllNuageL2DomainResources").
			Info("FetchAllL2Domains from " + vsd.NetPartition)
		l2doms, err := FetchAllL2Domains(enterprise)
		if err != nil {
			return err
		}
		nuageL2Domains = append(nuageL2Domains, l2doms...)
		nuageL2DomainMap = make(map[string]*vspk.L2Domain)
		for _, l2dom := range l2doms {
			nuageL2DomainMap[l2dom.ID] = l2dom
		}

		logrus.WithField("func", "dumpAllNuageL2DomainResources").
			Info("FetchAllDomains from " + vsd.NetPartition)
		domains, err := FetchAllDomains(enterprise)
		if err != nil {
			return err
		}
		for _, domain := range domains {
			logrus.WithField("func", "dumpAllNuageL2DomainResources").
				Info("FetchAllSubnets from " + vsd.NetPartition)
			subnets, err := FetchAllSubnets(domain)
			if err != nil {
				return err
			}
			nuageSubnets = append(nuageSubnets, subnets...)
		}
		nuageSubnetMap = make(map[string]*vspk.Subnet)
		for _, subnet := range nuageSubnets {
			nuageSubnetMap[subnet.ID] = subnet
		}
	}

	return nil
}

func scanResForSubnetBaseOnNeutron() {
	for _, neutronSubnet := range neutronSubnets {
		neutronL2domMapping := neutronL2domMappingSubnetIDMap[neutronSubnet.ID]
		if neutronL2domMapping == nil {
			logrus.WithFields(logrus.Fields{"func": "scanResForSubnetBaseOnNeutron", "object": "neutron"}).
				Warningf("nuage_subnet_l2dom_mapping.subnet_id %s was not found", neutronSubnet.ID)
			continue
		}

		if neutronL2domMapping.NuageL2domTmpltID.Valid {
			nuageL2DomainTemplate := nuageL2DomainTemplateMap[neutronL2domMapping.NuageL2domTmpltID.String]
			if nuageL2DomainTemplate == nil {
				logrus.WithFields(logrus.Fields{"func": "scanResForSubnetBaseOnNeutron", "object": "nuage"}).
					Warningf("l2domain template %s was not found", neutronL2domMapping.NuageL2domTmpltID.String)
			}

			nuageL2Domain := nuageL2DomainMap[neutronL2domMapping.NuageSubnetID]
			if nuageL2Domain == nil {
				logrus.WithFields(logrus.Fields{"func": "scanResForSubnetBaseOnNeutron", "object": "nuage"}).
					Warningf("l2domain %s was not found", neutronL2domMapping.NuageSubnetID)
			}
		} else {
			nuageSubnet := nuageSubnetMap[neutronL2domMapping.NuageSubnetID]
			if nuageSubnet == nil {
				logrus.WithFields(logrus.Fields{"func": "scanResForSubnetBaseOnNeutron", "object": "nuage"}).
					Warningf("subnet %s was not found", neutronL2domMapping.NuageSubnetID)
			}
		}
	}
}

func scanResForSubnetBaseOnNuage() {
	for _, nuageL2DomainTemplate := range nuageL2DomainTemplates {
		neutronL2domMapping := neutronL2domMappingNuageL2domTmpltIDMap[nuageL2DomainTemplate.ID]
		if neutronL2domMapping == nil {
			logrus.WithFields(logrus.Fields{"func": "scanResForSubnetBaseOnNuage", "object": "neutron"}).
				Warningf("nuage_subnet_l2dom_mapping.nuage_l2dom_tmplt_id %s was not found", nuageL2DomainTemplate.ID)
			continue
		}
		neutronSubnet := neutronSubnetMap[neutronL2domMapping.SubnetID]
		if neutronSubnet == nil {
			logrus.WithFields(logrus.Fields{"func": "scanResForSubnetBaseOnNuage", "object": "neutron"}).
				Warningf("subnet.id %s was not found", neutronL2domMapping.SubnetID)
		}
	}

	for _, nuageL2Domain := range nuageL2Domains {
		neutronL2domMapping := neutronL2domMappingNuageSubnetIDMap[nuageL2Domain.ID]
		if neutronL2domMapping == nil {
			logrus.WithFields(logrus.Fields{"func": "scanResForSubnetBaseOnNuage", "object": "neutron"}).
				Warningf("nuage_subnet_l2dom_mapping.nuage_subnet_id %s was not found", nuageL2Domain.ID)
			continue
		}
		neutronSubnet := neutronSubnetMap[neutronL2domMapping.SubnetID]
		if neutronSubnet == nil {
			logrus.WithFields(logrus.Fields{"func": "scanResForSubnetBaseOnNuage", "object": "neutron"}).
				Warningf("subnet.id %s was not found", neutronL2domMapping.SubnetID)
		}
	}

	for _, nuageSubnet := range nuageSubnets {
		neutronL2domMapping := neutronL2domMappingNuageSubnetIDMap[nuageSubnet.ID]
		if neutronL2domMapping == nil {
			logrus.WithFields(logrus.Fields{"func": "scanResForSubnetBaseOnNuage", "object": "neutron"}).
				Warningf("nuage_subnet_l2dom_mapping.nuage_subnet_id %s was not found", nuageSubnet.ID)
			continue
		}
		neutronSubnet := neutronSubnetMap[neutronL2domMapping.SubnetID]
		if neutronSubnet == nil {
			logrus.WithFields(logrus.Fields{"func": "scanResForSubnetBaseOnNuage", "object": "neutron"}).
				Warningf("subnet.id %s was not found", neutronL2domMapping.SubnetID)
		}
	}
}

func scanResForSubnet() {
	err := dumpAllNeutronSubnetResources()
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "dumpAllNeutronSubnetResources", "object": "neutron"}).Error(err)
		return
	}

	err = dumpAllNuageL2DomainResources()
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "dumpAllNuageL2DomainResources", "object": "nuage"}).Error(err)
		return
	}

	scanResForSubnetBaseOnNeutron()
	scanResForSubnetBaseOnNuage()
}

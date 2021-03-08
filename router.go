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
	"strings"
)

// Neutron resources
var neutronRouters []Router
var neutronRouterMap map[string]*Router

// Nuage resources
var nuageDomains vspk.DomainsList
var nuageDomainMap map[string]*vspk.Domain

func dumpAllNeutronRouterResources() error {
	logrus.WithField("func", "dumpAllNeutronRouterResources").
		Info("SelectAllRouters")
	err := SelectAllRouters(&neutronRouters)
	if err != nil {
		return err
	}
	neutronRouterMap = make(map[string]*Router)
	for i := 0; i < len(neutronRouters); i++ {
		router := &neutronRouters[i]
		neutronRouterMap[router.ID] = router
	}

	return nil
}

func dumpAllNuageDomainResources() error {
	nuageDomainMap = make(map[string]*vspk.Domain)

	for _, vsd := range globalConfig.Vsds {
		me, err := StartSession(vsd.Username, vsd.Password, vsd.Organization, vsd.URL)
		if err != nil {
			return err
		}

		enterprise, err := FetchEnterpriseByName(me, vsd.NetPartition)
		if err != nil {
			return err
		}

		logrus.WithField("func", "dumpAllNuageDomainResources").
			Info("FetchAllDomains from " + vsd.NetPartition)
		domains, err := FetchAllDomains(enterprise)
		if err != nil {
			return err
		}
		nuageDomains = append(nuageDomains, domains...)
		for _, domain := range domains {
			if domain.ExternalID == "" {
				logrus.WithFields(logrus.Fields{"func": "dumpAllNuageDomainResources", "object": "nuage"}).
					Warningf("found domain %s with empty externalID", domain.ID)
				continue
			}
			if nuageDomainMap[domain.ExternalID] != nil {
				logrus.WithFields(logrus.Fields{"func": "dumpAllNuageDomainResources", "object": "nuage"}).
					Warning("found redundant domain " + domain.ExternalID)
				continue
			}
			nuageDomainMap[domain.ExternalID] = domain
		}
	}

	return nil
}

func scanResForRouterBaseOnNeutron() {
	for _, neutronRouter := range neutronRouters {
		var routerPorts []RouterPort
		err := SelectRouterPortsByRouterID(&routerPorts, neutronRouter.ID)
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "scanResForRouterBaseOnNeutron", "object": "neutron"}).Error(err)
		}

		if len(routerPorts) <= 0 {
			continue
		}

		var newarchAzRouterNuages []NewarchAzRouterNuage
		err = SelectNewarchAzRouterNuagesByRouterID(&newarchAzRouterNuages, neutronRouter.ID)
		for _, newarchAzRouterNuage := range newarchAzRouterNuages {
			if !newarchAzRouterNuage.AzName.Valid {
				logrus.WithFields(logrus.Fields{"func": "scanResForRouterBaseOnNeutron", "object": "neutron"}).
					Warning("az is null for router " + neutronRouter.ID)
				continue
			}

			cmsID := GetCMSID(globalConfig, newarchAzRouterNuage.AzName.String)
			if cmsID == "" {
				logrus.WithFields(logrus.Fields{"func": "scanResForRouterBaseOnNeutron", "object": "neutron"}).
					Warning("cms_id was not found for AZ" + newarchAzRouterNuage.AzName.String)
				continue
			}

			externalID := neutronRouter.ID + "@" + cmsID
			nuageDomain := nuageDomainMap[externalID]
			if nuageDomain == nil {
				logrus.WithFields(logrus.Fields{"func": "scanResForRouterBaseOnNeutron", "object": "nuage"}).
					Warningf("domain %s was not found", externalID)
			}
		}
	}
}

func scanResForRouterBaseOnNuage() {
	for _, nuageDomain := range nuageDomains {
		if nuageDomain.ExternalID == "" {
			logrus.WithFields(logrus.Fields{"func": "scanResForRouterBaseOnNuage", "object": "nuage"}).
				Warningf("found domain %s with empty externalID", nuageDomain.ID)
			continue
		}
		neutronRouterID := strings.Split(nuageDomain.ExternalID, "@")[0]
		if neutronRouterID == "" {
			logrus.WithFields(logrus.Fields{"func": "scanResForRouterBaseOnNuage", "object": "nuage"}).
				Warning("invalid domain externalID " + nuageDomain.ExternalID)
			continue
		}
		neutronRouter := neutronRouterMap[neutronRouterID]
		if neutronRouter == nil {
			logrus.WithFields(logrus.Fields{"func": "scanResForRouterBaseOnNuage", "object": "neutron"}).
				Warningf("router.id %s was not found", neutronRouterID)
			continue
		}
	}
}

func scanResForRouter() {
	err := dumpAllNeutronRouterResources()
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "dumpAllNeutronRouterResources", "object": "neutron"}).Error(err)
		return
	}

	err = dumpAllNuageDomainResources()
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "dumpAllNuageDomainResources", "object": "nuage"}).Error(err)
		return
	}

	scanResForRouterBaseOnNeutron()
	scanResForRouterBaseOnNuage()
}

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
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

type Network struct {
	ID                    string         `db:"id"`
	AvailabilityZoneHints sql.NullString `db:"availability_zone_hints"`
}

type Subnet struct {
	ID        string `db:"id"`
	NetworkID string `db:"network_id"`
}

type NuageSubnetL2domMapping struct {
	SubnetID          string         `db:"subnet_id"`
	NuageSubnetID     string         `db:"nuage_subnet_id"`
	NuageL2domTmpltID sql.NullString `db:"nuage_l2dom_tmplt_id"`
}

type Router struct {
	ID string `db:"id"`
}

type RouterPort struct {
	RouterID string `db:"router_id"`
	PortID   string `db:"port_id"`
}

type NewarchAzRouterNuage struct {
	RouterID      sql.NullString `db:"router_id"`
	AzName        sql.NullString `db:"az_name"`
	NuageRouterID sql.NullString `db:"nuage_router_id"`
}

func OpenDB(username string, password string, ipAddr string, port uint16, dbName string) error {
	dsName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", username, password, ipAddr, port, dbName)
	db, err := sqlx.Open("mysql", dsName)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func SelectAllNetworks(networks *[]Network) error {
	return DB.Select(networks, "select id, availability_zone_hints from networks")
}

func SelectAllSubnets(subnets *[]Subnet) error {
	return DB.Select(subnets, "select id, network_id from subnets")
}

func SelectAllNuageSubnetL2domMappings(l2domMappings *[]NuageSubnetL2domMapping) error {
	return DB.Select(l2domMappings, "select subnet_id, nuage_subnet_id, nuage_l2dom_tmplt_id from nuage_subnet_l2dom_mapping")
}

func SelectAllRouters(routers *[]Router) error {
	return DB.Select(routers, "select id from routers")
}

func SelectRouterPortsByRouterID(routerPorts *[]RouterPort, routerID string) error {
	return DB.Select(routerPorts, "select router_id, port_id from routerports where router_id=?", routerID)
}

func SelectNewarchAzRouterNuagesByRouterID(newarchAzRouterNuage *[]NewarchAzRouterNuage, routerID string) error {
	return DB.Select(newarchAzRouterNuage, "select router_id, az_name, nuage_router_id from newarch_az_router_nuage where router_id=?", routerID)
}

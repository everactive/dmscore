// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Identity Service
 * Copyright 2019 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it
 * under the terms of the GNU Affero General Public License version 3, as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY,
 * SATISFACTORY QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.
 * See the GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package postgres

import (
	"database/sql"
	"github.com/everactive/dmscore/iot-identity/datastore"
	"github.com/everactive/dmscore/iot-identity/domain"
	"github.com/everactive/dmscore/iot-identity/models"
)

// OrganizationNew creates a new organization
func (s *Store) OrganizationNew(org datastore.OrganizationNewRequest) (string, error) {
	// var id int64
	var orgID = datastore.GenerateID()
	res := s.gormDB.Create(&models.Organization{
		OrgId:       orgID,
		Name:        org.Name,
		CountryName: org.CountryName,
		RootCert:    string(org.ServerCert),
		RootKey:     string(org.ServerKey),
	})

	if res.Error != nil {
		return "", res.Error
	}

	return orgID, nil
}

// OrganizationList fetches existing organizations
func (s *Store) OrganizationList() ([]domain.Organization, error) {
	var id int64
	rows, err := s.Query(listOrganizationSQL)
	if err != nil {
		datastore.Logger.Errorf("Error retrieving organizations: %v\n", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			datastore.Logger.Error(err)
		}
	}(rows)

	items := []domain.Organization{}
	for rows.Next() {
		item := domain.Organization{}
		err := rows.Scan(&id, &item.ID, &item.Name, &item.RootCert)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// OrganizationGet fetches an organization by ID
func (s *Store) OrganizationGet(orgID string) (*domain.Organization, error) {
	var id int64
	var countryName string
	org := domain.Organization{}

	err := s.QueryRow(getOrganizationSQL, orgID).Scan(&id, &org.ID, &org.Name, &countryName, &org.RootCert, &org.RootKey)
	if err != nil {
		datastore.Logger.Errorf("Error retrieving organization %v: %v\n", orgID, err)
	}
	return &org, err
}

// OrganizationGetByName fetches an organization by name
func (s *Store) OrganizationGetByName(name string) (*domain.Organization, error) {
	var id int64
	var countryName string
	org := domain.Organization{}

	err := s.QueryRow(getOrganizationByNameSQL, name).Scan(&id, &org.ID, &org.Name, &countryName, &org.RootCert, &org.RootKey)
	if err != nil {
		datastore.Logger.Errorf("Error retrieving organization `%v`: %v\n", name, err)
	}
	return &org, err
}

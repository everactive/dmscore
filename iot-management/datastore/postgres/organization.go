// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Management Service
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
	"fmt"
	"github.com/everactive/dmscore/models"
	log "github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-management/datastore"
)

// createOrganizationTable creates the database table for devices with its indexes.
func (s *Store) createOrganizationTable() error {
	_, err := s.Exec(createOrganizationTableSQL)
	return err
}

// createOrganizationUserTable creates the database table for devices with its indexes.
func (s *Store) createOrganizationUserTable() error {
	_, err := s.Exec(createOrganizationUserTableSQL)
	return err
}

// OrganizationGet returns an organization
func (s *Store) OrganizationGet(orgIDOrName string) (datastore.Organization, error) {
	org := models.Organization{}
	res := s.gormDB.Where("code = ?", orgIDOrName).Or("name", orgIDOrName).Find(&org)
	if res.Error != nil {
		return datastore.Organization{}, res.Error
	}

	return datastore.Organization{
		OrganizationID: org.OrganizationID,
		Name:           org.Name,
	}, nil
}

// OrgUserAccess checks if the user has permissions to access the organization
func (s *Store) OrgUserAccess(orgID, username string, role int) bool {
	// Superusers can access all accounts
	if role == datastore.Superuser {
		return true
	}

	var linkExists bool
	err := s.QueryRow(organizationUserAccessSQL, orgID, username).Scan(&linkExists)
	if err != nil {
		log.Printf("Error verifying the account-user link: %v\n", err)
		return false
	}
	return linkExists
}

// OrganizationsForUser returns the organizations a user can access
func (s *Store) OrganizationsForUser(username string) ([]datastore.Organization, error) {
	var sqlStatement string

	// Check if the user is a superuser
	user, err := s.GetUser(username)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %v", err)
	}

	// No restrictions for the superuser
	if user.Role == datastore.Superuser {
		sqlStatement = listOrganizationsSQL
	} else {
		sqlStatement = listUserOrganizationsSQL
	}

	rows, err := s.Query(sqlStatement, username)

	if err != nil {
		log.Printf("Error retrieving database accounts: %v\n", err)
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	return rowsToOrganizations(rows)
}

// OrganizationForUserToggle toggles the user access to an organization
func (s *Store) OrganizationForUserToggle(orgID, username string) error {
	result, err := s.Exec(deleteOrganizationUserAccessSQL, orgID, username)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		// Create the link
		_, err := s.Exec(createOrganizationUserAccessSQL, orgID, username)
		if err != nil {
			return err
		}
	}

	return nil
}

func rowsToOrganizations(rows *sql.Rows) ([]datastore.Organization, error) {
	orgs := []datastore.Organization{}

	for rows.Next() {
		org := datastore.Organization{}
		err := rows.Scan(&org.OrganizationID, &org.Name)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}

	return orgs, nil
}

// OrganizationCreate creates a new organization
func (s *Store) OrganizationCreate(org datastore.Organization) error {
	res := s.gormDB.Create(&models.Organization{
		OrganizationID: org.OrganizationID,
		Name:           org.Name,
	})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// OrganizationUpdate updates an organization
func (s *Store) OrganizationUpdate(org datastore.Organization) error {
	_, err := s.Exec(updateOrganizationSQL, org.OrganizationID, org.Name)
	return err
}

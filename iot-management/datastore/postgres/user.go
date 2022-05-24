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

	log "github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-management/datastore"
)

// CreateUser creates a new user
func (s *Store) CreateUser(user datastore.User) (int64, error) {
	var createdUserID int64

	err := s.QueryRow(createUserSQL, user.Username, user.Name, user.Email, user.Role).Scan(&createdUserID)
	if err != nil {
		log.Printf("Error creating user `%s`: %v\n", user.Username, err)
	}

	return createdUserID, err
}

// UserUpdate updates a user
func (s *Store) UserUpdate(user datastore.User) error {
	_, err := s.Exec(updateUserSQL, user.Username, user.Name, user.Email, user.Role)
	return err
}

// UserDelete removes a user
func (s *Store) UserDelete(username string) error {
	_, err := s.Exec(deleteUserSQL, username)
	return err
}

// UserList lists existing users
func (s *Store) UserList() ([]datastore.User, error) {
	rows, err := s.Query(listUsersSQL)
	if err != nil {
		log.Printf("Error retrieving database users: %v\n", err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	return s.rowsToUsers(rows)
}

// GetUser gets an existing user
func (s *Store) GetUser(username string) (datastore.User, error) {
	row := s.QueryRow(getUserSQL, username)
	user, err := s.rowToUser(row)
	if err != nil {
		log.Printf("Error retrieving user %v: %v\n", username, err)
	}
	return user, err
}

func (s *Store) rowToUser(row *sql.Row) (datastore.User, error) {
	user := datastore.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Role)
	if err != nil {
		return datastore.User{}, err
	}

	return user, nil
}

func (s *Store) rowsToUser(rows *sql.Rows) (datastore.User, error) {
	user := datastore.User{}
	err := rows.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Role)
	if err != nil {
		return datastore.User{}, err
	}

	return user, nil
}

func (s *Store) rowsToUsers(rows *sql.Rows) ([]datastore.User, error) {
	users := []datastore.User{}

	for rows.Next() {
		user, err := s.rowsToUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

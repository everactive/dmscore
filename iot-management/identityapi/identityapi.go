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

// Package identityapi provides an interface and types for interacting with the Identity REST API
package identityapi

import (
	"github.com/go-resty/resty/v2"
)

// Client is a client for the identity API
type Client interface {
}

// ClientAdapter adapts our expectations to device twin API
type ClientAdapter struct {
	URL    string
	client *resty.Client
}

// NewClientAdapter creates an adapter to access the device twin service
func NewClientAdapter(u string) (*ClientAdapter, error) {
	adapter := &ClientAdapter{
		URL:    u,
		client: resty.New(),
	}

	return adapter, nil
}

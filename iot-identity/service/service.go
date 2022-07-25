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

// Package service implements the Identity interface and data access methods
package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/everactive/dmscore/config/keys"
	"github.com/everactive/dmscore/iot-identity/models"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/everactive/dmscore/iot-identity/datastore"
	"github.com/everactive/dmscore/iot-identity/domain"
	"github.com/snapcore/snapd/asserts"
)

const (
	AccountKeyGetURL = "https://api.snapcraft.io/api/v1/snaps/assertions/account-key/"
)

// Logger is a logger specific to the service that can be swapped out and only affect it
var Logger = log.StandardLogger()

// Identity interface for the service
type Identity interface {
	RegisterOrganization(req *RegisterOrganizationRequest) (string, error)
	RegisterDevice(req *RegisterDeviceRequest) (string, error)
	DeleteDevice(deviceID string) (string, error)
	OrganizationList() ([]domain.Organization, error)
	DeviceList(orgID string) ([]domain.Enrollment, error)
	DeviceGet(orgID, deviceID string) (*domain.Enrollment, error)
	DeviceUpdate(orgID, deviceID string, req *DeviceUpdateRequest) error

	EnrollDevice(req *EnrollDeviceRequest) (*domain.Enrollment, error)
}

// IdentityService implementation of the identity use cases
type IdentityService struct {
	DB                       datastore.DataStore
	allowedSignKeyIDs        []string
	allowedSignKeyPublicKeys map[string]asserts.PublicKey
}

// NewIdentityService creates an implementation of the identity use cases
func NewIdentityService(db datastore.DataStore) *IdentityService {
	ids := &IdentityService{
		DB: db,
	}

	allowedSignKeyIDs := viper.GetStringSlice(keys.ValidSHA384Keys)
	ids.allowedSignKeyPublicKeys = make(map[string]asserts.PublicKey)

	for _, key := range allowedSignKeyIDs {
		restyClient := resty.New()
		req := restyClient.NewRequest()
		req.SetHeader("Accept", "application/x.ubuntu.assertion")
		resp, err := req.Get(AccountKeyGetURL + key)
		if err != nil {
			log.Error("Cannot retrieve account key assertion for: %s cannot continue, will not be able to accept devices for that key", key)
			break
		}

		accountKeyAssertion, err := asserts.Decode(resp.Body())
		if err != nil {
			log.Error("Cannot decode account key assertion for: %s cannot continue, will not be able to accept devices for that key", key)
			break
		}

		pubKey, err := asserts.DecodePublicKey(accountKeyAssertion.Body())
		if err != nil {
			log.Error("Cannot decode public key assertion for: %s cannot continue, will not be able to accept devices for that key", key)
			break
		}

		ids.allowedSignKeyPublicKeys[key] = pubKey
	}

	return ids
}

// EnrollDevice connects an IoT device with the service
func (id IdentityService) EnrollDevice(req *EnrollDeviceRequest) (*domain.Enrollment, error) {
	// Validate fields
	if req.Model.Type().Name != asserts.ModelType.Name {
		return nil, fmt.Errorf("the model assertion is an unexpected type")
	}
	if req.Serial.Type().Name != asserts.SerialType.Name {
		return nil, fmt.Errorf("the serial assertion is an unexpected type")
	}

	if req.Model.Header("brand-id") != req.Serial.Header("brand-id") {
		return nil, fmt.Errorf("the brand-id of the model and serial assertion do not match")
	}
	if req.Model.Header("model") != req.Serial.Header("model") {
		return nil, fmt.Errorf("the model name of the model and serial assertion do not match")
	}

	// Create the enrollment request
	enroll := datastore.DeviceEnrollRequest{
		Brand:        req.Model.Header("brand-id").(string),
		Model:        req.Model.Header("model").(string),
		SerialNumber: req.Serial.Header("serial").(string),
		DeviceKey:    req.Serial.Header("device-key").(string),
	}
	if req.Model.Header("store") != nil {
		enroll.StoreID = req.Model.Header("store").(string)
	}

	return id.enroll(&enroll, req.Model, req.Serial)
}

// Enroll connects an IoT device with the service
func (id IdentityService) enroll(enroll *datastore.DeviceEnrollRequest, model asserts.Assertion, serial asserts.Assertion) (*domain.Enrollment, error) {
	autoRegistrationEnabled := viper.GetBool(keys.AutoRegistrationEnabled)

	// Get the enrollment for the device, this will exist whether the device itself
	// has enrolled or not, when a device is registered this is partially created with that info
	dev, err := id.DB.DeviceGet(enroll.Brand, enroll.Model, enroll.SerialNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting device: %w", err)
	}

	if dev != nil {
		switch {
		case dev.Status == models.StatusWaiting:
			// this will result in the device being created before function returns
			break
		case dev.Status == models.StatusEnrolled:
			return nil, fmt.Errorf("`%s/%s/%s` is already enrolled", enroll.Brand, enroll.Model, enroll.SerialNumber)
		case dev.Status == models.StatusDisabled:
			return nil, fmt.Errorf("`%s/%s/%s` is disabled", enroll.Brand, enroll.Model, enroll.SerialNumber)
		}
	} else {
		if errors.Is(err, sql.ErrNoRows) && autoRegistrationEnabled {

			canAutoRegister := id.checkAutoRegistrationEligibility(model, serial)
			if !canAutoRegister {
				return nil, fmt.Errorf("`%s/%s/%s` is not eligible for auto-registration", enroll.Brand, enroll.Model, enroll.SerialNumber)
			}

			// if we couldn't find the partial enrollment (registration data) AND
			// auto-registration is enabled, then we will register the device and then enroll it
			orgID, err := id.getDefaultOrgID()
			if err != nil {
				return nil, fmt.Errorf("getting default org ID: %w", err)
			}

			register := &RegisterDeviceRequest{
				OrganizationID: orgID,
				Brand:          enroll.Brand,
				Model:          enroll.Model,
				SerialNumber:   enroll.SerialNumber,
			}

			if _, err := id.RegisterDevice(register); err != nil {
				return nil, fmt.Errorf("auto-registering device in enrollment: %w", err)
			}

			// Now that the device is registered without error, it will be created below
		}
	}

	// Enroll the device, this should happen for one of two reasons:
	// 1. a device is registered (partially enrolled) and is in the waiting state
	// 2. a device is not registered but autoregistration is enabled AND it satisfies the auto-registration criteria
	return id.DB.DeviceEnroll(*enroll)
}

func (id IdentityService) checkAutoRegistrationEligibility(model asserts.Assertion, serial asserts.Assertion) bool {
	isModelAssertionGood := id.checkKey(model)
	isSerialAssertionGood := id.checkKey(serial)

	return isModelAssertionGood && isSerialAssertionGood
}

func (id IdentityService) checkKey(model asserts.Assertion) bool {
	if pk, ok := id.allowedSignKeyPublicKeys[model.SignKeyID()]; ok {
		err := asserts.SignatureCheck(model, pk)
		if err != nil {
			log.Error("Failed signature check: ", err)
			return false
		}
	}

	return true
}

func (id IdentityService) getDefaultOrgID() (string, error) {
	orgs, err := id.OrganizationList()
	if err != nil {
		return "", fmt.Errorf("error getting OrganizationList: %w", err)
	}

	for _, org := range orgs {
		if org.Name == viper.GetString(keys.DefaultOrganization) {
			return org.ID, nil
		}
	}

	return "", fmt.Errorf("org not found: %s in %v", viper.GetString(keys.DefaultOrganization), orgs)
}

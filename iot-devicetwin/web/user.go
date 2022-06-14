package web

import (
	"encoding/json"
	"net/http"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	log "github.com/sirupsen/logrus"
)

// The snapd api message body has an action field that can be either "create" or "remove"
// Constants for these actions are defined here to not confuse them with actions in actions.go

// UserCreateAction represents the action of creating a user. This constant is used with the snapd api.
const UserCreateAction = "create"

// UserRemoveAction represents the action of removing a user. This constant is used with the snapd api.
const UserRemoveAction = "remove"

// UserAdd is the API call to create a user on a device
func (wb Service) UserAdd(w http.ResponseWriter, r *http.Request, vars varLookup) {

	orgID := vars("orgid")
	deviceID := vars("id")

	log.Tracef("UserAdd: orgid=%s device id=%s", orgID, deviceID)

	if r == nil {
		log.Error("error in json decoding for DeviceLogs: nil request")
		formatStandardResponse("UserAdd", "invalid request", w)
		return
	}

	var user messages.DeviceUser
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Error("error in json decoding for DeviceLogs: ", err)
		formatStandardResponse("UserAdd", "invalid json", w)
		return
	}

	if user.Email == "" {
		log.Error("error in JSON body: missing email field")
		formatStandardResponse("UserAdd", "invalid json", w)
		return
	}

	if user.Action != UserCreateAction {
		log.Error("error in JSON body: invalid or empty action")
		formatStandardResponse("UserAdd", "invalid json", w)
		return
	}

	err = wb.Controller.User(orgID, deviceID, user)
	if err != nil {
		log.Println("Error adding user to device:", err)
		formatStandardResponse("UserAdd", "Error adding user to device", w)
		return
	}

	formatStandardResponse("", "", w)
}

// UserRemove is the API call to remove a user on a device
func (wb Service) UserRemove(w http.ResponseWriter, r *http.Request, vars varLookup) {

	orgID := vars("orgid")
	deviceID := vars("id")

	log.Tracef("UserRemove: orgid=%s device id=%s", orgID, deviceID)

	var user messages.DeviceUser
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Error("error in json decoding for DeviceLogs: ", err)
		formatStandardResponse("UserRemove", "invalid json", w)
		return
	}

	if user.Username == "" {
		log.Error("error in JSON body: missing email field")
		formatStandardResponse("UserRemove", "invalid json", w)
		return
	}

	if user.Action != UserRemoveAction {
		log.Error("error in JSON body: invalid or empty action")
		formatStandardResponse("UserRemove", "invalid json", w)
		return
	}

	err = wb.Controller.User(orgID, deviceID, user)
	if err != nil {
		log.Println("Error removing user to device:", err)
		formatStandardResponse("UserRemove", "Error removing user to device", w)
		return
	}

	formatStandardResponse("", "", w)
}

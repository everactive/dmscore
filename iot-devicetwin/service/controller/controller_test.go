// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Device Twin Service
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

package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/everactive/dmscore/iot-devicetwin/pkg/actions"
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
	"github.com/everactive/dmscore/iot-devicetwin/service/mqtt"
	ksuid2 "github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_ActionHandler(t *testing.T) {
	m1 := []byte(`{"success": false, "message": "MOCK error"}`)
	m2 := []byte(`{"success": true, "action": "invalid"}`)

	type fields struct {
		MQTT       mqtt.Connect
		DeviceTwin devicetwin.DeviceTwin
	}
	type args struct {
		client MQTT.Client
		msg    MQTT.Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"valid", fields{&mqtt.MockConnect{}, &devicetwin.ManualMockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{}}},
		{"error-response", fields{&mqtt.MockConnect{}, &devicetwin.ManualMockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{Message: m1}}},
		{"invalid-action", fields{&mqtt.MockConnect{}, &devicetwin.ManualMockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{Message: m2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := Service{DeviceTwin: tt.fields.DeviceTwin}
			srv.ActionHandler(tt.args.msg)
		})
	}
}

func TestService_HealthHandler(t *testing.T) {
	m1 := []byte(`{"orgId": "abc", "deviceId": "aa111"}`)
	m2 := []byte(`{"orgId": "abc", "deviceId": "invalid"}`)
	m3 := []byte(`{"orgId": "abc", "deviceId": "new-device"}`)

	type fields struct {
		MQTT                    mqtt.Connect
		ReturnSoftDeletedDevice bool
	}
	type args struct {
		client         MQTT.Client
		msg            MQTT.Message
		deviceID       string
		orgID          string
		existingDevice bool
		isDeleted      bool
		waitOnChannel  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{name: "valid", args: args{msg: &mqtt.MockMessage{Message: m1}, orgID: "abc", deviceID: "aa111", existingDevice: true}},
		{name: "invalid-message", args: args{msg: &mqtt.MockMessage{}}},
		{name: "invalid-clientID", args: args{msg: &mqtt.MockMessage{Message: m2}}},
		{name: "new-clientID",
			args: args{
				msg:           &mqtt.MockMessage{Message: m3, TopicPath: "device/health/new-device"},
				deviceID:      "new-device",
				orgID:         "abc",
				waitOnChannel: true,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedDeviceMessage := &messages.Device{
				DeviceId: tt.args.deviceID,
			}
			unscopedController := &MockUnscopedController{}
			unscopedController.On("DeviceGetByID", tt.args.deviceID).Return(expectedDeviceMessage, tt.args.isDeleted, nil)

			deviceTwin := &devicetwin.MockDeviceTwin{}
			deviceTwin.On("Unscoped").Return(unscopedController)

			publishChan := make(chan mqtt.PublishMessage)
			srv := Service{DeviceTwin: deviceTwin, publishChan: publishChan}

			if !tt.args.waitOnChannel {
				if tt.args.existingDevice == true {
					deviceTwin.On("HealthHandler", messages.Health{DeviceId: tt.args.deviceID, OrgId: tt.args.orgID}).Return(nil)

					srv.HealthHandler(tt.args.msg)
					return
				}

				return
			}

			ksuid := ksuid2.New()
			generateKSUID = func() ksuid2.KSUID {
				return ksuid
			}

			act := messages.SubscribeAction{
				Action: actions.Device,
				Id:     ksuid.String(),
			}

			deviceTwin.On("HealthHandler", messages.Health{DeviceId: tt.args.deviceID, OrgId: tt.args.orgID}).Return(errors.New("some error"))
			deviceTwin.On("ActionCreate", tt.args.orgID, tt.args.deviceID, act).Return(nil)

			go func() {
				srv.HealthHandler(tt.args.msg)
			}()

			data, err := json.Marshal(&act)
			if err != nil {
				assert.Error(t, err)
			}

			expectedMessage := mqtt.PublishMessage{
				Topic:   fmt.Sprintf("devices/sub/%s", tt.args.deviceID),
				Payload: string(data),
			}

			msg := <-publishChan
			if tt.args.deviceID != "" {
				assert.Equal(t, expectedMessage.Topic, msg.Topic)
			}
			assert.Equal(t, expectedMessage.Payload, msg.Payload)

			//got := len(deviceTwin.call)
			//if got != tt.want {
			//	t.Errorf("HealthHandler() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func Test_getClientID(t *testing.T) {
	type args struct {
		msg MQTT.Message
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"valid", args{&mqtt.MockMessage{}}, "aa111"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getClientID(tt.args.msg); got != tt.want {
				t.Errorf("getClientID() = %v, want %v", got, tt.want)
			}
		})
	}
}

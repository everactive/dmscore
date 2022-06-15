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
	"testing"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
	"github.com/everactive/dmscore/iot-devicetwin/service/mqtt"
)

func TestService_SubscribeToActions(t *testing.T) {
	type fields struct {
		MQTT       mqtt.Connect
		DeviceTwin devicetwin.DeviceTwin
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"valid", fields{&mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := Service{MQTT: tt.fields.MQTT, DeviceTwin: tt.fields.DeviceTwin}
			if err := srv.SubscribeToActions(); (err != nil) != tt.wantErr {
				t.Errorf("Service.SubscribeToActions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
		{"valid", fields{&mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{}}},
		{"error-response", fields{&mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{Message: m1}}},
		{"invalid-action", fields{&mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{Message: m2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := Service{MQTT: tt.fields.MQTT, DeviceTwin: tt.fields.DeviceTwin}
			srv.ActionHandler(tt.args.client, tt.args.msg)
		})
	}
}

func TestService_HealthHandler(t *testing.T) {
	m1 := []byte(`{"orgId": "abc", "deviceId": "aa111"}`)
	m2 := []byte(`{"orgId": "abc", "deviceId": "invalid"}`)
	m3 := []byte(`{"orgId": "abc", "deviceId": "new-device"}`)

	type fields struct {
		MQTT       mqtt.Connect
		DeviceTwin *devicetwin.MockDeviceTwin
	}
	type args struct {
		client MQTT.Client
		msg    MQTT.Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"valid", fields{&mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{Message: m1}}, 0},
		{"invalid-message", fields{&mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{}}, 0},
		{"invalid-clientID", fields{&mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{Message: m2}}, 0},
		{"new-clientID", fields{&mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{ReturnSoftDeletedDevice: false}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{Message: m3, TopicPath: "device/health/new-device"}}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := Service{MQTT: tt.fields.MQTT, DeviceTwin: tt.fields.DeviceTwin}

			srv.HealthHandler(tt.args.client, tt.args.msg)
			got := len(tt.fields.DeviceTwin.Actions)
			if got != tt.want {
				t.Errorf("HealthHandler() = %v, want %v", got, tt.want)
			}
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

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

// Package mqtt is for testing MQTT clients and logic
package mqtt

import (
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const (
	mockMessageID = 1000
)

// ManualMockClient mocks the MQTT client
type ManualMockClient struct {
	open bool
}

// IsConnected mocks the connect status
func (cli *ManualMockClient) IsConnected() bool {
	return cli.open
}

// IsConnectionOpen mocks the connect status
func (cli *ManualMockClient) IsConnectionOpen() bool {
	return cli.open
}

// Connect mocks connecting to broker
func (cli *ManualMockClient) Connect() MQTT.Token {
	cli.open = true
	return &ManualMockToken{}
}

// Disconnect mocks client close
func (cli *ManualMockClient) Disconnect(quiesce uint) {
	cli.open = false
}

// Publish mocks a publish message
func (cli *ManualMockClient) Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token {
	return &ManualMockToken{}
}

// Subscribe mocks a subscribe message
func (cli *ManualMockClient) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) MQTT.Token {
	return &ManualMockToken{}
}

// SubscribeMultiple mocks subscribe messages
func (cli *ManualMockClient) SubscribeMultiple(filters map[string]byte, callback MQTT.MessageHandler) MQTT.Token {
	return &ManualMockToken{}
}

// Unsubscribe mocks a unsubscribe message
func (cli *ManualMockClient) Unsubscribe(topics ...string) MQTT.Token {
	return &ManualMockToken{}
}

// AddRoute mocks routing
func (cli *ManualMockClient) AddRoute(topic string, callback MQTT.MessageHandler) {
}

// OptionsReader mocks the options reader (badly)
func (cli *ManualMockClient) OptionsReader() MQTT.ClientOptionsReader {
	return MQTT.NewClient(nil).OptionsReader()
}

// ManualMockToken implements a Token
type ManualMockToken struct{}

// Wait mocks the token wait
func (t *ManualMockToken) Wait() bool {
	return true
}

// WaitTimeout mocks the token wait timeout
func (t *ManualMockToken) WaitTimeout(time.Duration) bool {
	return true
}

// Error mocks a token error check
func (t *ManualMockToken) Error() error {
	return nil
}

// ManualMockConnect is a mock MQTT connection
type ManualMockConnect struct{}

// Publish mocks a MQTT publish method
func (c *ManualMockConnect) Publish(topic, payload string) error {
	return nil
}

// Subscribe mocks a MQTT subscribe method
func (c *ManualMockConnect) Subscribe(topic string, callback MQTT.MessageHandler) error {
	return nil
}

// Close mocks a MQTT close method
func (c *ManualMockConnect) Close() {}

// ManualMockMessage implements an MQTT message
type ManualMockMessage struct {
	Message   []byte
	TopicPath string
}

// Duplicate mocks a duplicate message check
func (m *ManualMockMessage) Duplicate() bool {
	panic("implement me")
}

// Qos mocks the QoS flag
func (m *ManualMockMessage) Qos() byte {
	panic("implement me")
}

// Retained mocks the retained flag
func (m *ManualMockMessage) Retained() bool {
	panic("implement me")
}

// Topic mocks the topic
func (m *ManualMockMessage) Topic() string {
	if len(m.TopicPath) > 0 {
		return m.TopicPath
	}
	return "device/pub/aa111"
}

// MessageID mocks the message ID
func (m *ManualMockMessage) MessageID() uint16 {
	return mockMessageID
}

// Payload mocks the payload retrieval
func (m *ManualMockMessage) Payload() []byte {
	return m.Message
}

// Ack mocks the message ack
func (m *ManualMockMessage) Ack() {
	panic("implement me")
}

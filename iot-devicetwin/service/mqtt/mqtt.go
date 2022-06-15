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

package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/everactive/dmscore/iot-devicetwin/config"
)

// Constants for connecting to the MQTT broker
const (
	quiesce        = 250
	QOSAtMostOnce  = byte(0)
	QOSAtLeastOnce = byte(1)
	// QOSExactlyOnce = byte(2)
)

var conn *Connection
var client MQTT.Client

// Connect is the interface for an MQTT connection
type Connect interface {
	Publish(topic, payload string) error
	Subscribe(topic string, callback MQTT.MessageHandler) error
	Close()
}

// Connection for MQTT protocol
type Connection struct {
	client   MQTT.Client
	clientID string
}

// GetConnection fetches or creates an MQTT connection
func GetConnection(url, port string, connect *config.MQTTConnect, onConnect func(c MQTT.Client)) (*Connection, error) {
	if conn == nil {
		// Create the client
		localClient, err := newClient(url, port, connect, onConnect)
		if err != nil {
			return nil, err
		}

		// Create a new connection
		conn = &Connection{
			client:   localClient,
			clientID: connect.ClientID,
		}
	}

	// Check that we have a live connection
	if conn.client.IsConnectionOpen() {
		return conn, nil
	}

	// Connect to the MQTT broker
	if token := conn.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return conn, nil
}

// newClient creates a new MQTT client
func newClient(url, port string, connect *config.MQTTConnect, onConnect func(c MQTT.Client)) (MQTT.Client, error) {
	// Return the active client, if we have one
	if client != nil {
		return client, nil
	}

	// Generate a new MQTT client
	formattedURL := fmt.Sprintf("ssl://%s:%s", url, port)
	log.Println("Connect to the MQTT broker", formattedURL)

	// Generate the TLS config from the enrollment credentials
	tlsConfig, err := newTLSConfig(connect)
	if err != nil {
		return nil, err
	}

	// Set up the MQTT client options
	opts := MQTT.NewClientOptions()
	opts.AddBroker(formattedURL)
	opts.SetClientID(connect.ClientID)
	opts.SetTLSConfig(tlsConfig)
	opts.AutoReconnect = true
	opts.OnConnect = onConnect

	opts.OnConnectionLost = func(_ MQTT.Client, err error) {
		log.Printf("connection lost: %s", err)
	}

	client = MQTT.NewClient(opts)
	return client, nil
}

// newTLSConfig sets up the certificates from the enrollment record
func newTLSConfig(connect *config.MQTTConnect) (*tls.Config, error) {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(connect.RootCA)

	// Import client certificate/key pair
	cert, err := tls.X509KeyPair(connect.ClientCert, connect.ClientKey)
	if err != nil {
		return nil, err
	}

	// Create tls.Config with desired TLS properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certPool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		// nolint: gosec
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}, nil
}

// Publish sends data to the MQTT broker
func (c *Connection) Publish(topic, payload string) error {
	log.Printf("Attempting publish to %s with payload %s\n", topic, payload)
	token := c.client.Publish(topic, QOSAtLeastOnce, false, payload)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Subscribe starts a new subscription, providing a message handler for the topic
func (c *Connection) Subscribe(topic string, callback MQTT.MessageHandler) error {
	token := c.client.Subscribe(topic, QOSAtLeastOnce, callback)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Close closes the connection to the MQTT broker
func (c *Connection) Close() {
	c.client.Disconnect(quiesce)
}

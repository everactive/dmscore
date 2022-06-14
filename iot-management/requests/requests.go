// Package requests provides functionality common to manipulating REST requests
package requests

import (
	"errors"
	"strings"

	"github.com/everactive/dmscore/iot-management/config/configkey"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

// DefaultAddAuthorization provides a simple method to check that the provider is setup correctly if none is specified it is a panic, cannot run without explicit provider configuration
func DefaultAddAuthorization(_ *resty.Request) error {
	if strings.ToLower(viper.GetString(configkey.ClientTokenProvider)) != "disabled" {
		panic(errors.New("authorization token provider incorrectly configured, please set token provider"))
	}
	return nil
}

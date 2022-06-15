package controller

import (
	"testing"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
	"github.com/everactive/dmscore/iot-devicetwin/service/mqtt"
)

func TestService_User(t *testing.T) {
	type args struct {
		orgID    string
		clientID string
		user     messages.DeviceUser
	}

	tests := []struct {
		name    string
		wantErr bool
		args    args
	}{
		{
			name:    "validAddUser",
			wantErr: false,
			args: args{
				orgID:    "abc",
				clientID: "a111",
				user: messages.DeviceUser{
					Email:        "someemail",
					Sudoer:       true,
					ForceManaged: true,
				},
			},
		},
		{
			name:    "validRemoveUser",
			wantErr: false,
			args: args{
				orgID:    "abc",
				clientID: "a111",
				user: messages.DeviceUser{
					Username: "someusername",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := Service{MQTT: &mqtt.MockConnect{}, DeviceTwin: &devicetwin.MockDeviceTwin{}}
			if err := srv.User(tt.args.orgID, tt.args.clientID, &tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Service.User() test: error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

package controller

import (
	"github.com/everactive/dmscore/iot-devicetwin/service/mqtt"
	"sync"
	"testing"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
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
			publishChan := make(chan mqtt.PublishMessage)
			srv := Service{DeviceTwin: &devicetwin.ManualMockDeviceTwin{}, publishChan: publishChan}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				if err := srv.User(tt.args.orgID, tt.args.clientID, tt.args.user); (err != nil) != tt.wantErr {
					t.Errorf("Service.User() test: error = %v, wantErr %v", err, tt.wantErr)
				}
				wg.Done()
			}()

			_ = <-publishChan
			wg.Wait()
		})
	}

}

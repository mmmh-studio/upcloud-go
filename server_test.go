package upcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/mmmh-studio/upcloud-go/client"
)

func TestCreateServer(t *testing.T) {
	svc := newTestService()

	req := CreateServerRequest{
		Zone:     "de-fra1",
		Title:    "upcloud-go-test",
		Hostname: "upcloud.test.mmh.studio",
		StorageDevices: CreateServerStorageDevices{
			Devices: []CreateServerStorageDevice{
				{
					Action: "create",
					Size:   10,
					Title:  "root",
				},
			},
		},
		Interfaces: CreateServerInterfaces{
			Interfaces: CreateServerInterfacesInner{
				Interfaces: []CreateServerInterface{
					{
						IPAddresses: CreateServerIPAddresses{
							IPAddress: []CreateServerIPAddress{
								{
									Family: "IPv4",
								},
							},
						},
						Type: "public",
					},
					{
						IPAddresses: CreateServerIPAddresses{
							IPAddress: []CreateServerIPAddress{
								{
									Family: "IPv4",
								},
							},
						},
						Type: "utility",
					},
				},
			},
		},
	}

	server, err := svc.CreateServer(req)
	if err != nil {
		t.Fatal(err)
	}

	if err := teardownServer(svc, server.UUID); err != nil {
		t.Fatal(err)
	}
}

func teardownServer(svc *Service, uuid string) error {
	server, err := svc.WaitForServerState(WaitForServerStateRequest{
		Timeout:    time.Minute * 10,
		UUID:       uuid,
		WaitStates: []string{"started", "stopped"},
	})
	if err != nil {
		fmt.Printf("%#v\n", err)
		fmt.Printf("%q", err.(*client.Error).Body())
		return err
	}

	if server.State == "started" {
		if err := svc.StopServer(StopServerRequest{UUID: uuid}); err != nil {
			return err
		}

		if _, err := svc.WaitForServerState(WaitForServerStateRequest{
			Timeout:    time.Minute * 10,
			UUID:       uuid,
			WaitStates: []string{"stopped"},
		}); err != nil {
			return err
		}
	}

	if err := svc.DeleteServer(DeleteServerRequest{UUID: uuid}); err != nil {
		return err
	}

	return nil
}

package upcloud

import (
	"testing"
)

func TestCreateNetwork(t *testing.T) {
	var (
		svc = newTestService()

		name = "upcloud-go-test-network"
		zone = "de-fra1"
	)

	network, err := svc.CreateNetwork(CreateNetworkRequest{
		Name: name,
		Zone: zone,
		IPNetworks: []CreateIPNetwork{
			{
				Address: "10.1.0.0/22",
				DHCP:    true,
				Family:  "IPv4",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := network.Name, name; have != want {
		t.Errorf("have %q, want %q", have, want)
	}

	if have, want := network.Zone, zone; have != want {
		t.Errorf("have %q, want %q", have, want)
	}

	if err := teardownNetwork(svc, network.UUID); err != nil {
		t.Fatal(err)
	}
}

func TestListNetworksInZone(t *testing.T) {
	svc := newTestService()

	_, err := svc.ListNetworksInZone(ListNetworksInZoneRequest{
		Zone: "de-fra1",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func teardownNetwork(svc *Service, uuid string) error {
	return svc.DeleteNetwork(DeleteNetworkRequest{UUID: uuid})
}

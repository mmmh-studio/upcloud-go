package upcloud

import (
	"encoding/json"
	"fmt"

	"github.com/mmmh-studio/upcloud-go/client"
)

type Network struct {
	IPNetworks IPNetworks `json:"ip_networks"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	UUID       string     `json:"uuid"`
	Zone       string     `json:"zone"`
}

type IPNetworks struct {
	Networks []IPNetwork `json:"ip_network"`
}

type IPNetwork struct {
	Address          string   `json:"address"`
	DHCP             string   `json:"dhcp"`
	DHCPDefaultRoute string   `json:"dhcp_default_route"`
	DHCPNS           []string `json:"dhcp_dns"`
	Family           string   `json:"family"`
	Gateway          string   `json:"gateway"`
}

type CreateNetworkRequest struct {
	Name       string
	Zone       string
	IPNetworks []CreateIPNetwork
}

func (r CreateNetworkRequest) MarshalJSON() ([]byte, error) {
	type ipNetworks struct {
		Networks []CreateIPNetwork `json:"ip_network,omitempty"`
	}
	type inner struct {
		Name       string     `json:"name"`
		Zone       string     `json:"zone"`
		IPNetworks ipNetworks `json:"ip_networks,omitempty"`
	}

	wrapper := struct {
		Network inner `json:"network"`
	}{
		Network: inner{
			Name: r.Name,
			Zone: r.Zone,
			IPNetworks: ipNetworks{
				Networks: r.IPNetworks,
			},
		},
	}

	return json.Marshal(wrapper)
}

func (r CreateNetworkRequest) Path() string {
	return fmt.Sprintf("%s/network", apiURL)
}

func (r CreateNetworkRequest) Validate() []client.FieldError {
	es := []client.FieldError{}

	switch r.Zone {
	case "de-fra1", "fi-hel1", "fi-hel2", "nl-ams1", "sg-sin1", "uk-lon1", "us-chi1", "us-sjo1":
		// Zone is supproted.
	default:
		es = append(es, client.FieldError{
			Name:        "Zone",
			Description: fmt.Sprintf("'%s' not supported", r.Zone),
		})
	}

	return es
}

type CreateIPNetwork struct {
	Address string
	DHCP    bool
	Family  string
}

func (r CreateIPNetwork) MarshalJSON() ([]byte, error) {
	wrapper := struct {
		Address string `json:"address"`
		DHCP    string `json:"dhcp"`
		Family  string `json:"family"`
	}{
		Address: r.Address,
		DHCP:    "no",
		Family:  r.Family,
	}

	if r.DHCP {
		wrapper.DHCP = "yes"
	}

	return json.Marshal(wrapper)
}

type createNetworkResponse struct {
	Network Network
}

func (r *createNetworkResponse) UnmarshalJSON(raw []byte) error {
	var f struct {
		Network Network `json:"network"`
	}

	if err := json.Unmarshal(raw, &f); err != nil {
		return err
	}

	r.Network = f.Network

	return nil
}

type DeleteNetworkRequest struct {
	UUID string
}

func (r DeleteNetworkRequest) Path() string {
	return fmt.Sprintf("%s/network/%s", apiURL, r.UUID)
}

type GetNetworkDetailsRequest struct {
	UUID string
}

func (r GetNetworkDetailsRequest) Path() string {
	return fmt.Sprintf("%s/network/%s", apiURL, r.UUID)
}

type getNetworkDetailsResponse struct {
	network Network
}

func (r *getNetworkDetailsResponse) UnmarshalJSON(raw []byte) error {
	var res struct {
		Network Network `json:"network"`
	}

	if err := json.Unmarshal(raw, &res); err != nil {
		return err
	}

	r.network = res.Network

	return nil
}

type ListNetworksInZoneRequest struct {
	Zone string
}

func (r ListNetworksInZoneRequest) Path() string {
	return fmt.Sprintf("%s/network?zone=%s", apiURL, r.Zone)
}

type listNetworksInZoneResponse struct {
	Networks []Network
}

func (r *listNetworksInZoneResponse) UnmarshalJSON(raw []byte) error {
	var res struct {
		Wrapper struct {
			Networks []Network `json:"network"`
		} `json:"networks"`
	}

	if err := json.Unmarshal(raw, &res); err != nil {
		return err
	}

	r.Networks = res.Wrapper.Networks

	return nil
}

func (s *Service) CreateNetwork(req CreateNetworkRequest) (*Network, error) {
	res := &createNetworkResponse{}

	if err := s.client.Post(req, res); err != nil {
		return nil, err
	}

	return &res.Network, nil
}

func (s *Service) GetNetworkDetails(req GetNetworkDetailsRequest) (*Network, error) {
	res := &getNetworkDetailsResponse{}

	if err := s.client.Get(req, res); err != nil {
		return nil, err
	}

	return &res.network, nil
}

func (s *Service) DeleteNetwork(req DeleteNetworkRequest) error {
	return s.client.Delete(req)
}

func (s *Service) ListNetworksInZone(req ListNetworksInZoneRequest) ([]Network, error) {
	res := &listNetworksInZoneResponse{}

	if err := s.client.Get(req, res); err != nil {
		return nil, err
	}

	return res.Networks, nil
}

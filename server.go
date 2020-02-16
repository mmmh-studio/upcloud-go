package upcloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mmmh-studio/upcloud-go/client"
)

type Server struct {
	State string `json:"state"`
	UUID  string `json:"uuid"`
}

type CreateServerRequest struct {
	// The zone in which the server will be hosted, e.g. fi-hel1. See Zones.
	//
	// A valid zone identifierA valid zone identifier.
	Zone string `json:"zone"`
	// A short, informational description.
	//
	// 0-64 characters
	Title string `json:"title"`
	// A valid domain name, e.g. host.example.com. The maximum length is 128 characters.
	Hostname string `json:"hostname"`
	// The storage_devices block contains storage_device blocks that define the
	// attached storages.
	//
	// 1-8 storage_device blocks
	StorageDevices CreateServerStorageDevices `json:"storage_devices"`
	// All interfaces wanted for the server.
	//
	// An array of 1-10 interface objects.
	Interfaces CreateServerInterfaces `json:"networking"`
}

func (r CreateServerRequest) MarshalJSON() ([]byte, error) {
	type inner struct {
		Zone           string                     `json:"zone"`
		Title          string                     `json:"title"`
		Hostname       string                     `json:"hostname"`
		StorageDevices CreateServerStorageDevices `json:"storage_devices"`
		Interfaces     CreateServerInterfaces     `json:"networking"`
	}

	wrapper := struct {
		Server inner `json:"server"`
	}{
		Server: inner{
			Zone:           r.Zone,
			Title:          r.Title,
			Hostname:       r.Hostname,
			StorageDevices: r.StorageDevices,
			Interfaces:     r.Interfaces,
		},
	}

	return json.Marshal(wrapper)
}

func (r CreateServerRequest) Path() string {
	return fmt.Sprintf("%s/server", apiURL)
}

func (r CreateServerRequest) Validate() []client.FieldError {
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

type CreateServerInterfaces struct {
	Interfaces CreateServerInterfacesInner `json:"interfaces"`
}

type CreateServerInterfacesInner struct {
	Interfaces []CreateServerInterface `json:"interface"`
}

type CreateServerInterface struct {
	IPAddresses CreateServerIPAddresses `json:"ip_addresses"`
	Type        string                  `json:"type"`
	Network     string                  `json:"network,omitempty"`
}

type CreateServerIPAddresses struct {
	IPAddress []CreateServerIPAddress `json:"ip_address"`
}

type CreateServerIPAddress struct {
	Family string `json:"family"`
}

type CreateServerStorageDevices struct {
	Devices []CreateServerStorageDevice `json:"storage_device"`
}

type CreateServerStorageDevice struct {
	// The method used to create or attach the specified storage.
	//
	// create / clone / attach
	Action string `json:"action"`
	// if action is create The size of the storage device in gigabytes. This attribute is applicable only if action is create or clone.
	//
	// 10-1024
	Size int `json:"size"`
	// if action is clone or attach The UUID of the storage device to be attached or cloned. Applicable only if action is attach or clone.
	//
	// A valid storage UUID
	Storage string `json:"storage,omitempty"`
	// A short, informational description for the storage.
	//
	// 0-64 characters
	Title string `json:"title"`
}

type createServerResponse struct {
	Server Server
}

func (r *createServerResponse) UnmarshalJSON(raw []byte) error {
	var f struct {
		Server Server `json:"server"`
	}

	if err := json.Unmarshal(raw, &f); err != nil {
		return err
	}

	r.Server = f.Server

	return nil
}

type DeleteServerRequest struct {
	UUID string
}

func (r DeleteServerRequest) Path() string {
	return fmt.Sprintf("%s/server/%s?storage=1", apiURL, r.UUID)
}

type GetServerDetailsRequest struct {
	UUID string
}

func (r GetServerDetailsRequest) Path() string {
	return fmt.Sprintf("%s/server/%s", apiURL, r.UUID)
}

type getServerDetailsResponse struct {
	Server Server
}

func (r *getServerDetailsResponse) UnmarshalJSON(raw []byte) error {
	var res struct {
		Server Server `json:"server"`
	}

	if err := json.Unmarshal(raw, &res); err != nil {
		return err
	}

	r.Server = res.Server

	return nil
}

type StopServerRequest struct {
	Timeout time.Duration
	Type    string
	UUID    string
}

func (r StopServerRequest) MarshalJSON() ([]byte, error) {
	type stopServer struct {
		StopType string `json:"stop_type"`
		Timeout  string `json:"timeout"`
	}

	wrapper := struct {
		StopServer stopServer `json:"stop_server"`
	}{
		StopServer: stopServer{
			StopType: "hard",
			Timeout:  "60",
		},
	}

	return json.Marshal(wrapper)
}

func (r StopServerRequest) Path() string {
	return fmt.Sprintf("%s/server/%s/stop", apiURL, r.UUID)
}

func (r StopServerRequest) Validate() []client.FieldError {
	return nil
}

type WaitForServerStateRequest struct {
	Timeout    time.Duration
	UUID       string
	WaitStates []string
}

func (s *Service) CreateServer(req CreateServerRequest) (*Server, error) {
	res := createServerResponse{}

	if err := s.client.Post(req, &res); err != nil {
		return nil, err
	}

	return &res.Server, nil
}

func (s *Service) DeleteServer(req DeleteServerRequest) error {
	return s.client.Delete(req)
}

func (s *Service) GetServerDetails(req GetServerDetailsRequest) (*Server, error) {
	res := &getServerDetailsResponse{}

	if err := s.client.Get(req, res); err != nil {
		return nil, err
	}

	return &res.Server, nil
}

func (s *Service) StopServer(req StopServerRequest) error {
	return s.client.Post(req, nil)
}

func (s *Service) WaitForServerState(req WaitForServerStateRequest) (*Server, error) {
	var (
		attempts      = 0
		sleepDuration = time.Second * 3
	)

	for {
		attempts++
		time.Sleep(sleepDuration)

		server, err := s.GetServerDetails(GetServerDetailsRequest{
			UUID: req.UUID,
		})
		if err != nil {
			return nil, err
		}

		for _, state := range req.WaitStates {
			if server.State == state {
				return server, nil
			}
		}

		if time.Duration(attempts)*sleepDuration >= req.Timeout {
			return nil, fmt.Errorf("timeout reached waitin for state change of server '%s' to %v", req.UUID, req.WaitStates)
		}
	}
}

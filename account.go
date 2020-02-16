package upcloud

import (
	"encoding/json"
	"fmt"
)

type Account struct {
	Username string `json:"username"`
}

type GetAccountRequest struct{}

func (r GetAccountRequest) Path() string {
	return fmt.Sprintf("%s/account", apiURL)
}

type getAccountResponse struct {
	Account Account
}

func (r *getAccountResponse) UnmarshalJSON(raw []byte) error {
	var res struct {
		Account Account `json:"account"`
	}

	if err := json.Unmarshal(raw, &res); err != nil {
		return err
	}

	r.Account = res.Account

	return nil
}

func (s *Service) GetAccount() (*Account, error) {
	res := &getAccountResponse{}

	if err := s.client.Get(GetAccountRequest{}, res); err != nil {
		return nil, err
	}

	return &res.Account, nil
}

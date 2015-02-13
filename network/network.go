package network

import (
	"encoding/json"

	"github.com/gertd/go-openstack/identity"
	"github.com/parnurzeal/gorequest"
)

type networkResp struct {
	Networks []Network `json:"networks"`
}

type Network struct {
	Id                  string   `json:"id"`
	Name                string   `json:"name"`
	Status              string   `json:"status"`
	Subnets             []string `json:"subnets"`
	TenantId            string   `json:"tenant_id"`
	RouterExternal      bool     `json:"router:external"`
	AdminStateUp        bool     `json:"admin_state_up"`
	Shared              bool     `json:"shared"`
	PortSecurityEnabled bool     `json:"port_security_enabled"`
}

func GetNetworks(url string, token identity.Token) (networks []Network, err error) {

	req := gorequest.New()

	_, body, errs := req.Get(url+"/v2.0/networks").
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", token.Id).
		End()

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	var nw = networkResp{}
	if err = json.Unmarshal([]byte(body), &nw); err != nil {
		return
	}

	networks = nw.Networks
	err = nil
	return
}

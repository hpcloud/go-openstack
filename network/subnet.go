// subnet.go
package network

import (
	"encoding/json"
	"fmt"

	"github.com/gertd/go-openstack/identity"
	"github.com/parnurzeal/gorequest"
)

type subnetsResp struct {
	Subnets []Subnet `json:"subnets"`
}

type Subnet struct {
	Id              string           `json:"id"`
	Name            string           `json:"name"`
	NetworkId       string           `json:"network_id"`
	TenantId        string           `json:"tenant_id"`
	EnableDHCP      bool             `json:"enable_dhcp"`
	DnsNameserver   []string         `json:"dns_nameservers"`
	AllocationPools []AllocationPool `json:"allocation_pools"`
	HostRoutes      []string         `json:"host_routes"`
	IPVersion       int              `json:"ip_version"`
	GatewayIP       string           `json:"gateway_ip"`
	CIDR            string           `json:"cidr"`
}

type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func GetSubnets(auth identity.Auth) (subnets []Subnet, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/v2.0/subnets",
		auth.EndpointList["network"])

	_, body, errs := req.Get(url).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", auth.Access.Token.Id).
		End()

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	var sn = subnetsResp{}
	if err = json.Unmarshal([]byte(body), &sn); err != nil {
		return
	}

	subnets = sn.Subnets
	err = nil
	return
}

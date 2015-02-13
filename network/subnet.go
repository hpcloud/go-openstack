package network

import (
	"encoding/json"
	"fmt"

	"github.com/gertd/go-openstack/identity"
	"github.com/parnurzeal/gorequest"
)

type subnetsResp struct {
	Subnets []Subnet `json:"subnets,omitempty"`
}

type Subnet struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	NetworkId  string `json:"network_id,omitempty"`
	TenantId   string `json:"tenant_id,omitempty"`
	EnableDHCP bool   `json:"enable_dhcp,omitempty"`
	//DnsNameserver   []string         `json:"dns_nameservers,omitempty"`
	//AllocationPools []AllocationPool `json:"allocation_pools,omitempty"`
	//HostRoutes      []string         `json:"host_routes,omitempty"`
	IPVersion int    `json:"ip_version,omitempty"`
	GatewayIP string `json:"gateway_ip,omitempty"`
	CIDR      string `json:"cidr,omitempty"`
}

type AllocationPool struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
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
		fmt.Println(err)
		return
	}

	subnets = sn.Subnets
	err = nil
	return
}

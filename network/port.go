// port.go
package network

import (
	"encoding/json"
	"fmt"

	"github.com/gertd/go-openstack/identity"
	"github.com/parnurzeal/gorequest"
)

type portsResp struct {
	Ports []Port `json:"ports"`
}

type portResp struct {
	Port Port `json:"port"`
}

type portReq struct {
	Port Port `json:"port,omitempty"`
}

type ByName []Port

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type Port struct {
	Id                  string    `json:"id,omitempty"`
	Name                string    `json:"name,omitempty"`
	Status              string    `json:"status,omitempty"`
	AdminStateUp        bool      `json:"admin_state_up,omitempty"`
	PortSecurityEnabled bool      `json:"port_security_enabled,omitempty"`
	DeviceId            string    `json:"device_id,omitempty"`
	DeviceOwner         string    `json:"device_owner,omitempty"`
	NetworkId           string    `json:"network_id,omitempty"`
	TenantId            string    `json:"tenant_id,omitempty"`
	MacAddress          string    `json:"mac_address,omitempty"`
	FixedIPs            []FixedIP `json:"fixed_ips,omitempty"`
	SecurityGroups      []string  `json:"security_groups,omitempty"`
}

type FixedIP struct {
	SubnetId  string `json:"subnet_id,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
}

func GetPorts(auth identity.Auth) (ports []Port, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/v2.0/ports",
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

	var p = portsResp{}
	if err = json.Unmarshal([]byte(body), &p); err != nil {
		return
	}

	ports = p.Ports
	err = nil
	return
}

func GetPort(url string, token identity.Token, id string) (port Port, err error) {

	req := gorequest.New()

	reqUrl := fmt.Sprintf("%s/v2.0/ports/%s", url, id)

	_, body, errs := req.Get(reqUrl).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", token.Id).
		End()

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	if err = json.Unmarshal([]byte(body), &port); err != nil {
		return
	}

	err = nil
	return
}

func DeletePort(auth identity.Auth, id string) (err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/v2.0/ports/%s",
		auth.EndpointList["network"],
		id)

	_, _, errs := req.Delete(url).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", auth.Access.Token.Id).
		End()

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	err = nil
	return
}

func CreatePort(auth identity.Auth, port Port) (resPort Port, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/v2.0/ports",
		auth.EndpointList["network"])

	reqPort := portReq{port}

	b, err := json.Marshal(reqPort)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, body, errs := req.Post(url).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", auth.Access.Token.Id).
		Send(string(b)).
		End()

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	var portResp = portResp{}
	if err = json.Unmarshal([]byte(body), &portResp); err != nil {
		return
	}

	err = nil
	resPort = portResp.Port
	return
}

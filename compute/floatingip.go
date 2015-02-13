package compute

import (
	"encoding/json"
	"fmt"

	"github.com/gertd/go-openstack"
	"github.com/gertd/go-openstack/identity"

	"github.com/parnurzeal/gorequest"
)

type floatingIPResp struct {
	FloatingIP FloatingIP `json:"floating_ip"`
}

type floatingIPsResp struct {
	FloatingIPs []FloatingIP `json:"floating_ips"`
}

type FloatingIP struct {
	Id         string `json:"id"`
	IP         string `json:"ip"`
	InstanceId string `json:"instance_id"`
	FixedIP    string `json:"fixed_ip"`
	Pool       string `json:"pool"`
}

func GetFloatingIPs(auth identity.Auth) (floating_ips []FloatingIP, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/os-floating-ips",
		auth.EndpointList["compute"])

	resp, body, errs := req.Get(url).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", auth.Access.Token.Id).
		End()

	if err = openstack.CheckHttpResponseStatusCode(resp.StatusCode); err != nil {
		return
	}

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	var r = floatingIPsResp{}
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return
	}

	floating_ips = r.FloatingIPs
	err = nil
	return
}

func CreateFloatingIP(auth identity.Auth) (floatingIP FloatingIP, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/os-floating-ips",
		auth.EndpointList["compute"])

	var reqBody = `{"pool": "Ext-Net"}`

	resp, body, errs := req.Post(url+"/os-floating-ips").
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", auth.Access.Token.Id).
		Send(reqBody).
		End()

	if err = openstack.CheckHttpResponseStatusCode(resp.StatusCode); err != nil {
		return
	}

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	var r = floatingIPResp{}
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return
	}

	floatingIP = r.FloatingIP
	err = nil
	return
}

func DeleteFloatingIP(auth identity.Auth, floatingIP FloatingIP) (err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/os-floating-ips/%s",
		auth.EndpointList["compute"],
		floatingIP.Id)

	resp, _, errs := req.Delete(url).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", auth.Access.Token.Id).
		End()

	if err = openstack.CheckHttpResponseStatusCode(resp.StatusCode); err != nil {
		return
	}

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	err = nil
	return
}

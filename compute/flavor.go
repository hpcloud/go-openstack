package compute

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gertd/go-openstack"
	"github.com/gertd/go-openstack/identity"

	"github.com/parnurzeal/gorequest"
)

type flavorsResp struct {
	Flavors []Flavor `json:"flavors"`
}

type flavorsDetailResp struct {
	Flavors []FlavorDetail `json:"flavors"`
}

type Flavor struct {
	Id    string           `json:"id"`
	Name  string           `json:"name"`
	Links []openstack.Link `json:"links"`
}

type FlavorDetail struct {
	Id         string           `json:"id"`
	Name       string           `json:"name"`
	Links      []openstack.Link `json:"links"`
	Ram        int              `json:"ram"`
	VCPUs      int              `json:"vcpus"`
	Swap       string           `json:"swap"`
	RxtxFactor interface{}      `json:"rxtx_factor"`
	Ephemeral  interface{}      `json:"OS-FLV-EXT-DATA:ephemeral"`
	Disk       int              `json:"disk"`
}

func GetFlavors(auth identity.Auth) (flavors []Flavor, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/flavors",
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

	var r = flavorsResp{}
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return
	}

	flavors = r.Flavors
	err = nil
	return
}

func GetFlavorsDetail(auth identity.Auth) (flavors []FlavorDetail, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/flavors/detail",
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

	var r = flavorsDetailResp{}
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return
	}

	flavors = r.Flavors
	err = nil
	return
}

func GetFlavor(auth identity.Auth, name string) (flavor Flavor, err error) {

	flavors, err := GetFlavors(auth)
	if err != nil {
		return
	}

	for _, v := range flavors {
		if v.Name == name {
			flavor = v
			err = nil
			return
		}
	}

	err = errors.New(fmt.Sprintf("flavor %s not found", name))
	return
}

func GetFlavorDetail(auth identity.Auth, name string) (flavor FlavorDetail, err error) {

	flavors, err := GetFlavorsDetail(auth)
	if err != nil {
		return
	}

	for _, v := range flavors {
		if v.Name == name {
			flavor = v
			err = nil
			return
		}
	}

	err = errors.New(fmt.Sprintf("flavor %s not found", name))
	return
}

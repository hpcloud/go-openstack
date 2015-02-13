package compute

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gertd/go-openstack"
	"github.com/gertd/go-openstack/identity"

	"github.com/parnurzeal/gorequest"
)

type serversResp struct {
	Servers []ServerInfo `json:"servers"`
}

type serverDetailResp struct {
	Servers []Server `json:"servers"`
}

type serverResp struct {
	Server Server `json:"server"`
}

type ServerInfo struct {
	Id    string           `json:"id"`
	Name  string           `json:"name"`
	Links []openstack.Link `json:"links"`
}

type Server struct {
	Id               string               `json:"id"`
	Name             string               `json:"name"`
	Status           string               `json:"status"`
	Created          interface{}          `json:"created"`
	Updated          interface{}          `json:"updated"`
	HostId           string               `json:"hostId"`
	Addresses        map[string][]Address `json:"addresses"`
	Links            []openstack.Link     `json:"links"`
	Image            interface{}          `json:"image"`
	Flavor           interface{}          `json:"flavor"`
	TaskState        string               `json:"OS-EXT-STS:task_state"`
	VMState          string               `json:"OS-EXT-STS:vm_state"`
	PowerState       int                  `json:"OS-EXT-STS:power_state"`
	AvailabilityZone string               `json:"OS-EXT-AZ:availability_zone:"`
	UserId           string               `json:"user_id"`
	TenantId         string               `json:"tenant_id"`
	AccessIPv4       string               `json:"accessIPv4"`
	AccessIPv6       string               `json:"accessIPv6"`
	ConfigDrive      string               `json:"config_drive"`
	Progress         int                  `json:"progress"`
	MetaData         map[string]string    `json:"metadata"`
	AdminPass        string               `json:"adminPass"`
}

type serverReq struct {
	Server NewServer `json:"server"`
}

type NewServer struct {
	Name           string          `json:"name,omitempty"`
	ImageRef       string          `json:"imageRef,omitempty"`
	KeyName        string          `json:"key_name,omitempty"`
	FlavorRef      string          `json:"flavorRef,omitempty"`
	MinCount       int             `json:"min_count,omitempty"`
	MaxCount       int             `json:"max_count,omitempty"`
	UserData       string          `json:"user_data,omitempty"`
	Network        []Network       `json:"networks,omitempty"`
	SecurityGroups []SecurityGroup `json:"security_groups,omitempty"`
}

type Network struct {
	Uuid string `json:"uuid,omitempty"`
	Port string `json:"port,omitempty"`
}

type SecurityGroups struct {
	SecurityGroups []SecurityGroup `json:"security_groups"`
}

type SecurityGroup struct {
	Name string `json:"name,omitempty"`
}

type Address struct {
	Addr    string `json:"addr"`
	Version int    `json:"version"`
	Type    string `json:"OS-EXT-IPS:type"`
	MacAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
}

type ByName []ServerInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func GetServers(auth identity.Auth) (servers []ServerInfo, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/servers",
		auth.EndpointList["compute"])

	_, body, errs := req.Get(url).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", auth.Access.Token.Id).
		End()

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	var sr = serversResp{}
	if err = json.Unmarshal([]byte(body), &sr); err != nil {
		return
	}

	servers = sr.Servers
	err = nil
	return
}

func GetServersDetail(auth identity.Auth) (servers []Server, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/servers/detail",
		auth.EndpointList["compute"])

	_, body, errs := req.Get(url).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", auth.Access.Token.Id).
		End()

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	var r = serverDetailResp{}
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return
	}

	servers = r.Servers
	err = nil
	return
}

func GetServer(auth identity.Auth, id string) (server Server, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/servers/%s",
		auth.EndpointList["compute"],
		id)

	resp, respBody, errs := req.Get(url).
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

	r := serverResp{}
	if err = json.Unmarshal([]byte(respBody), &r); err != nil {
		return
	}

	server = r.Server
	err = nil
	return
}

func DeleteServer(auth identity.Auth, id string) (err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/servers/%s",
		auth.EndpointList["compute"],
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

func CreateServer(auth identity.Auth, newServer NewServer) (server Server, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/servers",
		auth.EndpointList["compute"])

	r := serverReq{newServer}

	b, err := json.Marshal(r)
	if err != nil {
		return
	}

	resp, body, errs := req.Post(url).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", auth.Access.Token.Id).
		Send(string(b)).
		End()

	if !(resp.StatusCode == 201 || resp.StatusCode == 202) {
		err = errors.New("Error: status code != 201 or 202, object not created (" + resp.Status + ")")
		return
	}

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	serverResp := serverResp{}
	if err = json.Unmarshal([]byte(body), &serverResp); err != nil {
		return
	}

	server = serverResp.Server
	err = nil
	return
}

func ServerAction(auth identity.Auth, id string, action string, key string, value string) (err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/servers/%s/action",
		auth.EndpointList["compute"],
		id)

	var reqBody = fmt.Sprintf(`
	{
		"%s":
		{
			"%s": "%s"
		}
	}`, action, key, value)

	resp, _, errs := req.Post(url).
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

	err = nil
	return
}

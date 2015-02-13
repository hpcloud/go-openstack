package identity

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gertd/go-openstack"
	"github.com/parnurzeal/gorequest"
)

// request
type AuthenticationReq struct {
	Auth AuthenticateReq `json:"auth,omitempty"`
}

type AuthenticateReq struct {
	TenantName          string              `json:"tenantName,omitempty"`
	TenantId            string              `json:"tenantId,omitempty"`
	PasswordCredentials PasswordCredentials `json:"passwordCredentials,omitempty"`
	//TokenCredentials    TokenCredentials    `json:"token,omitempty"`
}

type PasswordCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type TokenCredentials struct {
	Id string `json:"id,omitempty"`
}

// response
type Auth struct {
	Access       Access `json:"access"`
	EndpointList map[string]string
}

type Access struct {
	Token          Token     `json:"token"`
	User           User      `json:"user"`
	ServiceCatalog []Service `json:"serviceCatalog"`
}

type Token struct {
	Id      string    `json:"id"`
	Expires time.Time `json:"expires"`
	Tenant  Tenant    `json:"tenant"`
}

type Tenant struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	Id          string
	Name        string
	Roles       []Role
	Roles_links []string
}

type Role struct {
	Id       string
	Name     string
	TenantId string
}

type Service struct {
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	Endpoints       []Endpoint `json:"endpoints"`
	Endpoints_links []string
}

type Endpoint struct {
	TenantId    string
	PublicURL   string
	InternalURL string
	Region      string
	VersionId   string
	VersionInfo string
	VersionList string
}

func Authenticate(openStackConfig openstack.OpenStackConfig) (auth Auth, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/tokens", openStackConfig.AuthUrl)

	authReq := AuthenticationReq{}

	if len(openStackConfig.TenantId) > 0 {
		authReq.Auth.TenantId = openStackConfig.TenantId
	} else if len(openStackConfig.TenantName) > 0 {
		authReq.Auth.TenantName = openStackConfig.TenantName
	}

	authReq.Auth.PasswordCredentials = PasswordCredentials{openStackConfig.Username, openStackConfig.Password}

	b, err := json.Marshal(authReq)
	if err != nil {
		return
	}

	_, body, errs := req.Post(url).
		Set(`Accept-Encoding`, `gzip,deflate`).
		Set(`Accept`, `application/json`).
		Set(`Content-Type`, `application/json`).
		Send(string(b)).
		End()

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	if err = json.Unmarshal([]byte(body), &auth); err != nil {
		return
	}

	if !auth.Access.Token.Expires.After(time.Now()) {
		err = errors.New("Error: The AuthN token is expired")
		return
	}

	auth.EndpointList = auth.endpointList(openStackConfig.Region)

	err = nil
	return
}

func (auth Auth) endpointList(region string) (list map[string]string) {

	list = make(map[string]string)

	for _, v := range auth.Access.ServiceCatalog {
		for _, endPoint := range v.Endpoints {
			if len(region) == 0 || (len(region) > 0 && endPoint.Region == region) {
				list[v.Type] = endPoint.PublicURL
			}
		}
	}

	return list
}

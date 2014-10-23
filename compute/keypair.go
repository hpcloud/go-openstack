// keypair.go
package compute

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gertd/go-openstack"
	"github.com/gertd/go-openstack/identity"

	"github.com/parnurzeal/gorequest"
)

type keyPairsResp struct {
	KeyPairs []Keypair `josn:"keypairs"`
}

type Keypair struct {
	KeyPair KeyPairEntry `json:"keypair"`
}

type KeyPairEntry struct {
	Name        string `json:"name"`
	PublicKey   string `json:"public_key"`
	FingerPrint string `json:"fingerprint"`
}

type keyPairDetailResp struct {
	KeyPair KeyPairDetail `json:"keypair"`
}

type KeyPairDetail struct {
	Name        string      `json:"name"`
	PublicKey   string      `json:"public_key"`
	FingerPrint string      `json:"fingerprint"`
	UserId      string      `json:"user_id"`
	CreatedAt   interface{} `json:"created_at"`
	UpdatedAt   interface{} `json:"updated_at"`
	DeletedAt   interface{} `json:"deleted_at"`
	IsDeleted   bool        `json:"deleted"`
	Id          int         `json:"id"`
}

func GetKeypairs(auth identity.Auth) (keypairs []Keypair, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/os-keypairs",
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

	var kp = keyPairsResp{}
	if err = json.Unmarshal([]byte(body), &kp); err != nil {
		return
	}

	keypairs = kp.KeyPairs
	err = nil
	return
}

func GetKeypair(auth identity.Auth, name string) (keypair KeyPairDetail, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/os-keypairs/%s",
		auth.EndpointList["compute"],
		name)

	resp, body, errs := req.Get(url).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", auth.Access.Token.Id).
		End()

	if resp.StatusCode == 404 {
		err = errors.New(fmt.Sprintf("keypair %s not found", name))
		return
	}

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	var r = keyPairDetailResp{}
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return
	}

	keypair = r.KeyPair
	err = nil
	return
}

// image.go
package image

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gertd/go-openstack"
	"github.com/gertd/go-openstack/identity"

	"github.com/parnurzeal/gorequest"
)

type imagesResp struct {
	Images []Image `json:"images"`
}

type imagesDetailResp struct {
	Images []ImageDetail `json:"images"`
}

type Image struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	ContainerFormat string `json:"container_format"`
	DiskFormat      string `json:"disk_format"`
	CheckSum        string `json:"checksum"`
	Size            int    `json:"size"`
}

type ImageDetail struct {
	Id              string            `json:"id"`
	Name            string            `json:"name"`
	Status          string            `json:"status"`
	Owner           string            `json:"owner"`
	Created         interface{}       `json:"created_at,omitempty"`
	Updated         interface{}       `json:"updated_at,omitempty"`
	Deleted         interface{}       `json:"deleted_at,omitempty"`
	IsDeleted       bool              `json:"deleted"`
	ContainerFormat string            `json:"container_format"`
	DiskFormat      string            `json:"disk_format"`
	CheckSum        string            `json:"checksum"`
	Size            int               `json:"size"`
	Protected       bool              `json:"protected"`
	IsPublic        bool              `json:"is_public"`
	MinDisk         int               `json:"minDisk"`
	MinRam          int               `json:"minRam"`
	Properties      map[string]string `json:"properties"`
}

func GetImages(url string, token identity.Token) (images []Image, err error) {

	req := gorequest.New()

	resp, body, errs := req.Get(url+"/images").
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", token.Id).
		End()

	if err = openstack.CheckHttpResponseStatusCode(resp.StatusCode); err != nil {
		return
	}

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	var r = imagesResp{}
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return
	}

	images = r.Images
	err = nil
	return
}

func GetImagesDetail(url string, token identity.Token) (images []ImageDetail, err error) {

	err = nil
	return
}

func GetImage(auth identity.Auth, name string) (image Image, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/images?limit=20&name=%s",
		auth.EndpointList["image"],
		name)

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

	var r = imagesResp{}
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return
	}

	if len(r.Images) == 0 {
		err = errors.New(fmt.Sprintf("image %s not found", name))
		return
	}
	if len(r.Images) > 1 {
		err = errors.New(fmt.Sprintf("image %s multiple entries found", name))
		return
	}

	image = r.Images[0]

	err = nil
	return
}

func GetImageDetail(auth identity.Auth, name string) (image ImageDetail, err error) {

	req := gorequest.New()

	url := fmt.Sprintf("%s/images/detail?limit=20&name=%s",
		auth.EndpointList["image"],
		name)

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

	var r = imagesDetailResp{}
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return
	}

	if len(r.Images) == 0 {
		err = errors.New(fmt.Sprintf("image %s not found", name))
		return
	}
	if len(r.Images) > 1 {
		err = errors.New(fmt.Sprintf("image %s multiple entries found", name))
		return
	}

	image = r.Images[0]

	err = nil
	return
}

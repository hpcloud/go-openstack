// openstack.go
package openstack

import (
	"errors"
	"log"
	"os"
)

type OpenStackConfig struct {
	AuthUrl    string
	TenantId   string
	TenantName string
	Username   string
	Password   string
	Region     string
}

type Link struct {
	HRef string `json:"href"`
	Rel  string `json:"rel"`
}

func (config OpenStackConfig) Log() {
	log.Printf("%-20s - %s\n", "OS_AUTH_URL", config.AuthUrl)
	log.Printf("%-20s - %s\n", "OS_TENANT_ID", config.TenantId)
	log.Printf("%-20s - %s\n", "OS_TENANT_NAME", config.TenantName)
	log.Printf("%-20s - %s\n", "OS_USERNAME", config.Username)
	log.Printf("%-20s - %s\n", "OS_REGION_NAME", config.Region)

}

func InitializeFromEnv() (config OpenStackConfig, err error) {

	var c = OpenStackConfig{}

	c.AuthUrl = os.Getenv("OS_AUTH_URL")
	c.TenantId = os.Getenv("OS_TENANT_ID")
	c.TenantName = os.Getenv("OS_TENANT_NAME")
	c.Username = os.Getenv("OS_USERNAME")
	c.Password = os.Getenv("OS_PASSWORD")
	c.Region = os.Getenv("OS_REGION_NAME")

	if len(c.AuthUrl) == 0 {
		err = errors.New("Error: no authentication URL specified")
		return
	}
	if len(c.Username) == 0 {
		err = errors.New("Error: no username specified")
		return
	}
	if len(c.Password) == 0 {
		err = errors.New("Error: no password specified")
		return
	}
	if len(c.TenantName) == 0 {
		err = errors.New("Error: no tenant name specified")
		return
	}
	if len(c.TenantId) == 0 {
		err = errors.New("Error: no tenant ID specified")
		return
	}

	config = c
	err = nil
	return
}

//utility methods
func CheckHttpResponseStatusCode(statusCode int) error {
	switch statusCode {
	case 200, 201, 202, 204:
		return nil
	case 400:
		return errors.New("Error: response == 400 bad request")
	case 401:
		return errors.New("Error: response == 401 unauthorised")
	case 403:
		return errors.New("Error: response == 403 forbidden")
	case 404:
		return errors.New("Error: response == 404 not found")
	case 405:
		return errors.New("Error: response == 405 method not allowed")
	case 409:
		return errors.New("Error: response == 409 conflict")
	case 413:
		return errors.New("Error: response == 413 over limit")
	case 415:
		return errors.New("Error: response == 415 bad media type")
	case 422:
		return errors.New("Error: response == 422 unprocessable")
	case 429:
		return errors.New("Error: response == 429 too many request")
	case 500:
		return errors.New("Error: response == 500 instance fault / server err")
	case 501:
		return errors.New("Error: response == 501 not implemented")
	case 503:
		return errors.New("Error: response == 503 service unavailable")
	}
	return errors.New("Error: unexpected response status code")
}

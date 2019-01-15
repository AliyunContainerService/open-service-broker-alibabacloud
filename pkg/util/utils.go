package util

import (
	"encoding/json"
	"fmt"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/denverdino/aliyungo/metadata"
	"github.com/golang/glog"
)

// WriteResponse will serialize 'object' to the HTTP ResponseWriter
// using the 'code' as the HTTP status code
func WriteResponse(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

// WriteErrorResponse will serialize 'err' to the HTTP ResponseWriter
// with JSON formatted error response
// using the 'code' as the HTTP status code
func WriteErrorResponse(w http.ResponseWriter, code int, err error) {
	WriteResponse(w, code, &brokerapi.BrokerErrorResponse{
		Error:	     err.Error(),
		Description: err.Error(),
	})
}

// BodyToObject will convert the incoming HTTP request into the
// passed in 'object'
func BodyToObject(r *http.Request, object interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		return err
	}

	return nil
}

// ResponseBodyToObject will reading the HTTP response into the
// passed in 'object'
func ResponseBodyToObject(r *http.Response, object interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	glog.Info(string(body))

	err = json.Unmarshal(body, object)
	if err != nil {
		return err
	}

	return nil
}

// ExecCmd executes a command and returns the stdout + error, if any
func ExecCmd(cmd string) (string, error) {
	fmt.Println("command: " + cmd)

	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]

	out, err := exec.Command(head, parts...).CombinedOutput()
	if err != nil {
		fmt.Printf("Command failed with: %s\n", err)
		fmt.Printf("Output: %s\n", out)
		return "", err
	}
	return string(out), nil
}

// Fetch will do an HTTP GET to the passed in URL, returning
// HTTP Body of the response or any error
func Fetch(u string) (string, error) {
	fmt.Printf("Fetching: %s\n", u)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// FetchObject will do an HTTP GET to the passed in URL, returning
// the response parsed as a JSON object, or any error
func FetchObject(u string, object interface{}) error {
	r, err := http.Get(u)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		return err
	}
	return nil
}

type AccessMetaData struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
}

func GetAccessMetaData() (*AccessMetaData, error) {
	m := metadata.NewMetaData(nil)

	roleName := ""
	var err error
	if roleName, err = m.RoleName(); err != nil {
		glog.Errorf("Get role name error: %v", err.Error())
		return nil, err
	}
	role, err := m.RamRoleToken(roleName)
	if err != nil {
		glog.Errorf("Get STS Token error: %v", err.Error())
		return nil, err
	}

	accessMetaData := AccessMetaData{
		AccessKeyId:     role.AccessKeyId,
		AccessKeySecret: role.AccessKeySecret,
		SecurityToken:   role.SecurityToken,
	}
	return &accessMetaData, nil
}

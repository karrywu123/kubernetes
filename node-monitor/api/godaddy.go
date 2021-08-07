package api

import (
	"net/http"

	//"net/url"
	//"crypto/md5"
	"io/ioutil"

	//"fmt"
	//"strconv"
	"encoding/json"
	// "errors"
	//"bytes"
	//"reflect"
	// "strings"
)

type List_alldomains_res struct {
	CreatedAt           string `json:"createdAt"`
	DeletedAt           string `json:"deletedAt"`
	Domain              string `json:"domain"`
	DomainId            int64  `json:"domainId"`
	ExpirationProtected bool   `json:"expirationProtected"`
	Expires             string `json:"expires"`
	HoldRegistrar       bool   `json:"holdRegistrar"`
	Locked              bool   `json:"locked"`
	NameServers         string `json:"nameServers"`
	Privacy             bool   `json:"privacy"`
	RenewAuto           bool   `json:"renewAuto"`
	Renewable           bool   `json:"renewable"`
	Status              string `json:"status"`
	TransferProtected   bool   `json:"transferProtected"`
}

//common.Cfg.Godaddy_api.Head_Accept
//common.Cfg.Godaddy_api.Head_Type
//common.Cfg.Godaddy_api.Authorization

//获取godaddy账号所有账号
func List_alldomains(auth string) (records []List_alldomains_res, err error) {
	var values map[string]string
	values = make(map[string]string)

	// var records []List_alldomains_res

	client := &http.Client{}
	urls := "https://api.godaddy.com/v1/domains?limit=1000"

	//values := url.Values{"X-Auth-Email":{common.Cfg.Cloudflareapi.Email},"X-Auth-Key":{common.Cfg.Cloudflareapi.Auth_api_key},"X-Auth-User-Service-Key":{common.Cfg.Cloudflareapi.User_service_api_key},"Content-Type":{common.Cfg.Cloudflareapi.Content_type}}
	values["Accept"] = "application/json"
	values["Content-Type"] = "application/json"
	values["Authorization"] = auth

	req, _ := http.NewRequest("GET", urls, nil)

	for key, value := range values {
		req.Header.Add(key, value)
	}

	//resp, err := http.PostForm(urls,values)
	resp, err := client.Do(req)

	//fmt.Println(resp.StatusCode)

	if err == nil {
		res, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println(string(res))

		err = json.Unmarshal([]byte(res), &records)
	} else {
		return records, err
	}
	defer resp.Body.Close()

	return
}

// //获取某个域名详细信息
// func List_domain_detail(auth string,domain string)(res []byte, err error) {
//     var values map[string]string
//     values = make(map[string]string)
//     client := &http.Client{}
//     urls := "https://api.godaddy.com/v1/domains/" + domain

//     //values := url.Values{"X-Auth-Email":{common.Cfg.Cloudflareapi.Email},"X-Auth-Key":{common.Cfg.Cloudflareapi.Auth_api_key},"X-Auth-User-Service-Key":{common.Cfg.Cloudflareapi.User_service_api_key},"Content-Type":{common.Cfg.Cloudflareapi.Content_type}}
//     values["Accept"] = common.Cfg.Godaddy_api.Head_Accept
//     values["Content-Type"] = common.Cfg.Godaddy_api.Head_Type
//     values["Authorization"] = auth

//     req,_ := http.NewRequest("GET",urls,nil)

//     for key,value := range values {
//         req.Header.Add(key,value)
//     }

//     //resp, err := http.PostForm(urls,values)
//     resp, err := client.Do(req)
//     //fmt.Println(resp.StatusCode)

//     if err == nil {
//         res,err = ioutil.ReadAll(resp.Body)
//         resp.Body.Close()
//     }

//     return
// }

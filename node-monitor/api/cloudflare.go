package api

import (
	"net/http"

	//"net/url"
	//"crypto/md5"
	"io/ioutil"
	"node-monitor/common"

	//"fmt"
	//"strconv"
	"encoding/json"
	"errors"

	//"bytes"
	//"reflect"
	"strings"

	"node-monitor/libs/logs"
)

//zone返回相关结构体
type Res_cloudflare_zone struct {
	Success     bool                   `json:"success"`
	Errors      []string               `json:"errors"`
	Messages    []string               `json:"messages"`
	Result      []Result_value_zone    `json:"result"`
	Result_info Result_info_value_zone `json:"result_info"`
}
type Result_value_zone struct {
	Id                    string      `json:"id"`
	Name                  string      `json:"name"`
	Development_mode      int         `json:"development_mode"`
	Original_name_servers [2]string   `json:"original_name_servers"`
	Original_registrar    string      `json:"original_registrar"`
	Original_dnshost      string      `json:"original_dnshost"`
	Created_on            string      `json:"created_on"`
	Modified_on           string      `json:"modified_on"`
	Activated_on          string      `json:"activated_on"`
	Owner                 interface{} `json:"owner"`
	Account               interface{} `json:"account"`
	Permissions           []string    `json:"permissions"`
	Plan                  interface{} `json:"plan"`
	Plan_pending          interface{} `json:"plan_pending"`
	Status                string      `json:"status"`
	Paused                bool        `json:"paused"`
	Type                  string      `json:"type"`
	Name_servers          [2]string   `json:"name_servers"`
}
type Result_info_value_zone struct {
	Page        int `json:"page"`
	Per_page    int `json:"per_page"`
	Count       int `json:"count"`
	Total_count int `json:"total_count"`
}

//记录值相关结构体
type Res_cloudflare_record struct {
	Success     bool                     `json:"success"`
	Errors      interface{}              `json:"errors"`
	Messages    interface{}              `json:"messages"`
	Result      []Result_value_record    `json:"result"`
	Result_info Result_info_value_record `json:"result_info"`
}
type Res_cloudflare_record_add struct {
	Success  bool                `json:"success"`
	Errors   interface{}         `json:"errors"`
	Messages interface{}         `json:"messages"`
	Result   Result_value_record `json:"result"`
}
type Res_cloudflare_record_del struct {
	Success  bool        `json:"success"`
	Errors   interface{} `json:"errors"`
	Messages interface{} `json:"messages"`
	Result   interface{} `json:"result"`
}
type Result_value_record struct {
	Id          string      `json:"id"`
	Type        string      `json:"type"`
	Name        string      `json:"name"`
	Content     string      `json:"content"`
	Proxiable   bool        `json:"proxiable"`
	Proxied     bool        `json:"proxied"`
	Ttl         int         `json:"ttl"`
	Locked      bool        `json:"locked"`
	Zone_id     string      `json:"zone_id"`
	Zone_name   string      `json:"zone_name"`
	Created_on  string      `json:"created_on"`
	Modified_on string      `json:"modified_on"`
	Data        interface{} `json:"data"`
}
type Result_info_value_record struct {
	Page        int `json:"page"`
	Per_page    int `json:"per_page"`
	Count       int `json:"count"`
	Total_count int `json:"total_count"`
}

type Res_cloudflare_firewall struct {
	Success  bool        `json:"success"`
	Errors   interface{} `json:"errors"`
	Messages interface{} `json:"messages"`
	Result   interface{} `json:"result"`
}

//列出全部域名的zone_id
func List_zone_id() (res []byte, err error) {
	var values map[string]string
	values = make(map[string]string)
	client := &http.Client{}
	urls := "https://api.cloudflare.com/client/v4/zones?page=0&per_page=1000"

	//values := url.Values{"X-Auth-Email":{common.Cfg.Cloudflareapi.Email},"X-Auth-Key":{common.Cfg.Cloudflareapi.Auth_api_key},"X-Auth-User-Service-Key":{common.Cfg.Cloudflareapi.User_service_api_key},"Content-Type":{common.Cfg.Cloudflareapi.Content_type}}
	values["X-Auth-Email"] = common.Cfg.Cloudflareapi.Email
	values["X-Auth-Key"] = common.Cfg.Cloudflareapi.Auth_api_key
	values["X-Auth-User-Service-Key"] = common.Cfg.Cloudflareapi.User_service_api_key
	values["Content-Type"] = common.Cfg.Cloudflareapi.Content_type

	req, _ := http.NewRequest("GET", urls, nil)

	for key, value := range values {
		req.Header.Add(key, value)
	}

	//resp, err := http.PostForm(urls,values)
	resp, err := client.Do(req)

	//fmt.Println(resp.StatusCode)

	if err == nil {
		res, err = ioutil.ReadAll(resp.Body)

	} else {
		logs.Error("Failed to request,err,", err)
		return
	}
	defer resp.Body.Close()

	return
}

//获取某个域名的zone_id
func Get_domain_zone_id(domain string) (zone_id string, err error) {
	res, err := List_zone_id()
	flag := 0
	if err == nil {
		var records Res_cloudflare_zone
		if err := json.Unmarshal([]byte(res), &records); err == nil {
			if records.Success {
				for _, v := range records.Result {
					//fmt.Println(v.Id,v.Name)
					if domain == v.Name {
						zone_id = v.Id
						flag = 1
						break
					}
				}

			} else {
				flag = 1
			}
		} else {
			flag = 1
		}
	} else {
		flag = 1
	}

	if flag == 0 {
		err = errors.New("No have domain: " + domain + " in CF.")
	}

	return
}

//清理域名缓存
func Purge_domain_cache(domain string) (status int, err error) {
	type res_cache struct {
		Success  bool        `json:"success"`
		Errors   interface{} `json:"errors"`
		Messages interface{} `json:"messages"`
		Result   interface{} `json:"result"`
	}

	status = 0

	var records res_cache
	var values map[string]string
	var values1 map[string]bool
	values = make(map[string]string)
	values1 = make(map[string]bool)

	client := &http.Client{}
	zone_id, err := Get_domain_zone_id(domain)

	if err == nil {

		urls := "https://api.cloudflare.com/client/v4/zones/" + zone_id + "/purge_cache"

		//values := url.Values{"X-Auth-Email":{common.Cfg.Cloudflareapi.Email},"X-Auth-Key":{common.Cfg.Cloudflareapi.Auth_api_key},"X-Auth-User-Service-Key":{common.Cfg.Cloudflareapi.User_service_api_key},"Content-Type":{common.Cfg.Cloudflareapi.Content_type}}
		values["X-Auth-Email"] = common.Cfg.Cloudflareapi.Email
		values["X-Auth-Key"] = common.Cfg.Cloudflareapi.Auth_api_key
		values["X-Auth-User-Service-Key"] = common.Cfg.Cloudflareapi.User_service_api_key
		values["Content-Type"] = common.Cfg.Cloudflareapi.Content_type

		values1["purge_everything"] = true
		js, _ := json.Marshal(values1)

		req, _ := http.NewRequest("POST", urls, strings.NewReader(string(js)))

		for key, value := range values {
			req.Header.Add(key, value)
		}

		//resp, err := http.PostForm(urls,values)
		resp, err := client.Do(req)

		//fmt.Println(resp.StatusCode)

		if err == nil {
			res, _ := ioutil.ReadAll(resp.Body)
			// resp.Body.Close()
			if err := json.Unmarshal([]byte(res), &records); err == nil {
				if records.Success {
					status = 1
				}
			}

		} else {
			logs.Error("Failed to request,err,", err)
			return status, err
		}
		defer resp.Body.Close()
	}
	return
}

//获取某域名所有记录以及id
func List_records_id(domain string) (res []byte, err error) {
	var values map[string]string
	values = make(map[string]string)
	client := &http.Client{}
	zone_id, err := Get_domain_zone_id(domain)
	if err == nil {
		urls := "https://api.cloudflare.com/client/v4/zones/" + zone_id + "/dns_records?page=0&per_page=1000"

		//values := url.Values{"X-Auth-Email":{common.Cfg.Cloudflareapi.Email},"X-Auth-Key":{common.Cfg.Cloudflareapi.Auth_api_key},"X-Auth-User-Service-Key":{common.Cfg.Cloudflareapi.User_service_api_key},"Content-Type":{common.Cfg.Cloudflareapi.Content_type}}
		values["X-Auth-Email"] = common.Cfg.Cloudflareapi.Email
		values["X-Auth-Key"] = common.Cfg.Cloudflareapi.Auth_api_key
		values["X-Auth-User-Service-Key"] = common.Cfg.Cloudflareapi.User_service_api_key
		values["Content-Type"] = common.Cfg.Cloudflareapi.Content_type

		req, _ := http.NewRequest("GET", urls, nil)

		for key, value := range values {
			req.Header.Add(key, value)
		}

		//resp, err := http.PostForm(urls,values)
		resp, err := client.Do(req)

		//fmt.Println(resp.StatusCode)

		if err == nil {
			res, err = ioutil.ReadAll(resp.Body)
			// resp.Body.Close()
		} else {
			logs.Error("Failed to request,err,", err)
			return res, err
		}
		defer resp.Body.Close()
	}
	return
}

//获取某域名某记录值id
func Get_record_id(domain string, record string, record_value string) (record_id string, err error) {
	res, err := List_records_id(domain)
	flag := 0
	if err == nil {
		var records Res_cloudflare_record
		if err := json.Unmarshal([]byte(res), &records); err == nil {
			// fmt.Println(records)
			if records.Success {
				for _, v := range records.Result {
					//fmt.Println(strings.Split(v.Name,"." + domain)[0])
					if strings.Split(v.Name, "."+domain)[0] == record {
						if record_value == v.Content {
							flag = 1
							record_id = v.Id
						}
					}
				}

			} else {
				flag = 1
			}
		} else {
			flag = 1
			//fmt.Println(err)
		}
	} else {
		flag = 1
	}

	if flag == 0 {
		err = errors.New("No have record: " + record + " in domain: " + domain)
	}

	return
}

//添加域名记录
func Add_record_cf(domain string, typeof string, record string, record_value string) (status int, err error) {
	var values map[string]string
	values = make(map[string]string)

	var values1 map[string]interface{}
	values1 = make(map[string]interface{})

	status = 0

	client := &http.Client{}
	zone_id, err := Get_domain_zone_id(domain)
	if err == nil {
		urls := "https://api.cloudflare.com/client/v4/zones/" + zone_id + "/dns_records"

		//values := url.Values{"X-Auth-Email":{common.Cfg.Cloudflareapi.Email},"X-Auth-Key":{common.Cfg.Cloudflareapi.Auth_api_key},"X-Auth-User-Service-Key":{common.Cfg.Cloudflareapi.User_service_api_key},"Content-Type":{common.Cfg.Cloudflareapi.Content_type}}
		values["X-Auth-Email"] = common.Cfg.Cloudflareapi.Email
		values["X-Auth-Key"] = common.Cfg.Cloudflareapi.Auth_api_key
		values["X-Auth-User-Service-Key"] = common.Cfg.Cloudflareapi.User_service_api_key
		values["Content-Type"] = common.Cfg.Cloudflareapi.Content_type

		values1["type"] = typeof
		if record == "@" {
			values1["name"] = domain
		} else {
			values1["name"] = record + "." + domain
		}
		values1["content"] = record_value
		values1["ttl"] = 1
		values1["priority"] = 0
		values1["proxied"] = false
		js, _ := json.Marshal(values1)

		req, _ := http.NewRequest("POST", urls, strings.NewReader(string(js)))

		for key, value := range values {
			req.Header.Add(key, value)
		}

		//resp, err := http.PostForm(urls,values)
		resp, err := client.Do(req)

		//fmt.Println(resp.StatusCode)

		if err == nil {
			res, _ := ioutil.ReadAll(resp.Body)
			// resp.Body.Close()
			var records Res_cloudflare_record_add
			if err := json.Unmarshal([]byte(res), &records); err == nil {
				//fmt.Println(records)
				if records.Success {
					status = 1
				}
			}
		} else {
			logs.Error("Failed to request,err,", err)
			return status, err
		}
		defer resp.Body.Close()
	}
	return
}

//删除域名记录
func Del_record_cf(domain string, record string, record_value string) (status int, err error) {
	var values map[string]string
	values = make(map[string]string)

	status = 0

	client := &http.Client{}
	zone_id, err := Get_domain_zone_id(domain)
	record_id, err1 := Get_record_id(domain, record, record_value)
	if err == nil {
		if err1 == nil {
			urls := "https://api.cloudflare.com/client/v4/zones/" + zone_id + "/dns_records/" + record_id

			//values := url.Values{"X-Auth-Email":{common.Cfg.Cloudflareapi.Email},"X-Auth-Key":{common.Cfg.Cloudflareapi.Auth_api_key},"X-Auth-User-Service-Key":{common.Cfg.Cloudflareapi.User_service_api_key},"Content-Type":{common.Cfg.Cloudflareapi.Content_type}}
			values["X-Auth-Email"] = common.Cfg.Cloudflareapi.Email
			values["X-Auth-Key"] = common.Cfg.Cloudflareapi.Auth_api_key
			values["X-Auth-User-Service-Key"] = common.Cfg.Cloudflareapi.User_service_api_key
			values["Content-Type"] = common.Cfg.Cloudflareapi.Content_type

			req, _ := http.NewRequest("DELETE", urls, nil)

			for key, value := range values {
				req.Header.Add(key, value)
			}

			//resp, err := http.PostForm(urls,values)
			resp, err := client.Do(req)

			//fmt.Println(resp.StatusCode)

			if err == nil {
				res, _ := ioutil.ReadAll(resp.Body)
				// resp.Body.Close()
				var records Res_cloudflare_record_del
				if err := json.Unmarshal([]byte(res), &records); err == nil {
					//fmt.Println(records)
					if records.Success {
						status = 1
					}
				}
			} else {
				logs.Error("Failed to request,err,", err)
				return status, err
			}
			defer resp.Body.Close()
		}
	}
	return
}

//根据规则名字获取规则详细信息
// func GetCfFirewall(zoneid string, ruleDestrition string) (err error) {
// 	var values map[string]string

// 	values = make(map[string]string)
// 	client := &http.Client{}

// 	urls := "https://api.cloudflare.com/client/v4/zones/" + zoneid + "/firewall/rules"

// 	values["X-Auth-Email"] = common.Cfg.Cloudflareapi.Email
// 	values["X-Auth-Key"] = common.Cfg.Cloudflareapi.Auth_api_key
// 	values["X-Auth-User-Service-Key"] = common.Cfg.Cloudflareapi.User_service_api_key
// 	values["Content-Type"] = common.Cfg.Cloudflareapi.Content_type
// 	req, _ := http.NewRequest("GET", urls, nil)

// 	for key, value := range values {
// 		req.Header.Add(key, value)
// 	}

// 	//resp, err := http.PostForm(urls,values)
// 	resp, err := client.Do(req)
// 	//fmt.Println(resp.StatusCode)

// 	if err == nil {
// 		res, err := ioutil.ReadAll(resp.Body)

// 		flag := 0
// 		if err == nil {
// 			var records Res_cloudflare_zone
// 			if err := json.Unmarshal([]byte(res), &records); err == nil {
// 				if records.Success {
// 					for _, v := range records.Result {
// 						//fmt.Println(v.Id,v.Name)
// 						if domain == v.Name {
// 							zone_id = v.Id
// 							flag = 1
// 							break
// 						}
// 					}

// 				} else {
// 					flag = 1
// 				}
// 			} else {
// 				flag = 1
// 			}
// 		} else {
// 			flag = 1
// 		}

// 		if flag == 0 {
// 			err = errors.New("No have domain: " + domain + " in CF.")
// 		}
// 		resp.Body.Close()
// 	}

// 	return
// }

//设置自定义防火墙规则,js challenge
func SetCfFirewall(cfdomain string, defenddomain string, filterId string, expression string) (status int, err error) {
	var values map[string]string
	values = make(map[string]string)

	var values1 map[string]interface{}
	values1 = make(map[string]interface{})
	// var values2 map[string]interface{}
	// values2 = make(map[string]interface{})

	var values4 [1]interface{}

	client := &http.Client{}

	zoneId, err := Get_domain_zone_id(cfdomain)
	if err != nil {
		logs.Error("Failed to get zoneid for domain:", cfdomain, ",err,", err)
	} else {

		// urls := "https://api.cloudflare.com/client/v4/zones/" + zoneId + "/firewall/rules"
		urls := "https://api.cloudflare.com/client/v4/zones/" + zoneId + "/filters"
		values["X-Auth-Email"] = common.Cfg.Cloudflareapi.Email
		values["X-Auth-Key"] = common.Cfg.Cloudflareapi.Auth_api_key
		values["X-Auth-User-Service-Key"] = common.Cfg.Cloudflareapi.User_service_api_key
		values["Content-Type"] = common.Cfg.Cloudflareapi.Content_type

		values1["id"] = filterId
		values1["paused"] = false
		// values1["description"] = ruleDestrition
		values1["expression"] = expression

		values4[0] = values1

		js, _ := json.Marshal(values4)
		req, _ := http.NewRequest("PUT", urls, strings.NewReader(string(js)))
		// fmt.Println(string(js))

		for key, value := range values {
			req.Header.Add(key, value)
		}

		//resp, err := http.PostForm(urls,values)
		resp, err := client.Do(req)

		//fmt.Println(resp.StatusCode)
		status = 0
		if err == nil {
			res, _ := ioutil.ReadAll(resp.Body)
			// fmt.Println(string(res))

			var records Res_cloudflare_firewall
			if err := json.Unmarshal([]byte(res), &records); err == nil {
				//fmt.Println(records)
				if records.Success {
					status = 1
				}
			} else {
				logs.Error("Failed to parse results from CF,err,", err)
			}

			// resp.Body.Close()
		} else {
			logs.Error("Failed to request,err,", err)
			return status, err
		}
		defer resp.Body.Close()
	}

	return
}

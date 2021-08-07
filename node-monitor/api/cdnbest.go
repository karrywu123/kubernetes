package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	//"reflect"
	"strings"

	"node-monitor/libs/logs"
)

type Cdnbest_res struct {
	Status struct {
		Code       int    `json:"code"`
		Created_at string `json:"created_at"`
		Message    string `json:"message"`
	} `json:"status"`
	Rows []struct {
		Uid   string      `json:"uid"`
		Ngid  string      `json:"ngid"`
		Vhost string      `json:"vhost"`
		Name  string      `json:"name"`
		Id    string      `json:"id"`
		Value interface{} `json:"value"`
	} `json:"rows"`
}

func interface2string(inter interface{}) string {
	tempStr := ""
	switch inter.(type) {
	case string:
		tempStr = inter.(string)
		break
	case float64:
		tempStr = strconv.FormatFloat(inter.(float64), 'f', -1, 64)
		break
	case int64:
		tempStr = strconv.FormatInt(inter.(int64), 10)
		break
	case int:
		tempStr = strconv.Itoa(inter.(int))
		break
	}
	return tempStr
}

// type Res_status struct {
//     Code                int                  `json:"code"`
//     Created_at          string               `json:"created_at"`
//     Message             string               `json:"message"`
// }

// type Res_rows struct {
//     Uid                 string               `json:"uid"`
//     Ngid                string               `json:"ngid"`
//     Vhost               string               `json:"vhost"`
//     Name                string               `json:"name"`
//     Id                  string               `json:"id"`
//     Value               interface{}        `json:"value"`
// }

//获取站点配置
func Get_settings_site(site string) (res []byte, err error) {
	urls := "http://ip/api/?c=site&a=getVhostSettingList"
	data := []byte("123456getVhostSettingList14863uLmFAs4sqc6p5xAe")
	s := fmt.Sprintf("%x", md5.Sum(data))
	values := url.Values{"r": {"123456"}, "s": {s}, "uid": {"14863"}, "vhost": {site}, "name": {site}}
	resp, err := http.PostForm(urls, values)
	if err == nil {
		res, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}

	return
}

//增加301重写
func Add_urlrewrite(site string, code string, domain string) (status int, err error) {
	var res_status Cdnbest_res
	status = 0
	urls := "http://ip/api/?c=site&a=addUrlRewrtite"
	data := []byte("123456addUrlRewrtite14863uLmFAs4sqc6p5xAe")
	s := fmt.Sprintf("%x", md5.Sum(data))
	values := url.Values{"r": {"123456"}, "s": {s}, "uid": {"14863"}, "vhost": {site}, "name": {site}, "code": {code}, "url": {"http://" + domain + "(.*)$"}, "target": {"https://" + domain + "$1"}}
	resp, err := http.PostForm(urls, values)

	if err == nil {
		res, _ := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal([]byte(res), &res_status); err == nil {
			if res_status.Status.Code == 1 {
				status = 1
			}
		}
	} else {
		logs.Error("Failed to request,err,", err)
		return
	}
	defer resp.Body.Close()

	return
}

//删除301重写
func Del_urlrewrite(site string, domain string, domain_urlwrite_id int) (status int, err error) {
	var res_status Cdnbest_res
	status = 0
	urls := "http://ip/api/?c=site&a=deleteConfig"
	data := []byte("123456deleteConfig14863uLmFAs4sqc6p5xAe")
	s := fmt.Sprintf("%x", md5.Sum(data))
	values := url.Values{"r": {"123456"}, "s": {s}, "uid": {"14863"}, "vhost": {site}, "name": {site}, "type": {"redirect"}, "id": {strconv.Itoa(domain_urlwrite_id)}}
	resp, err := http.PostForm(urls, values)

	if err == nil {
		res, _ := ioutil.ReadAll(resp.Body)
		// resp.Body.Close()
		if err := json.Unmarshal([]byte(res), &res_status); err == nil {
			if res_status.Status.Code == 1 {
				status = 1
			}
		}
	} else {
		logs.Error("Failed to request,err,", err)
		return
	}
	defer resp.Body.Close()

	return

}

//判断某域名是否添加了301重写
func Judge_domain_urlwrite(site string, domain string) (status int, err error) {
	var settings Cdnbest_res
	flag := 0
	res, err := Get_settings_site(site)
	if err == nil {
		if err := json.Unmarshal([]byte(res), &settings); err == nil {
			for _, value := range settings.Rows {
				if value.Name == "redirect" {
					v := value.Value.(map[string]interface{})
					tmp := strings.Split(interface2string(v["host"]), "/")[2]
					get_domain := strings.Split(tmp, "(")[0]
					if get_domain == domain {
						flag = 1
						break
					}
				}
			}
		}
	}
	if flag == 1 {
		status = 1
	} else {
		status = 0
	}

	return
}

//获取域名301重写id
func Get_domain_urlwrite_id(site string, domain string) (id string, err error) {
	var settings Cdnbest_res
	res, err := Get_settings_site(site)
	if err == nil {
		if err := json.Unmarshal([]byte(res), &settings); err == nil {
			for _, value := range settings.Rows {
				if value.Name == "redirect" {
					v := value.Value.(map[string]interface{})
					tmp := strings.Split(interface2string(v["host"]), "/")[2]
					get_domain := strings.Split(tmp, "(")[0]
					if get_domain == domain {
						id = value.Id
						break
					}
				}
			}
		}
	}
	return
}

//判断站点是否配置证书
func Judge_site_https(site string) (status int, err error) {
	var settings Cdnbest_res
	status = 0
	res, err := Get_settings_site(site)
	if err == nil {
		if err := json.Unmarshal([]byte(res), &settings); err == nil {
			for _, value := range settings.Rows {
				if value.Name == "https" {
					v := value.Value.(map[string]interface{})
					cert := interface2string(v["certificate"])
					if cert != "" {
						status = 1
						break
					}
				}
			}
		}
	}

	return

}

package api

import (
	"net/http"
	"net/url"

	//"crypto/md5"
	"io/ioutil"
	"node-monitor/common"

	//"fmt"
	//"strconv"
	"encoding/json"
	//"reflect"
	//"strings"
)

type Dnspod_res struct {
	Status  Res_status    `json:"status"`
	Domain  Res_domain    `json:"domain"`
	Info    Res_info      `json:"info"`
	Records []Res_records `json:"records"`
}

type Res_status struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Created_at string `json:"created_at"`
}

type Res_domain struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Punycode   string   `json:"punycode"`
	Grade      string   `json:"grade"`
	Owner      string   `json:"owner"`
	Ext_status string   `json:"ext_status"`
	Ttl        int      `json:"ttl"`
	Min_ttl    int      `json:"min_ttl"`
	Dnspod_ns  []string `json:"dnspod_ns"`
}

type Res_info struct {
	Sub_domains  string `json:"sub_domains"`
	Record_total string `json:"record_total"`
}

type Res_records struct {
	Id             string      `json:"id"`
	Ttl            string      `json:"ttl"`
	Value          string      `json:"value"`
	Enabled        string      `json:"enabled"`
	Status         string      `json:"status"`
	Updated_on     string      `json:"updated_on"`
	Name           string      `json:"name"`
	Line           string      `json:"line"`
	Line_id        string      `json:"line_id"`
	Type           string      `json:"type"`
	Weight         interface{} `json:"weight"`
	Monitor_status string      `json:"monitor_status"`
	Remark         string      `json:"remark"`
	Use_aqb        string      `json:"use_aqb"`
	Mx             string      `json:"mx"`
	Hold           string      `json:"hold"`
}

type Res_Domain struct {
	Status Res_status `json:"status"`
	// Info    Res_domain_info `json:"info"`
	Domain Res_domain_info `json:"domain"`
}
type Res_domain_info struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Punycode          string `json:"punycode"`
	Grade             string `json:"grade"`
	Grade_title       string `json:"grade_title"`
	Status            string `json:"status"`
	Ext_status        string `json:"ext_status"`
	Records           string `json:"records"`
	Group_id          string `json:"group_id"`
	Is_mark           string `json:"is_mark"`
	Remark            string `json:"remark"`
	Is_vip            string `json:"is_vip"`
	Searchengine_push string `json:"searchengine_push"`
	User_id           string `json:"user_id"`
	Created_on        string `json:"created_on"`
	Updated_on        string `json:"updated_on"`
	Ttl               string `json:"ttl"`
	Cname_speedup     string `json:"cname_speedup"`
	Owner             string `json:"owner"`
	Vip_start_at      string `json:"vip_start_at"`
	Vip_end_at        string `json:"vip_end_at"`
	Vip_auto_renew    string `json:"vip_auto_renew"`
}

// type Res_domain_info struct {
// 	Domain_total    int    `json:"domain_total"`
// 	All_total       int    `json:"all_total"`
// 	Mine_total      int    `json:"mine_total"`
// 	Share_total     string `json:"share_total"`
// 	Vip_total       int    `json:"vip_total"`
// 	Ismark_total    int    `json:"ismark_total"`
// 	Pause_total     int    `json:"pause_total"`
// 	Error_total     int    `json:"error_total"`
// 	Lock_total      int    `json:"lock_total"`
// 	Spam_total      int    `json:"spam_total"`
// 	Vip_expire      int    `json:"vip_expire"`
// 	Share_out_total int    `json:"share_out_total"`
// }
// type Res_domains struct {
// 	Id                int64  `json:"id"`
// 	Status            string `json:"status"`
// 	Grade             string `json:"grade"`
// 	Group_id          string `json:"group_id"`
// 	Searchengine_push string `json:"searchengine_push"`
// 	Is_mark           string `json: "is_mark"`
// 	Ttl               string `json:"ttl"`
// 	Cname_speedup     string `json: "cname_speedup"`
// 	Remark            string `json:"remark"`
// 	Created_on        string `json:"created_on"`
// 	Updated_on        string `json:"updated_on"`
// 	Punycode          string `json:"punycode"`
// 	Ext_status        string `json:"ext_status"`
// 	Name              string `json:"name"`
// 	Grade_title       string `json: "grade_title"`
// 	Is_vip            string `json: "is_vip"`
// 	Owner             string `json:"owner"`
// 	Records           string `json: "records"`
// 	Auth_to_anquanbao bool   `json:"auth_to_anquanbao"`
// }

//增加记录
func Add_record(domain string, sub_domain string, record_type string, value string, status0 string) (status int, err error) {
	var res_status Dnspod_res
	status = 0
	urls := "https://dnsapi.cn/Record.Create"
	values := url.Values{"login_token": {common.Cfg.Dnspod.Login_token}, "format": {common.Cfg.Dnspod.Format}, "lang": {common.Cfg.Dnspod.Lang}, "error_on_empty": {common.Cfg.Dnspod.Error_on_empty}, "domain": {domain}, "sub_domain": {sub_domain}, "record_type": {record_type}, "value": {value}, "record_line": {"默认"}, "status": {status0}}
	resp, err := http.PostForm(urls, values)

	if err == nil {
		res, _ := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal([]byte(res), &res_status); err == nil {
			if res_status.Status.Code == "1" {
				status = 1
			}
		}
	} else {
		return status, err
	}
	defer resp.Body.Close()

	return
}

//删除记录
func Del_record(domain string, record string, value string) (status int, err error) {
	var res_status Dnspod_res
	status = 0
	urls := "https://dnsapi.cn/Record.Remove"
	id, err := Select_record_id(domain, record, value)
	if err == nil {
		values := url.Values{"login_token": {common.Cfg.Dnspod.Login_token}, "format": {common.Cfg.Dnspod.Format}, "lang": {common.Cfg.Dnspod.Lang}, "error_on_empty": {common.Cfg.Dnspod.Error_on_empty}, "domain": {domain}, "record_id": {id}}
		resp, err := http.PostForm(urls, values)

		if err == nil {
			res, _ := ioutil.ReadAll(resp.Body)

			if err := json.Unmarshal([]byte(res), &res_status); err == nil {
				if res_status.Status.Code == "1" {
					status = 1
				}
			}
		} else {
			return status, err
		}
		defer resp.Body.Close()
	}

	return
}

//查询站点所有记录信息
func Select_record_info(domain string) (res []byte, err error) {
	urls := "https://dnsapi.cn/Record.List"
	values := url.Values{"login_token": {common.Cfg.Dnspod.Login_token}, "format": {common.Cfg.Dnspod.Format}, "lang": {common.Cfg.Dnspod.Lang}, "error_on_empty": {common.Cfg.Dnspod.Error_on_empty}, "domain": {domain}}
	resp, err := http.PostForm(urls, values)

	if err == nil {
		res, err = ioutil.ReadAll(resp.Body)

	} else {
		return res, err
	}
	defer resp.Body.Close()

	return
}

//获取域名详细信息
func Select_domains_info(domain string) (res []byte, err error) {
	urls := "https://dnsapi.cn/Domain.Info"
	values := url.Values{"login_token": {common.Cfg.Dnspod.Login_token}, "format": {common.Cfg.Dnspod.Format}, "lang": {common.Cfg.Dnspod.Lang}, "error_on_empty": {common.Cfg.Dnspod.Error_on_empty}, "domain": {domain}}
	resp, err := http.PostForm(urls, values)

	if err == nil {
		res, err = ioutil.ReadAll(resp.Body)

	} else {
		return res, err
	}
	defer resp.Body.Close()

	return
}

//获取站点某记录的id
func Select_record_id(domain string, record string, value string) (id string, err error) {
	var records Dnspod_res
	res0, err := Select_record_info(domain)
	if err == nil {
		if err := json.Unmarshal([]byte(res0), &records); err == nil {
			if records.Status.Code == "1" {
				for _, values := range records.Records {
					if values.Value == value {
						if values.Name == record {
							id = values.Id
							break
						}
					}
				}

			}
		}
	}
	return
}

//改变解析记录是否开启,status:enable|disable
func Set_record_status(domain string, record_id string, status0 string) (status int, err error) {
	var res_status Dnspod_res
	status = 0
	urls := "https://dnsapi.cn/Record.Status"
	values := url.Values{"login_token": {common.Cfg.Dnspod.Login_token}, "format": {common.Cfg.Dnspod.Format}, "lang": {common.Cfg.Dnspod.Lang}, "error_on_empty": {common.Cfg.Dnspod.Error_on_empty}, "domain": {domain}, "record_id": {record_id}, "status": {status0}}
	resp, err := http.PostForm(urls, values)

	if err == nil {
		res, _ := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal([]byte(res), &res_status); err == nil {
			if res_status.Status.Code == "1" {
				status = 1
			}
		}
	} else {
		return status, err
	}
	defer resp.Body.Close()

	return
}

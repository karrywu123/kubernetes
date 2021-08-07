package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"node-monitor/common"
	"strconv"
	"time"

	"node-monitor/libs/logs"

	"github.com/chanyipiaomiao/hltool"
)

type ResDnscom struct {
	code    int            `json:"code"`
	message string         `json:"message"`
	data    ReponseProduct `json:data`
}
type ReponseProduct struct {
	Total    int           `json:"total"`
	Cpage    int           `json:"cpage"`
	LastPage int           `json:"lastPage"`
	List     []ProductList `json:"list"`
}
type ProductList struct {
	Id          int    `json:"id"`
	PackageID   int    `json:"packageID"`
	PackageName string `json:"packageName"`
	Domain      string `json:"domain"`
	DomainID    int    `json:"domainID"`
	EndTime     string `json:"endTime"`
}

// lastStr := "apiKey=*********e5cf1d8f3f966b&domain=dns.com&timestamp=1521005892ecb4ff0e877a8329*******7e9ae673"
// lastStr := "cenusdesign"

func getHash(t string, domain string, page string, perpage string) (h string) {
	lastStr := "apiKey=" + common.Cfg.Dnscomapi.Key + "&page=" + page + "&per_page=" + perpage + "&timestamp=" + t + common.Cfg.Dnscomapi.Secret

	h = hltool.GetMD5(lastStr)

	// fmt.Println(h)

	return

}

//GetProductInfo 查询所有收费解析的情况
func GetProductInfo() (res []byte, err error) {
	page := "1"
	perPage := "1000"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	// timestamp := time.Now().Unix()
	urlstr := "&apiKey=" + common.Cfg.Dnscomapi.Key + "&apiSecret=" + common.Cfg.Dnscomapi.Secret + "&timestamp=" + timestamp + "&hash=" + getHash(timestamp, "123.com", page, perPage)
	urls := "https://www.dns.com/api/product/package/user?page=1&per_page=1000" + urlstr
	// values := url.Values{"apiKey": {common.Cfg.Dnscomapi.Key}, "timestamp": {timestamp}, "hash": {getHash(timestamp, "123.com")}, "page": {page}, "perPage": {perPage}}
	// values := url.Values{"apiKey": {common.Cfg.Dnscomapi.Key}, "timestamp": {timestamp}, "hash": {getHash(timestamp, "123.com")}}
	// values := url.Values{}

	fmt.Println(urls)
	// var values1 map[string]interface{}
	// values1 = make(map[string]interface{})

	// client := &http.Client{}

	// values1["apiKey"] = common.Cfg.Dnscomapi.Key
	// values1["timestamp"] = timestamp
	// values1["hash"] = getHash(strconv.FormatInt(timestamp, 10), "123.com")
	// values1["page"] = page
	// values1["perPage"] = perPage
	// js, _ := json.Marshal(values1)
	// req, _ := http.NewRequest("POST", urls, strings.NewReader(string(js)))
	// fmt.Println(getHash(timestamp, "dns.com"))

	// resp, err := client.Do(req)

	// resp, err := http.PostForm(urls, values)
	resp, err := http.Get(urls)

	if err == nil {
		res, err = ioutil.ReadAll(resp.Body)

	} else {
		logs.Error("Failed to request,err", err)
		return res, err
	}
	defer resp.Body.Close()

	return
}

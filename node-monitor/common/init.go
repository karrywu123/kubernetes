package common

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"node-monitor/libs/logs"

	_ "github.com/go-sql-driver/mysql"
	//"monitor"
	//go get github.com/go-sql-driver/mysql
	//"net"
	//"time"
	"github.com/bitly/go-simplejson"
)

var (
	Cfg           *XmlConfigApp
	Workdir       string
	Logname_today string
	// Nodes_all     map[string]string

	//kangle-control api配置
	api_url     string = "http://kangle-api.clubs999.com"
	conten_type string = "application/json"
	token       string = "***************yayIsInVwZGF0ZVRpbWUiOjE2MjE5MzY0NTF9.LKqXt7L7mBtR8jI7htn2ehlqyYsKBJ_7neizKFAngsI"
)

type XmlConfigApp struct {
	Name             string                  `xml:"name,attr"`
	Log              LogValues               `xml:"log"`
	Dnspod           DnspodPubValues         `xml:"dnspodapi"`
	Dnscomapi        DnscomapiValues         `xml:"dnscomapi"`
	Cloudflareapi    CloudflareValues        `xml:"cloudflareapi"`
	Cloudflaredefend []CloudflaredefendVaues `xml:"cloudflaredefend>filter"`
	Cdnbestapi       CdnbestapiValues        `xml:"cdnbestapi"`
	Cname            string                  `xml:"cname"`
	Domains_adjust   string                  `xml:"domains_adjust"`
	Skypemessage     SkypemessageValues      `xml:"skypemessage"`
	Telegrammessage  TelegrammessageValues   `xml:"telegrammessage"`
	Monitor_dur      string                  `xml:"monitor_dur"`
	// Logpath           string                `xml:"logpath"`
	// Loglevel          string                `xml:"logLevel"`
	// Logname           string                `xml:"logname"`
	Defend_nodes      string              `xml:"defend_nodes"`
	Defend_time       int                 `xml:"defend_time"`
	Cdnbestdb         CdnbestdbValues     `xml:"cdnbestdb"`
	Dnspod_buy_domain string              `xml:"dnspod_buy_domain"`
	Godaddy_apis      []Godaddy_apiValues `xml:"godaddy_apis>authorization"`
	Certs_monitor     string              `xml:"certs_monitor"`
	Rtmpmonitor       []RtmpmonitorValues `xml:"rtmpmonitor>application"`
}

type CloudflaredefendVaues struct {
	Name         xml.Name `xml:"filter"`
	Filterid     string   `xml:"filterid"`
	Cfdomain     string   `xml:"cfdomain"`
	Defenddomain string   `xml:"defenddomain"`
	RandomKey    string   `xml:"randomKey"`
}

//LogValues 日志配置
type LogValues struct {
	LogLevel int `xml:"logLevel"`
	Maxdays  int `xml:"maxdays"`
}

type DnscomapiValues struct {
	Key    string `xml:"key"`
	Secret string `xml:"secret"`
}

type DnspodPubValues struct {
	Login_token    string `xml:"login_token"`
	Format         string `xml:"format"`
	Lang           string `xml:"lang"`
	Error_on_empty string `xml:"error_on_empty"`
}

type CloudflareValues struct {
	Email                string `xml:"email"`
	Auth_api_key         string `xml:"auth_api_key"`
	User_service_api_key string `xml:"user_service_api_key"`
	Content_type         string `xml:"content_type"`
}

type CdnbestapiValues struct {
	Host    string `xml:"host"`
	Uid     string `xml:"uid"`
	Skey    string `xml:"skey"`
	Product string `xml:"product"`
}

type SkypemessageValues struct {
	Url      string `xml:"url"`
	Receiver string `xml:"receiver"`
}

type TelegrammessageValues struct {
	Url       string `xml:"url"`
	Api_token string `xml:"api_token"`
	Msgid     string `xml:"msgid"`
}

type CdnbestdbValues struct {
	Dbtype     string `xml:"dbtype"`
	Datasource string `xml:"datasource"`
}

type Godaddy_apiValues struct {
	// Head_Accept    string `xml:"head_Accept"`
	// Head_Type      string `xml:"head_Type"`
	Name          xml.Name `xml:"authorization"`
	Authorization string   `xml:"auth"`
}

type RtmpmonitorValues struct {
	Name xml.Name `xml:"application"`

	Domain string `xml:"domain"`
	App    string `xml:"app"`
	Stream string `xml:"stream"`
}

func getExecPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	p, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return p, nil
}

func loadxmlconfig(xmlpath string) (*XmlConfigApp, error) {
	content, err := ioutil.ReadFile(xmlpath)
	if err != nil {
		return nil, err
	}
	var result XmlConfigApp
	err = xml.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func Get_nodes_all() (nodes_all map[string]string, err error) {
	//获取cdnbest节点
	nodes_all = make(map[string]string)
	nodeGetRes, err := GetAllNode()
	if err != nil {
		logs.Error("Failed to get all node from kangle api,err,", err)
		return
	} else {
		statuscode := nodeGetRes.Get("statuscode").MustInt()
		if statuscode == 0 {
			nodeList := nodeGetRes.Get("data").Get("data").MustArray()
			for _, node := range nodeList {

				name := node.(map[string]interface{})["name"].(string)
				ipaddr := node.(map[string]interface{})["ipaddr"].(string)
				// subT, _ := node.(map[string]interface{})["sub"].(json.Number).Int64()
				// sub := int(subT)
				statusT, _ := node.(map[string]interface{})["status"].(json.Number).Int64()
				status := int(statusT)
				// comment := node.(map[string]interface{})["comment"].(string)
				if status == 1 { //状态为0时不必监控
					nodes_all[ipaddr] = name
				}
				// logs.Info(name, ipaddr, sub, status, comment)

			}
		} else {
			logs.Error("Failed to get all node from kangle api,response,", nodeGetRes)
		}
	}
	// db, _ := sql.Open(Cfg.Cdnbestdb.Dbtype, Cfg.Cdnbestdb.Datasource)
	// err = db.Ping()
	// defer db.Close()
	// if err != nil {
	// 	logs.Error("Fail to connect to cdnbest db", err)
	// 	//fmt.Println(err)
	// } else {
	// 	rows, err1 := db.Query("SELECT host,mem FROM nodes WHERE 1=1")
	// 	//fmt.Println(err)
	// 	defer rows.Close()
	// 	if err1 != nil {
	// 		logs.Error("Fail to get cdnbest nodes,err", err)
	// 	} else {
	// 		for rows.Next() {
	// 			var node_ip string
	// 			//var node_nickname string
	// 			var node_mem string
	// 			err2 := rows.Scan(&node_ip, &node_mem)
	// 			if err2 == nil {
	// 				nodes_all[node_ip] = node_mem
	// 			}
	// 		}
	// 	}
	// }
	nodes_all["ip"] = "dev-aws001"
	nodes_all["ip"] = "dev-aws002"
	// nodes_all["ip"] = "qiniu006"
	// nodes_all["ip"] = "aliyun003"
	// nodes_all["ip"] = "qiniu007"
	// nodes_all["ip"] = "bwg001"

	return

}

func Init_config() {
	execPath, err := getExecPath()
	if err != nil {
		log.Fatal("Fail to get work directory: %v", err)
	}
	workDir := path.Dir(strings.Replace(execPath, "\\", "/", -1))
	Workdir = workDir
	//fmt.Println(workDir)
	Cfg, err = loadxmlconfig(path.Join(Workdir, "config.xml"))
	if err != nil {
		log.Fatal("Fail to parse 'config.xml': %v", err)
	}

	//获取cdnbest节点
	// Get_nodes_all()

	//判断日志目录和日志文件是否存在
	//Logname_today = Cfg.Logname + "-" + time.Now().Format("2006-01-02")
	// Logname_today = Cfg.Logname
	// _, err = os.Stat(path.Join(Workdir, Cfg.Logpath))
	// if err != nil {
	// 	os.Mkdir(path.Join(Workdir, Cfg.Logpath), 0755)
	// }
	// _, err = os.Stat(path.Join(Workdir,Cfg.Path.Log_marks))
	// if err != nil {
	//     os.Mkdir(path.Join(Workdir,Cfg.Path.Log_marks),0755)
	// }
	// _, err = os.Stat(path.Join(Workdir, Cfg.Logpath+"/"+Logname_today))
	// if err != nil {
	// 	createlogFile, err := os.Create(path.Join(Workdir, Cfg.Logpath+"/"+Logname_today))
	// 	defer createlogFile.Close()
	// 	if err != nil {
	// 		log.Fatal("Fail to create log file", err)
	// 	}
	// }

}

func GetAllNode() (result *simplejson.Json, err error) {

	var values map[string]string

	values = make(map[string]string)

	client := &http.Client{}

	urls := api_url + "/v1/node/comment"

	values["API-TOKEN"] = token
	values["Content-Type"] = conten_type

	req, _ := http.NewRequest("GET", urls, nil)

	for key, value := range values {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)

	if err == nil {
		res, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		var err1 error
		result, err1 = simplejson.NewJson(res)
		if err1 != nil {
			err = err1
		}

	}

	return

}

func SetNodeComment(name string, ipaddr string, sub int, comment string) (result *simplejson.Json, err error) {
	var values map[string]string
	values = make(map[string]string)
	var values1 map[string]interface{}
	values1 = make(map[string]interface{})

	client := &http.Client{}

	urls := api_url + "/v1/node/comment"

	values["API-TOKEN"] = token
	values["Content-Type"] = conten_type

	values1["name"] = name
	values1["ipaddr"] = ipaddr
	values1["sub"] = sub
	values1["comment"] = comment
	js, _ := json.Marshal(values1)

	req, _ := http.NewRequest("POST", urls, strings.NewReader(string(js)))

	for key, value := range values {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)

	if err == nil {
		res, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		var err1 error
		result, err1 = simplejson.NewJson(res)
		if err1 != nil {
			err = err1
		}

	}

	return

}

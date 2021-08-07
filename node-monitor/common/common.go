package common

import (
	// "fmt"
	// "io/ioutil"
	"log"
	"net"
	"time"
	//"net/http"
	//"net/url"
	//"strconv"
	//"strings"
	//"api"
	//"os"
	//"path"
	//"bufio"
)

//判断ip是否存活
func Judge_ip_alive(ip string) (status int, err error) {
	status = 1
	conn, err := net.DialTimeout("tcp", ip+":80", time.Millisecond*500)

	if err != nil {
		status = 0
		// conn = nil
		return status, err
	}
	defer conn.Close()

	return status, err
}

//解析域名
func Resolve_domain(domainname string) (ips []string, err error) {
	ips, err = net.LookupHost(domainname)
	if err != nil {
		log.Fatal("Fail to resolve domain: ", domainname, ",", err)
	}
	return
}

//判断字符串是否符合site_id规则
func Judge_lineid(lineid string) (status int) {
	status = 0
	length := len(lineid)

	if length == 7 || length == 10 {
		if lineid[0:4] == "line" {
			status = 1
		}
	}

	return
}

//扫描所有节点
// func Scan_all_nodes() {
//     //fmt.Println(Nodes_field)
//     for ip,field := range Nodes_field {
//         status,_ := Judge_ip_alive(ip)

//         if status == 1 {
//             fmt.Println(ip + "-" + field[1] + "-" +"is alive")
//             //fmt.Println(123)
//         } else {
//             fmt.Println(ip + "-" + field[1] + "-" + "is not alive")
//         }

//     }
// }

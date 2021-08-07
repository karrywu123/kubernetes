package monitor

import (
	"encoding/json"
	"errors"
	"node-monitor/api"
	"node-monitor/common"
	"time"

	"node-monitor/libs/logs"

	_ "github.com/go-sql-driver/mysql"

	// "strconv"
	"net/http"
	"net/url"
	"strings"
	//"io/ioutil"
)

// var (
// 	Nodes_status       map[string]int
// 	Lines_records_info map[string][]Records_Type
// 	//logFile interface{}
// 	//公共err
// 	// err          error
// 	// Https_status map[string][]Https_status_type
// 	// Log_lock     sync.Mutex
// )

type Records_Type struct {
	Value   string
	Enabled string
	Type    string
	Id      string
}

// type Https_status_type struct {
// 	Domains         string
// 	Cert_comment    string
// 	Cf_https_status string
// }

//从dnspod获取所有解析信息
func getall_records_from_dnspod() (sites_records_info map[string][]Records_Type, err error) {
	sites_records_info = make(map[string][]Records_Type)
	var records_all api.Dnspod_res
	records_all_res, err := api.Select_record_info(common.Cfg.Cname)
	// fmt.Println("response from dnspod", string(records_all_res))
	if err == nil {
		if err := json.Unmarshal([]byte(records_all_res), &records_all); err == nil {
			if records_all.Status.Code == "1" {
				var sites []string
				for _, values := range records_all.Records {
					flag := 0
					for _, value := range sites {
						if values.Name == value {
							flag = 1
						}
					}
					//fmt.Println(values.Name)
					if common.Judge_lineid(values.Name) == 1 {
						if flag == 0 {
							sites = append(sites, values.Name)
						}

						//sites_records_info[values.Name] = [4]string{values.Value,values.Enabled,values.Type,values.Id}
					}
				}
				for _, site := range sites {
					for _, values := range records_all.Records {
						if values.Name == site {
							sites_records_info[site] = append(sites_records_info[site], Records_Type{values.Value, values.Enabled, values.Type, values.Id})
						}
					}
				}

			}
		} else {
			logs.Error("Failed to parse info from dnspod,err,", err, ",respone,", string(records_all_res))
		}
	} else {
		logs.Error("Failed to get info from dnspod,err,", err, ",respone,", string(records_all_res))
	}
	return
}

//判断节点是否存活
func judge_nodes_alive(nodes_all map[string]string) (status map[string]int) {
	status = make(map[string]int)
	var ips []string

	for key, _ := range nodes_all {
		ips = append(ips, key)
	}
	//fmt.Println(ips)
	for _, ip := range ips {
		status[ip] = 1
		//s := make(map[string] []int)
		count := make(map[string]int)
		for i := 0; i < 10; i++ {
			res, _ := common.Judge_ip_alive(ip)
			//s[ip] = append(s[ip],res)
			if res == 1 {
				count[ip] = count[ip] + 1
			}
		}
		if count[ip] < 8 {
			status[ip] = 0
		}
	}

	return
}

// //写日志,log_level: ERROR>WARNING>DEBG>INFO
// func Write_log(error_message string, err error, log_level string) {
// 	flag := 0
// 	switch common.Cfg.Loglevel {
// 	case "ERROR":
// 		if log_level == "ERROR" {
// 			flag = 1
// 		}
// 	case "WARNING":
// 		if log_level == "ERROR" || log_level == "WARNING" {
// 			flag = 1
// 		}
// 	case "DEBG":
// 		if log_level == "ERROR" || log_level == "WARNING" || log_level == "DEBG" {
// 			flag = 1
// 		}
// 	case "INFO":
// 		flag = 1
// 	default:
// 		log.Fatal("Wrong log level")
// 	}
// 	if flag == 1 {
// 		Log_lock.Lock()
// 		logFile, _ := os.OpenFile(path.Join(common.Workdir, common.Cfg.Logpath+"/"+common.Logname_today), os.O_WRONLY|os.O_APPEND, 0644)
// 		//fmt.Println(err1,common.Logname_today)
// 		defer logFile.Close()
// 		if err == nil {
// 			logwrite := log.New(logFile, "["+log_level+"] [", log.LstdFlags)
// 			logwrite.Println(error_message)
// 			//fmt.Println("0")
// 		} else {
// 			logwrite := log.New(logFile, "["+log_level+"] [", log.LstdFlags)
// 			logwrite.Println("\x08]", error_message, "====>", err)
// 			//fmt.Println("1")
// 		}
// 		fmt.Println(error_message)
// 		Log_lock.Unlock()
// 	}

// }

//modify_cdnbest_nodes_comment 调整cdnbest备注
func modify_cdnbest_nodes_comment(lines_records_info map[string][]Records_Type) {
	nodes_all := make(map[string]string)
	nodes_all_sub := make(map[string]int)
	nodes_all_name := make(map[string]string)

	// db, _ := sql.Open(common.Cfg.Cdnbestdb.Dbtype, common.Cfg.Cdnbestdb.Datasource)
	// err := db.Ping()
	// defer db.Close()
	// if err != nil {
	// 	logs.Error("Fail to connect to cdnbest db", err)
	// } else {
	// rows, err := db.Query("SELECT host,mem FROM nodes WHERE 1=1")
	// defer rows.Close()
	// if err != nil {
	// 	logs.Error("Fail to get nodes info from cdnbest db", err)
	// 	//log.Fatal("Fail to get nodes info from cdnbest db",err)
	// } else {
	// 	for rows.Next() {
	// 		var host string
	// 		//var nickname string
	// 		var mem string
	// 		//var ngid int
	// 		err := rows.Scan(&host, &mem)
	// 		if err != nil {
	// 			logs.Error("Fail to get nodes info from cdnbest db", err)
	// 		} else {
	// 			nodes_all[host] = mem
	// 		}
	// 	}
	// }
	//fmt.Println(nodes_all)
	nodeGetRes, err := common.GetAllNode()
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
				subT, _ := node.(map[string]interface{})["sub"].(json.Number).Int64()
				sub := int(subT)
				statusT, _ := node.(map[string]interface{})["status"].(json.Number).Int64()
				status := int(statusT)
				comment := node.(map[string]interface{})["comment"].(string)
				if status == 1 { //状态为0时不必监控
					nodes_all[ipaddr] = comment
					nodes_all_sub[ipaddr] = sub
					nodes_all_name[ipaddr] = name
				}
				// logs.Info(name, ipaddr, sub, status, comment)

			}
		} else {
			logs.Error("Failed to get all node from kangle api,response,", nodeGetRes)
			return
		}
	}

	for line, records := range lines_records_info {
		var record_ip []string
		flag := 0
		for _, record := range records {
			if record.Type == "A" {
				record_ip = append(record_ip, record.Value)
				flag = 1
			}
		}
		if flag == 1 {
			for _, ip := range record_ip {
				flag1 := 0
				for keys, _ := range nodes_all {
					if keys == ip {
						flag1 = 1
					}
				}
				if flag1 == 1 {
					comment_list := strings.Split(nodes_all[ip], "|")
					flag2 := 0
					for _, mem := range comment_list {
						if mem == line {
							flag2 = 1
						}
					}
					if flag2 == 0 {
						comment_list = append(comment_list, line)
						comment := strings.Join(comment_list, "|")
						// stmt, _ := db.Prepare("UPDATE nodes SET mem = ? WHERE host = ?")
						// stmt.Exec(comment, ip)
						setnodecommentres, err := common.SetNodeComment(nodes_all_name[ip], ip, nodes_all_sub[ip], comment)
						logs.Info("Setting comment for node,", setnodecommentres, err)
					}
				}
			}
			for host, mem := range nodes_all {
				comment_list := strings.Split(mem, "|")
				flag3 := 0
				for _, ip := range record_ip {
					if ip == host {
						flag3 = 1
					}
				}
				if flag3 == 1 {
					//fmt.Println("flag3=",flag3)
					continue
				} else {
					//fmt.Println("flag3=",flag3)
					flag4 := 0
					var del_index int
					for index, comm := range comment_list {
						if comm == line {
							flag4 = 1
							del_index = index
						}
					}
					if flag4 == 1 {

						comment_list = append(comment_list[:del_index], comment_list[del_index+1:]...)
						comment := strings.Join(comment_list, "|")
						// stmt, _ := db.Prepare("UPDATE nodes SET mem = ? WHERE host = ?")
						// stmt.Exec(comment, host)
						setnodecommentres, err := common.SetNodeComment(nodes_all_name[host], host, nodes_all_sub[host], comment)
						logs.Info("Setting comment for node,", setnodecommentres, err)
					}
				}
			}
		} else {
			err = errors.New("Wrong node")
			logs.Warn(line+"没有分配节点", err)
		}
	}
	// }
}

//初始化
func init_monitor() (nodes_status map[string]int, lines_records_info map[string][]Records_Type, err error) {
	//获取cdnbest节点
	nodes_all, err := common.Get_nodes_all()
	if err != nil {
		logs.Error("Failed to get all nodes,err,", err)
		return
	}

	//获取所有解析情况
	lines_records_info, err = getall_records_from_dnspod()
	if err != nil {
		logs.Error("Fail to get info from Dnspod,err,", err)
		return

	}

	if err == nil {

		//获取节点是否存活信息
		nodes_status = judge_nodes_alive(nodes_all)

	} else {
		logs.Error("Failed to get all nodes from CDNBest database,err: ", err)
		return
	}

	return

}

//判断线路是否有防御ip,没有自动添加,但是添加后默认为暂停状态
func judge_line_defendnode(line string, lines_records_info map[string][]Records_Type) {
	defend_ips := strings.Split(common.Cfg.Defend_nodes, ",")
	for _, ip := range defend_ips {
		flag := 0
		for _, record := range lines_records_info[line] {
			if record.Value == ip {
				flag = 1
			}
		}
		if flag == 0 {
			_, err := api.Add_record(common.Cfg.Cname, line, "A", ip, "disable")
			if err != nil {
				logs.Error("Failed to add defend nodes for line: ", line)
			}
		}
	}
}

//监控特殊线路，直接切换到CF，万利相关
func monitorToCFLine(wanliCname, wanliLine, normalCname, CFCname string) {
	logs.Info("Monitoring line into CF,", wanliCname, wanliLine, normalCname, CFCname)

	normalId, err1 := api.Select_record_id(wanliCname, wanliLine, normalCname)
	CFlId, err2 := api.Select_record_id(wanliCname, wanliLine, CFCname)

	if err1 == nil && err2 == nil {
		defendcdn01RecordInfo, err := get_domain_records_from_dnspod(wanliCname)
		if err == nil {
			wanliDefendMode := 0
			for _, r := range defendcdn01RecordInfo[wanliLine] {
				if r.Value == CFCname && r.Enabled == "1" {
					wanliDefendMode = 1
				}
			}
			if wanliDefendMode == 0 {
				go func() {
					var message string
					message = "线路可能被攻击: " + "切换线路到CF线路一小时,请注意查看,线路:" + normalCname + "|" + CFCname
					logs.Warn(message)
					api.Set_record_status(wanliCname, normalId, "disable")
					api.Set_record_status(wanliCname, CFlId, "enable")
					status, err3 := Send_Telegram_message(message)
					if status == 0 {
						logs.Error("Failed to send message,err,", err3)
					}

					time.Sleep(60 * time.Minute)

					message = "线路可能被攻击: " + "还原线路到普通线路,请注意查看,线路:" + normalCname
					logs.Warn(message)
					api.Set_record_status(wanliCname, CFlId, "disable")
					api.Set_record_status(wanliCname, normalId, "enable")
					status, err4 := Send_Telegram_message(message)
					if status == 0 {
						logs.Error("Failed to send message,err,", err4)
					}

				}()

			}
		} else {
			logs.Error("Failed to get records info, domain,", wanliCname, ",err,", err)
		}

	} else {
		logs.Error("Failed to get records ID,err1,", err1, ",err2,", err2)
	}

}

func monitor(nodes_status map[string]int, lines_records_info map[string][]Records_Type) {
	for line, records := range lines_records_info {
		// fmt.Println(Nodes_status)
		// fmt.Println(line,records)

		// 判断线路是否有防御节点，如果没有，就自动添加
		judge_line_defendnode(line, lines_records_info)

		/*records:
		  Value string
		  Enabled string
		  Type string
		  Id string */

		//判断线路存活的节点少于2个时，开启防御节点
		defend_ips := strings.Split(common.Cfg.Defend_nodes, ",")
		line_alive := 0
		for _, record := range records {
			flag := 0
			for _, ip := range defend_ips {
				if record.Value == ip {
					flag = 1
				}
			}
			if flag == 0 {
				if nodes_status[record.Value] == 1 {
					line_alive = line_alive + 1
				}
			}
		}

		if line_alive < 2 {
			logs.Warn(line + " is under attacked maybe. Open defending node!!!")
			status, err1 := Send_Telegram_message(line + " is under attacked maybe. Open defending node!!!")
			if status == 0 {
				logs.Error(line+" is under attacked maybe. Open defending node!!! ---发送消息失败", err1)
			}

			//万利线路不切换到自己的高防
			if line == "line_u16ct" || line == "line_se7s4" || line == "line_s4x04" || line == "line_v2pw1" || line == "line_c3ql5" || line == "line_gi4hn" || line == "line_81pay" {
				// fmt.Println(defendFlag)
				//线路单独处理,当防御节点挂掉时,切换到CF线路,并且防御一小时
				//K1
				if line == "line_u16ct" {
					wanliCname := "defendcdn01.com"
					wanliLine := "k1"
					normalCname := "line_u16ct.defendcdn02.com."
					CFCname := "defend.defend002.com."

					monitorToCFLine(wanliCname, wanliLine, normalCname, CFCname)
				}
				//AP
				if line == "line_se7s4" {
					wanliCname := "defendcdn01.com"
					wanliLine := "wanli"
					normalCname := "line_se7s4.defendcdn02.com."
					CFCname := "defend.defend002.com."

					monitorToCFLine(wanliCname, wanliLine, normalCname, CFCname)
				}
				//江湖
				if line == "line_s4x04" {
					wanliCname := "defendcdn01.com"
					wanliLine := "kj"
					normalCname := "line_s4x04.defendcdn02.com."
					CFCname := "defend.defend002.com."

					monitorToCFLine(wanliCname, wanliLine, normalCname, CFCname)
				}
				//jh
				if line == "line_v2pw1" {
					wanliCname := "defendcdn01.com"
					wanliLine := "jh"
					normalCname := "line_v2pw1.defendcdn02.com."
					CFCname := "defend.defend002.com."

					monitorToCFLine(wanliCname, wanliLine, normalCname, CFCname)
				}
				//jh1
				if line == "line_c3ql5" {
					wanliCname := "defendcdn01.com"
					wanliLine := "jh1"
					normalCname := "line_c3ql5.defendcdn02.com."
					CFCname := "defend.defend002.com."

					monitorToCFLine(wanliCname, wanliLine, normalCname, CFCname)
				}
				//king855
				if line == "line_gi4hn" {
					wanliCname := "defendcdn01.com"
					wanliLine := "king855"
					normalCname := "line_gi4hn.defendcdn02.com."
					CFCname := "defend.defend002.com."
					monitorToCFLine(wanliCname, wanliLine, normalCname, CFCname)
				}
				//82pay
				if line == "line_81pay" {
					wanliCname := "defendcdn01.com"
					wanliLine := "81pay"
					normalCname := "line_81pay.defendcdn02.com."
					CFCname := "defend.defend002.com."
					monitorToCFLine(wanliCname, wanliLine, normalCname, CFCname)
				}
			} else {
				for _, ip := range defend_ips {
					// defendFlag = 1
					recorid, err := api.Select_record_id(common.Cfg.Cname, line, ip)
					if err == nil {
						if nodes_status[ip] == 1 {
							api.Set_record_status(common.Cfg.Cname, recorid, "enable")
							api.Add_record_cf(common.Cfg.Cname, "A", line, ip)
						} else {
							// defendFlag = 0
							logs.Warn(line+" defending node: "+ip+" is unreachable!!!", err)
						}
					}
				}
				// for _, ip := range defend_ips {
				// 	for _, record := range records {
				// 		if record.Value == ip {
				// 			if record.Enabled == "1" {
				// 				beego.Error(line+" is not under attacked. Close defending node: "+ip, err)
				// 				status, err1 := Send_Telegram_message(line + " is not under attacked. Close defending node: " + ip)
				// 				if status == 0 {
				// 					beego.Error(line+" is not under attacked. Close defending node: "+ip+" ---发送消息失败", err1)
				// 				}
				// 				recorid, err := api.Select_record_id(common.Cfg.Cname, line, ip)
				// 				if err == nil {
				// 					api.Set_record_status(common.Cfg.Cname, recorid, "disable")
				// 					api.Del_record_cf(common.Cfg.Cname, line, ip)
				// 				}
				// 			}
				// 		}
				// 	}

				// }
			}
		} else {
			for _, ip := range defend_ips {
				for _, record := range records {
					if record.Value == ip {
						if record.Enabled == "1" {
							logs.Error(line + " is not under attacked. Close defending node: " + ip)
							status, err1 := Send_Telegram_message(line + " is not under attacked. Close defending node: " + ip)
							if status == 0 {
								logs.Error(line+" is not under attacked. Close defending node: "+ip+" ---发送消息失败", err1)
							}
							recorid, err := api.Select_record_id(common.Cfg.Cname, line, ip)
							if err == nil {
								api.Set_record_status(common.Cfg.Cname, recorid, "disable")
								api.Del_record_cf(common.Cfg.Cname, line, ip)
							}
						}
					}
				}

			}

			// // var defendFlag int
			// for _, ip := range defend_ips {
			// 	// defendFlag = 1
			// 	recorid, err := api.Select_record_id(common.Cfg.Cname, line, ip)
			// 	if err == nil {
			// 		if Nodes_status[ip] == 1 {
			// 			api.Set_record_status(common.Cfg.Cname, recorid, "enable")
			// 			api.Add_record_cf(common.Cfg.Cname, "A", line, ip)
			// 		} else {
			// 			// defendFlag = 0
			// 			beego.Warn(line+" defending node: "+ip+" is unreachable!!!", err)
			// 		}
			// 	}
			// }
		}

		//判断线路的ip是否正常，不正常关闭解析
		for _, record := range records {
			if nodes_status[record.Value] == 0 {
				// fmt.Println(line + ": " + record.Value + " is unreachable!!!")
				if record.Enabled == "1" {
					flag := 0
					for _, ip := range defend_ips {
						if record.Value == ip {
							flag = 1
						}
					}
					if flag == 0 {
						logs.Warn(line+": "+record.Value+" is unreachable, now close it!!!", "WARNING")
						status, err1 := Send_Telegram_message(line + ": " + record.Value + " is unreachable, now close it!!!")
						if status == 0 {
							logs.Error(line+": "+record.Value+" is unreachable, now close it!!! ---发送消息失败", err1)
						}
						recorid, err := api.Select_record_id(common.Cfg.Cname, line, record.Value)
						if err == nil {
							api.Set_record_status(common.Cfg.Cname, recorid, "disable")
							api.Del_record_cf(common.Cfg.Cname, line, record.Value)
						}
					}
				}
			} else {
				if record.Enabled == "0" {
					flag := 0
					for _, ip := range defend_ips {
						if record.Value == ip {
							flag = 1
						}
					}
					if flag == 0 {

						logs.Warn(line+": "+record.Value+" becomes reachable, now open it!!!", "WARNING")
						status, err1 := Send_Telegram_message(line + ": " + record.Value + " becomes reachable, now open it!!!")
						if status == 0 {
							logs.Error(line+": "+record.Value+" becomes reachable, now open it!!! ---发送消息失败", err1)
						}
						recorid, err := api.Select_record_id(common.Cfg.Cname, line, record.Value)
						if err == nil {
							api.Set_record_status(common.Cfg.Cname, recorid, "enable")
							api.Add_record_cf(common.Cfg.Cname, "A", line, record.Value)
						}
					}
				}
			}
		}
	}
}

// //发送skype报警消息
// func Send_Skype_message(message string) (status int, err error) {
// 	status = 0
// 	var urls string
// 	//urls = common.Cfg.Skypemessage.Url + "?receiver=" + common.Cfg.Skypemessage.Receiver + "&message=" + message
// 	//fmt.Println(urls)
// 	urls = common.Cfg.Skypemessage.Url
// 	values := url.Values{"receiver": {common.Cfg.Skypemessage.Receiver}, "message": {message}}
// 	resp, err := http.PostForm(urls, values)
// 	//defer resp.Body.Close()
// 	if err == nil {
// 		status = 1
// 		//fmt.Println(err)
// 		resp.Body.Close()
// 	}
// 	return
// }

//发送telegram消息
func Send_Telegram_message(message string) (status int, err error) {
	status = 0
	var urls string
	var values map[string]string
	// var values1 map[string]interface{}
	values = make(map[string]string)
	// values1 = make(map[string]interface{})

	client := &http.Client{}
	// urls = common.Cfg.Skypemessage.Url + "?receiver=" + common.Cfg.Skypemessage.Receiver + "&message=" + message
	//fmt.Println(urls)
	urls = common.Cfg.Telegrammessage.Url + "?msgId=" + common.Cfg.Telegrammessage.Msgid + "&msg=" + url.QueryEscape(message)
	values["API-TOKEN"] = common.Cfg.Telegrammessage.Api_token
	// values1 := url.Values{"msgId":{common.Cfg.Telegrammessage.Msgid},"msg":{message}}
	// values1["msgId"]=common.Cfg.Telegrammessage.Msgid
	// values1["msg"]=message
	// js,_ := json.Marshal(values1)
	// fmt.Println(urls)
	// resp, err := http.PostForm(urls,values)
	req, _ := http.NewRequest("POST", urls, nil)
	for key, value := range values {
		req.Header.Add(key, value)
	}
	resp, err := client.Do(req)

	if err == nil {
		status = 1
		//fmt.Println(err)

	} else {
		logs.Error("Failed to send message to telegram,msg: ", message, ",err: ", err)
		return
	}
	defer resp.Body.Close()
	return
}

func get_domain_records_from_dnspod(domain string) (sites_records_info map[string][]Records_Type, err error) {
	sites_records_info = make(map[string][]Records_Type)
	var records_all api.Dnspod_res
	records_all_res, err := api.Select_record_info(domain)
	if err == nil {
		if err := json.Unmarshal([]byte(records_all_res), &records_all); err == nil {
			if records_all.Status.Code == "1" {
				var sites []string
				for _, values := range records_all.Records {

					flag := 0
					for _, value := range sites {
						if values.Name == value {
							flag = 1
						}
					}
					//fmt.Println(values.Name)
					if values.Type != "NS" {
						if flag == 0 {
							sites = append(sites, values.Name)
						}
					}
				}
				for _, site := range sites {
					for _, values := range records_all.Records {
						if values.Name == site {
							if values.Type != "NS" {
								sites_records_info[site] = append(sites_records_info[site], Records_Type{values.Value, values.Enabled, values.Type, values.Id})
							}
						}
					}
				}

				//fmt.Println(sites_records_info)

			}
		} else {
			logs.Error("Failed to parse resultds from dnspod,err,", err)
		}
	} else {
		logs.Error("Failed to get domain records from dnspod,err,", err)
	}
	return
}

//从CF获取所有的解析记录值
// type Records_Type struct {
//     Value string
//     Enabled string
//     Type string
//     Id string
// }
func get_domain_records_from_cf(domain string) (sites_records_info map[string][]Records_Type, err error) {
	sites_records_info = make(map[string][]Records_Type)
	var records_all api.Res_cloudflare_record
	records_all_res, err := api.List_records_id(domain)
	if err == nil {
		if err := json.Unmarshal([]byte(records_all_res), &records_all); err == nil {
			if records_all.Success {
				var sites []string
				for _, values := range records_all.Result {
					flag := 0
					for _, value := range sites {
						if strings.Split(values.Name, "."+domain)[0] == value {
							flag = 1
						}
					}
					//fmt.Println(values.Name)
					if values.Type != "NS" {
						if flag == 0 {
							sites = append(sites, strings.Split(values.Name, "."+domain)[0])
						}
					}
				}
				for _, site := range sites {
					for _, values := range records_all.Result {
						if strings.Split(values.Name, "."+domain)[0] == site {
							if values.Type != "NS" {
								sites_records_info[site] = append(sites_records_info[site], Records_Type{values.Content, "1", values.Type, values.Id})
							}
						}
					}
				}
			}
		} else {
			logs.Error("Failed to parse results from cf,err,", err)
		}
	} else {
		logs.Error("Failed to get domain records from cf,err,", err)
	}
	return
}

//同步DNSPod上面的部分域名的解析到CF上面
// type Records_Type struct {
//     Value string
//     Enabled string
//     Type string
//     Id string
// }
func adjust_DNSPod_to_CFrecord() {
	domains := strings.Split(common.Cfg.Domains_adjust, ",")
	for _, domain := range domains {
		var domain_res_dnspod map[string][]Records_Type
		var domain_res_cf map[string][]Records_Type
		domain_res_cf, err := get_domain_records_from_cf(domain)
		if err != nil {
			logs.Error("Fail to get info from cloudflare,err,", err)
		}
		domain_res_dnspod, err1 := get_domain_records_from_dnspod(domain)
		if err1 != nil {
			logs.Error("Fail to get info from Dnspod,err,", err1)
		}

		// fmt.Println("dnspod",domain_res_dnspod)
		// fmt.Println("cloudflare",domain_res_cf)
		if err == nil {
			if err1 == nil {
				//添加1
				for record, values := range domain_res_dnspod {
					var flag int
					flag = 0
					for record_cf, _ := range domain_res_cf {
						if record == record_cf {
							flag = 1
							break
						}
						if record == "@" {
							if record_cf == domain {
								flag = 1
								break
							}
						}
					}
					//fmt.Println(flag)
					if flag == 0 {
						for _, v := range values {
							if v.Enabled == "1" {
								//添加DNSPod有，而CF没有的记录，DNSPod暂停的记录不添加
								//fmt.Println(v.Value,v.Type)
								_, err2 := api.Add_record_cf(domain, v.Type, record, v.Value)
								//fmt.Println(err2)
								if err2 == nil {
									logs.Info("Add "+domain+"-"+record+"-"+v.Type+"-"+v.Value+" into CF successfully", err2)
								} else {
									logs.Error("Add "+domain+"-"+record+"-"+v.Type+"-"+v.Value+" into CF unsuccessfully", err2)
								}
							}
						}
					}
				}
				//添加2
				for record_cf, values := range domain_res_cf {
					var flag int
					flag = 0
					for record, _ := range domain_res_dnspod {
						if record == record_cf {
							flag = 1
							break
						}
						if record == "@" {
							if record_cf == domain {
								flag = 1
								break
							}
						}
					}
					if flag == 0 {
						//删除DNSPod没有，而CF有得记录
						//fmt.Println(record_cf,values)
						for _, v := range values {
							_, err2 := api.Del_record_cf(domain, record_cf, v.Value)
							//fmt.Println(err2)
							if err2 == nil {
								logs.Info("Del "+domain+"-"+record_cf+"-"+v.Type+"-"+v.Value+" from CF successfully", err2)
							} else {
								logs.Error("Del "+domain+"-"+record_cf+"-"+v.Type+"-"+v.Value+" from CF unsuccessfully", err2)
							}
						}
					}
				}
				//添加3
				for record, values := range domain_res_dnspod {
					//fmt.Println(record,values)
					for _, value := range values {
						if value.Enabled == "0" {
							var flag int
							flag = 0
							var record_adjust string
							var value_adjust string
							var type_adjust string
							for record_cf, values_cf := range domain_res_cf {
								if record == record_cf {
									for _, value_cf := range values_cf {
										var value_cf_Value string
										if value_cf.Type == "CNAME" {
											value_cf_Value = value_cf.Value + "."
										} else {
											value_cf_Value = value_cf.Value
										}
										if value_cf_Value == value.Value {
											flag = 1
											record_adjust = record_cf
											value_adjust = value_cf.Value
											type_adjust = value_cf.Type
											//fmt.Println(record_cf,value_cf)
											break
										}
									}
								}
								if record == "@" {
									if record_cf == domain {
										for _, value_cf := range values_cf {
											var value_cf_Value string
											if value_cf.Type == "CNAME" {
												value_cf_Value = value_cf.Value + "."
											} else {
												value_cf_Value = value_cf.Value
											}
											if value_cf_Value == value.Value {
												flag = 1
												record_adjust = record_cf
												value_adjust = value_cf.Value
												type_adjust = value_cf.Type
												//fmt.Println(record_cf,value_cf)
												break
											}
										}
									}
								}

							}
							if flag == 1 {
								//删除DNSPod已经暂停，而CF有的记录
								//fmt.Println(record_adjust,value)
								_, err2 := api.Del_record_cf(domain, record_adjust, value_adjust)
								//fmt.Println(err2)
								if err2 == nil {
									logs.Info("Del " + domain + "-" + record_adjust + "-" + type_adjust + "-" + value_adjust + " from CF successfully")
								} else {
									logs.Error("Del "+domain+"-"+record_adjust+"-"+type_adjust+"-"+value_adjust+" from CF unsuccessfully", err2)
								}
								break
							}
						}
					}

					//添加4
					for record, values := range domain_res_dnspod {

						for _, value := range values {
							var flag int
							flag = 0
							var record_adjust string
							var value_adjust string
							var type_adjust string
							//var dnspod_enable string
							for record_cf, values_cf := range domain_res_cf {
								if record_cf == record {
									for _, value_cf := range values_cf {
										if value.Enabled == "1" {
											var value_cf_Value string
											if value_cf.Type == "CNAME" {
												value_cf_Value = value_cf.Value + "."
											} else {
												value_cf_Value = value_cf.Value
											}
											//fmt.Println(value.Value==value_cf_Value)
											if value.Value == value_cf_Value {
												flag = 1
												break
											}
										}
									}
								}
								if record == "@" {
									if record_cf == domain {
										for _, value_cf := range values_cf {
											if value.Enabled == "1" {
												var value_cf_Value string
												if value_cf.Type == "CNAME" {
													value_cf_Value = value_cf.Value + "."
												} else {
													value_cf_Value = value_cf.Value
												}
												//fmt.Println(value.Value==value_cf_Value)
												if value.Value == value_cf_Value {
													flag = 1
													break
												}
											}
										}
									}
								}
							}
							if flag == 0 {
								if value.Enabled == "1" {
									//添加DNSPod已经库企企，而CF没有有的记录
									record_adjust = record
									value_adjust = value.Value
									type_adjust = value.Type
									//fmt.Println(record_adjust,value_adjust,type_adjust)
									_, err2 := api.Add_record_cf(domain, type_adjust, record_adjust, value_adjust)
									//fmt.Println(err2)
									if err2 == nil {
										logs.Info("Add " + domain + "-" + record_adjust + "-" + type_adjust + "-" + value_adjust + " into CF successfully")
									} else {
										logs.Error("Add "+domain+"-"+record_adjust+"-"+type_adjust+"-"+value_adjust+" into CF unsuccessfully", err2)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

//清理CF域名缓存，进行同步解析的域名以及defend001.com,defend002.com,defend003.com,defend004.com,defend005.com
func purge_cfdomain_cache() {
	//清理同步解析的域名的缓存
	domains := strings.Split(common.Cfg.Domains_adjust, ",")
	for _, domain := range domains {
		status, err := api.Purge_domain_cache(domain)
		if status == 1 {
			logs.Info("Purge cache of CF successfully,domain: ", domain)
		} else {
			logs.Error("Failed to purge cache of CF,domain: ", domain, ",err:", err)
		}
	}
	//清理防御域名缓存
	domainsDefend := [5]string{"defend001.com", "defend002.com", "defend003.com", "defend004.com", "defend005.com"}
	for _, domainDefend := range domainsDefend {
		status, err := api.Purge_domain_cache(domainDefend)
		if status == 1 {
			logs.Info("Purge cache of CF successfully,domain: ", domainDefend)
		} else {
			logs.Error("Failed to purge cache of CF,domain: ", domainDefend, ",err:", err)
		}
	}

}

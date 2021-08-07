package monitor

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"

	"node-monitor/libs/logs"

	"github.com/cavaliercoder/go-zabbix"
)

func zabbixNginxSockets() {
	// Default approach - without session caching
	url := "http://ip:port/api_jsonrpc.php"
	username := "Admin"
	password := "******"
	nginxIp := "ip"
	// domain := "zsyknk.com"
	// rr := "wsdata"
	// defendCname := "jianghu.defend002.com" //wsdata.zsyknk.com.wswebcdn.com

	session, err := zabbix.NewSession(url, username, password)
	// defer session.Client.Close()
	// defer session.Close()
	// defer session=nil

	if err != nil {
		logs.Error("Failed to Create new session,err,", err)
		return
	}

	version, err := session.GetVersion()

	if err != nil {
		logs.Error("Failed to get version,err,", err)
		return
	}

	logs.Info("Connected to Zabbix API v%s", version)

	//获取hosts
	hosts, err := session.GetHosts(zabbix.HostGetParams{})
	if err != nil {
		logs.Error("Failed to get hosts from zabbix server,err,", err)
		return
	}
	var nginxId string
	for _, host := range hosts {
		if host.Hostname == nginxIp {
			nginxId = host.HostID
		}
	}

	//获取item
	var itemParams zabbix.ItemGetParams
	itemParams.HostIDs = append(itemParams.HostIDs, nginxId)
	itemParams.GetParameters.TextSearch = make(map[string]string)
	itemParams.GetParameters.TextSearch["key_"] = "sockets"
	items, err := session.GetItems(itemParams)
	if err != nil {
		logs.Error("Failed to get item for nginx,err,", err)
		return
	}
	// fmt.Println(items)
	var nginxItemId int
	for _, item := range items {
		if item.ItemName == "All-sockets" {
			nginxItemId = item.ItemID
		}
	}
	// fmt.Println(nginxItemId)

	//获取历史值history,10个
	var historyParams zabbix.HistoryGetParams
	historyParams.ItemIDs = append(historyParams.ItemIDs, strconv.Itoa(nginxItemId))
	historyParams.History = 3
	historyParams.GetParameters = zabbix.GetParameters{}
	historyParams.GetParameters.ResultLimit = 10
	historyParams.GetParameters.SortField = append(historyParams.GetParameters.SortField, "clock")
	historyParams.GetParameters.SortOrder = "DESC"

	historyValue, err := session.GetHistories(historyParams)
	if err != nil {
		logs.Error("Failed to get history value from zabbix for nginx,err,", err)
		return
	}
	// fmt.Println(time.Now().Unix())
	// fmt.Println(historyValue)

	//如果总socket超标，就切换
	// socketNum := historyValue[0].Value
	socketNum, _ := strconv.ParseInt(historyValue[0].Value, 10, 64)
	logs.Info("Wanli nginx IP:", nginxIp, ",all sockets: ", socketNum)
	if socketNum > 8000 {
		// msg := "Wanli jianghu nginx origin IP: " + nginxIp + ",all sockets: " + historyValue[0].Value + ",切换tcp socket 域名cname线路到CloudFlare,域名: " + rr + "." + domain + ",cname: " + defendCname
		msg := "Wanli jianghu nginx origin IP: " + nginxIp + ",all sockets: " + historyValue[0].Value
		logs.Error(msg)
		Send_Telegram_message(msg)
		// changeWangliSocketCname(domain, rr, defendCname)
	}

	session = nil
	return

}

//***
//*****
func changeWangliSocketCname(domain string, rr string, defendCname string) {
	AccessKeyId := "********"
	AccessKeySecret := "*********"

	//获取域名记录信息
	recordResponse, err := getRecordOfDomain(domain)
	if err != nil {
		logs.Error("Failed to query record for domain,", domain, ",err,", err)
		return
	}
	var fanRecord alidns.Record
	for _, record := range recordResponse.DomainRecords.Record {
		if record.RR == rr {
			// fmt.Println(record)
			fanRecord = record
		}
	}
	// client, err := domain.NewClientWithAccessKey("", AccessKeyId, AccessKeySecret)
	// if err != nil {
	// 	beego.Error("Failed to create new clien of ali sdk,err,", err)
	// }
	client, err := alidns.NewClientWithAccessKey("", AccessKeyId, AccessKeySecret)

	// client, err := alidns.NewClientWithAccessKey("cn-hangzhou", "<accessKeyId>", "<accessSecret>")

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = fanRecord.RecordId
	request.RR = rr
	request.Type = fanRecord.Type
	request.Value = defendCname

	response, err := client.UpdateDomainRecord(request)
	if err != nil {
		logs.Error("Failed to update domain record,domain:", domain, ",record:", fanRecord, ",err,", err)
		return
	}
	logs.Info("Changed domain record to cf defend,domain:", domain, ",record value:", request.Value, ",err:", err, ",response is: ", response)
	// fmt.Printf("response is %#v\n", response)

}

func getRecordOfDomain(domain string) (response *alidns.DescribeDomainRecordsResponse, err error) {
	AccessKeyId := "*********"
	AccessKeySecret := "*********"
	client, err := alidns.NewClientWithAccessKey("", AccessKeyId, AccessKeySecret)

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	// request.DomainName = "xalrr.com"
	request.DomainName = domain

	response, err = client.DescribeDomainRecords(request)
	// if err != nil {
	// 	beego.Error("Failed to query record for domain,", domain, ",err,", err)
	// }
	return
	// fmt.Printf("response is %#v\n", response)
}

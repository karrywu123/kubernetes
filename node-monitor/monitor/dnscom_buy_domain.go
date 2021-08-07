package monitor

import (
	"node-monitor/api"

	"node-monitor/libs/logs"
)

func MonitorDnscomBuyDomains() {
	getproductinfofromdnscom()
}

func getproductinfofromdnscom() {
	res_all_res, err := api.GetProductInfo()
	if err != nil {
		logs.Error("Failed to get product info from dns.com,err", err, "ERROR")
		return
	}
	logs.Info(string(res_all_res))

}

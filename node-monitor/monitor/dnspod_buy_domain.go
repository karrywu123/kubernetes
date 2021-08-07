package monitor

import (
	"encoding/json"
	"node-monitor/api"
	"node-monitor/common"
	"strings"
	"time"

	"node-monitor/libs/logs"
)

// func judge_domain_buy(domain string) bool {
// 	domain_string := common.Cfg.Dnspod_buy_domain
// 	domains := strings.Split(domain_string, ",")
// 	flag := false
// 	for _, d := range domains {
// 		if domain == d {
// 			flag = true
// 		}
// 	}
// 	return flag
// }

func monitorVipDomains() {
	domain_string := common.Cfg.Dnspod_buy_domain
	domains := strings.Split(domain_string, ",")

	for _, domain := range domains {

		records_all_res, err := api.Select_domains_info(domain)
		if err != nil {
			logs.Error("Failed to get domain detail from dnspod,err:", err, ",domain: ", domain)
			return
		}
		var records_all api.Res_Domain
		err = json.Unmarshal([]byte(records_all_res), &records_all)
		if err != nil {
			logs.Error("Failed to unmarshal data from dnspod res,err:", err, ",domain: ", domain)
			return
		}
		// fmt.Println(domain, records_all.Domain.Name, records_all.Domain.Is_vip, records_all.Domain.Grade, records_all.Domain.Grade_title, records_all.Domain.Vip_start_at, records_all.Domain.Vip_end_at)
		logs.Info("Monitoring VIP domain of DNSPod,domain: ", domain, ",domain detail: ", records_all.Domain)
		if records_all.Domain.Is_vip != "yes" {
			msg := domain + " DNS is not vip already, please check"
			logs.Warn(msg)
			Send_Telegram_message(msg)
		}

		expire_hours := get_time_duri_dnspod(records_all.Domain.Vip_end_at)
		// fmt.Println(expire_hours, records_all.Domain.Vip_end_at)
		if expire_hours <= 720 && expire_hours > 0 {
			msg := "域名:" + domain + ", DNSPod VIP过期时间:" + records_all.Domain.Vip_end_at + ",DNSPod VIP 收费解析将要过期，请及时处理"
			logs.Warn(msg)
			Send_Telegram_message(msg)
		}
	}

}

func get_time_duri_dnspod(t string) (expirehours int) {
	// tmp := strings.Split(t, "T")
	// day := tmp[0]
	// tmp1 := strings.Split(tmp[1], ".")
	// time0 := tmp1[0]
	// time_expire := day + " " + time0
	t1, _ := time.Parse("2006-01-02 15:04:05", t)
	now := time.Now()
	expirehours = int(t1.Sub(now).Hours())
	return
}

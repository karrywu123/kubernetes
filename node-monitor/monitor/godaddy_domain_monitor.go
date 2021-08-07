package monitor

import (
	"node-monitor/api"
	"node-monitor/common"
	"strings"
	"time"

	"node-monitor/libs/logs"
)

//监控godaddy域名
func monitor_domains_godaddy() {
	var result []api.List_alldomains_res
	var err error
	auths := make([]([]string), len(common.Cfg.Godaddy_apis))

	// get_time_duri("2020-01-18T00:53:47.000Z")

	for _, authri := range common.Cfg.Godaddy_apis {
		a := strings.Split(authri.Authorization, ",")
		auths = append(auths, a)
	}
	// var auths [6][]string = [6][]string{auth1,auth2,auth3,auth4,auth5,auth6}
	// fmt.Println(common.Cfg.Godaddy_api.Authorization1)
	for _, auth := range auths {
		if len(auth) == 0 {
			continue
		}
		logs.Info("Monitoring godaddy domains : " + auth[0])
		result, err = api.List_alldomains(auth[1])
		// fmt.Println(result)
		if err == nil {
			for _, item := range result {
				// fmt.Println(item.Domain,item.Expires,get_time_duri(item.Expires))
				logs.Info("Monitoring godaddy domain: ", item.Domain, ",godaddy account: ", auth[0], ",domain detail: ", ",Expires, ", item.Expires, ",Status, ", item.Status)
				if item.Status == "ACTIVE" {
					if get_time_duri(item.Expires) <= 720 && get_time_duri(item.Expires) > 0 {
						// if get_time_duri(item.Expires) <= 4800 {
						// Write_log(fmt.Sprintf("%d", get_time_duri(item.Expires)),nil,"ERROR")
						logs.Error("账号: " + auth[0] + ", 域名: " + item.Domain + ", 过期时间: " + item.Expires + ", 状态: " + item.Status + ", 域名将过期,请确认并且续费")
						Send_Telegram_message("账号: " + auth[0] + ", 域名: " + item.Domain + ", 过期时间: " + item.Expires + ", 状态: " + item.Status + ", 域名将过期,请确认并且续费")
					}
					if get_time_duri(item.Expires) <= 0 {
						// if get_time_duri(item.Expires) <= 4800 {
						// Write_log(fmt.Sprintf("%d", get_time_duri(item.Expires)),nil,"ERROR")
						logs.Error("账号: " + auth[0] + ", 域名: " + item.Domain + ", 过期时间: " + item.Expires + ", 状态: " + item.Status + ", 域名已过期,请确认并且续费")
						Send_Telegram_message("账号: " + auth[0] + ", 域名: " + item.Domain + ", 过期时间: " + item.Expires + ", 状态: " + item.Status + ", 域名已过期,请确认并且续费")
					}
				}
				// fmt.Println(item.Domain,item.Expires,item.Status)

			}
		} else {
			logs.Error("Failed to get results from godaddy: ", err)
			Send_Telegram_message("Failed to get results from godaddy " + auth[0])
		}

	}

}

//获取域名过期时间和当前时间时间差/hours
// func get_time_duri(t string) (expiredays int){
func get_time_duri(t string) (expirehours int) {
	tmp := strings.Split(t, "T")
	day := tmp[0]
	tmp1 := strings.Split(tmp[1], ".")
	time0 := tmp1[0]
	time_expire := day + " " + time0
	t1, _ := time.Parse("2006-01-02 15:04:05", time_expire)
	now := time.Now()
	expirehours = int(t1.Sub(now).Hours())
	return
}

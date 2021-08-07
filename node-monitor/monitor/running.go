package monitor

import (
	"node-monitor/common"
	"time"

	"node-monitor/libs/logs"
)

//Adjust_record cf调整
func Adjust_record() {
	for {
		adjust_DNSPod_to_CFrecord()
		time.Sleep(120 * time.Second)
	}
}

//Purge_cache 定时刷新CF缓存
func Purge_cache() {
	for {
		purge_cfdomain_cache()
		time.Sleep(300 * time.Second)
	}
}

//Monitor_domains 监控所有godaddy域名过期时间
func Monitor_domains() {
	for {
		monitor_domains_godaddy()
		time.Sleep(86400 * time.Second)
	}
}

//MonitorDnspodBuyDomains 监控DnsPod收费解析域名过期时间
func MonitorDnspodBuyDomains() {
	for {
		monitorVipDomains()
		time.Sleep(86400 * time.Second)
	}
}

//MonitorCerts 监控证书
func MonitorCerts() {
	for {
		monitorcerts()
		time.Sleep(86400 * time.Second)
		// time.Sleep(60 * time.Second)
	}
}

//MonitorRtmp 监控rtmp视频流
func MonitorRtmp() {
	for {
		Monitorrtmp()
		time.Sleep(60 * time.Second)
	}
}

//设置节点备注
func MonitorNodeComment() {
	for {
		//获取所有解析情况
		lines_records_info, err := getall_records_from_dnspod()
		if err != nil {
			logs.Error("Fail to get info from Dnspod,err,", err)
			return

		}
		//调整cdnbest备注
		modify_cdnbest_nodes_comment(lines_records_info)

		time.Sleep(60 * time.Second)
	}
}

//监控
func Monitor_nodes() {

	nodes_status, lines_records_info, err := init_monitor()

	if err == nil {
		monitor(nodes_status, lines_records_info)
		// logs.Debug(nodes_status, lines_records_info)

	} else {
		logs.Error("Failed to init monitor of CDNBest nodes,err: ", err)
	}

}

func RunningMonitorCDNBest() {
	for {
		t := time.Now()
		logs.Info("Start Monitoring CDNBest, start from: " + t.Format("2006-01-02 15:04:05"))

		Monitor_nodes()

		logs.Info("Finish Monitoring CDNBest, time:", time.Since(t), ", start time:", t)

		// var time_duri int
		// time_duri, _ = strconv.Atoi(common.Cfg.Monitor_dur)
		// time.Sleep(time.Duration(time_duri) * time.Minute)
		time.Sleep(10 * time.Second)

	}

}

//CF domain defend域名防御
/*
rule名字 filterId
*.k1668.vip-jianghu-masterapi 7asasa3a2892
*/
func Cf_domain_defend() {
	for {
		for _, cloudflaredefends := range common.Cfg.Cloudflaredefend {
			go adjust_CF_defend_strategy(cloudflaredefends.Cfdomain, cloudflaredefends.Defenddomain, cloudflaredefends.Filterid, cloudflaredefends.RandomKey)
		}
		// go adjust_CF_defend_strategy("defend002.com", "k1668.vip", "asasasdc3a2892")
		time.Sleep(80 * time.Second)
	}
}

//zabbix api 监控 nginx socket数量 ----万利
func ZabbixDefend() {
	for {
		zabbixNginxSockets()

		time.Sleep(10 * time.Second)
	}

}

package main

import (
	"node-monitor/common"
	"node-monitor/monitor"
	"runtime"
	"sync"

	//"api"
	//"strconv"
	//"encoding/json"
	"net/http"
	_ "net/http/pprof"

	"node-monitor/libs/logs"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	//初始化配置
	common.Init_config()

	//初始化日志
	// common.InitLog()
	common.InitLogNew()

	//开启pprof
	go func() {
		ip := "0.0.0.0:6666"
		err := http.ListenAndServe(ip, nil)
		if err != nil {
			logs.Error("Failed to listen address for pprof,err,", err)
		}
	}()

	// spec := fmt.Sprintf("*/%s * * * *", common.Cfg.Monitor_dur)

	// c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(nil)))
	// defer c.Stop()
	// // c.AddFunc("*/5 * * * * ?", func() { running() })
	// // c.AddFunc("*/5 * * * * ?", func() { running() })
	// // c.AddFunc("*/5 * * * * ?", func() { running() })
	// // c.AddFunc("*/5 * * * * ?", func() { running() })
	// // c.AddFunc("*/5 * * * * ?", func() { running() })
	// // c.AddFunc("*/5 * * * * ?", func() { running() })
	// // c.AddFunc("*/5 * * * * ?", func() { running() })
	// // c.AddFunc("*/5 * * * * ?", func() { running() })
	// _, err := c.AddFunc("*/60 * * * * ?", func() { monitor.Monitorrtmp() })
	// c.Start()
	// if err != nil {
	// 	beego.Error("Failed to add func for monitorrtmp,err,", err)
	// } else {
	// 	c.Start()

	// }
	// select {}

	//开启监控
	go monitor.Adjust_record()
	go monitor.Purge_cache()
	go monitor.Monitor_domains()
	go monitor.MonitorDnspodBuyDomains()
	go monitor.MonitorCerts()
	go monitor.MonitorRtmp()
	go monitor.Cf_domain_defend()
	go monitor.ZabbixDefend()
	go monitor.MonitorNodeComment()

	//开启节点监控，并且阻塞主进程
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		monitor.RunningMonitorCDNBest()

		wg.Done()
	}()

	wg.Wait()

}

// func running() {

// 	beego.Info("cron running")
// }

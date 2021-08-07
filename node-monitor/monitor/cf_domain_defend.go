package monitor

import (
	"crypto/md5"
	"fmt"
	"node-monitor/api"
	"strconv"
	"time"

	"node-monitor/libs/logs"
)

/*
加密方法:
域名: *.k1668.vip,暂时不支持https,等测试好了，就可以配置https
加密方法:
Key:    SfPEwqR2p5wl0bpH
timestamp: timestamp.Now / 100 X 100 + 100 X 2    (取200s之后的整数时间戳)
md5(Key + timestamp)
取小写
最终域名 :    md5 加密值 取前10位  + ".k1668.vip"
*/

// var RandomKey string = "SfPEwqR2p5wl0bpH"

func getDomainsFromNow(defenddomain string, rangdomkey string) (domains []string) {
	// var domains []string

	t := time.Now().Unix() / 100 * 100

	var i int64
	for i = 0; i < 36; i++ {

		tStr := strconv.FormatInt(t+100*i, 10)
		PriStr := rangdomkey + tStr
		Md5Str := fmt.Sprintf("%x", md5.Sum([]byte(PriStr)))
		// fmt.Println(PriStr, fmt.Sprintf("%x", Md5Str))
		domain := Md5Str[:10] + "." + defenddomain
		domains = append(domains, domain)
		// fmt.Println(PriStr)
		logs.Debug("Pri:", PriStr, ",Md5:", Md5Str, ",domain:", domain)

	}

	return
}

func adjust_CF_defend_strategy(cfdomain string, defenddomain string, filterId string, rangdomkey string) {
	domains := getDomainsFromNow(defenddomain, rangdomkey)
	// fmt.Println(domains, len(domains))
	var exStr string
	exStr = ""
	for j, d := range domains {
		if j == 0 {
			exStr = exStr + d
		} else {
			exStr = exStr + "\" \"" + d
		}
	}
	// fmt.Println(exStr)

	expression := "(http.host contains \"" + defenddomain + "\" and not http.host in {\"" + exStr + "\"})"
	// fmt.Println(expression)
	logs.Info("Updating firewall rules,", "cfdomain,"+cfdomain+",defenddomain,"+defenddomain+",expression,"+expression)
	status, err := api.SetCfFirewall(cfdomain, defenddomain, filterId, expression)
	if status != 1 || err != nil {
		msg := "Failed to set CloudFlare firewall rules to defend," + ",cfdomain," + cfdomain + ",defenddomain," + defenddomain + ",expression," + expression
		logs.Error(msg, ",err,", err)
		Send_Telegram_message(msg)
	}

}

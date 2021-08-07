package monitor

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"node-monitor/common"
	"strconv"
	"strings"
	"time"

	"node-monitor/libs/logs"
)

func GetDomainsFromCert(ip string, port int) (cert *x509.Certificate, err error) {
	dialer := net.Dialer{Timeout: time.Second * 3}
	conn, err := tls.DialWithDialer(&dialer, "tcp", fmt.Sprintf("%s:%d", ip, port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		logs.Error("Failed to connect: " + err.Error())
		return nil, err
	}
	defer conn.Close()
	state := conn.ConnectionState()
	cert = state.PeerCertificates[0]
	return cert, nil
}

func monitorcerts() {
	domain_string := common.Cfg.Certs_monitor
	domains := strings.Split(domain_string, ",")
	// fmt.Println(domains)
	for _, domain_port := range domains {
		domain := strings.Split(domain_port, ":")[0]
		port, _ := strconv.Atoi(strings.Split(domain_port, ":")[1])
		cert, err := GetDomainsFromCert(domain, port)
		if err != nil {
			logs.Error("Failed to get certs:"+domain_port, ",err: ", err)
			msg := "证书检查: Failed to get certs:" + domain_port + ", 请检查域名的证书端口是否正常"
			Send_Telegram_message(msg)
			// fmt.Println("Failed to get certs:", domain_port)
		} else {
			//fmt.Println(cert.Subject.CommonName, cert.DNSNames, cert.NotBefore, cert.NotAfter)
			logs.Info("Monitoring certs of domain: ", domain_port, ",domain certs detail: ", ",commonname,", cert.Subject.CommonName, ",DNSNames,", cert.DNSNames, ",Expires,", cert.NotAfter, ",cert issuer,", cert.Issuer)
			expire_hours := get_time_duri_certs(cert.NotAfter)
			// fmt.Println(expire_hours)
			// if expire_hours <= 720 && expire_hours > 0 {
			if expire_hours <= 720 {
				// if expire_hours <= 10000 && expire_hours > 0 {
				msg := "检测域名:" + domain_port + ", 域名CN:" + cert.Subject.CommonName + ", 证书支持域名:" + strings.Join(cert.DNSNames, "|") + ", 证书 过期时间:" + cert.NotAfter.String() + ", 证书将要过期，请及时处理"
				logs.Error(msg)
				Send_Telegram_message(msg)
			}
		}
	}
}

func get_time_duri_certs(t time.Time) (expirehours int) {
	// tmp := strings.Split(t, "T")
	// day := tmp[0]
	// tmp1 := strings.Split(tmp[1], ".")
	// time0 := tmp1[0]
	// time_expire := day + " " + time0
	// t1, _ := time.Parse("2006-01-02 15:04:05", t)
	now := time.Now()
	expirehours = int(t.Sub(now).Hours())
	return
}

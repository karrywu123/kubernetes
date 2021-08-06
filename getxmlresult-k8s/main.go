package main

import (
	"fmt"
	"getxmlresult-k8s/common"
	"html/template"
	"log"
	"os"
	"strings"
)

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("    ", "xml PATH:", "./project_config.xml")
	fmt.Println("    ", os.Args[0], "xml", "project")
	fmt.Println("    ", "or")
	fmt.Println("    ", os.Args[0], "yml", "project")
	fmt.Println("    ", "or")
	fmt.Println("    ", os.Args[0], "nginx", "project")

}

type nginxValues struct {
	ProjectName string
	Domains     []string
	IpPorts     []string
}

type ymlValues struct {
	Namespace   string
	ProjectName string
	Javaname    string
	Ports       []string
	Ip          string
}

func main() {
	if len(os.Args) != 3 {
		printHelp()
		log.Fatal("Wrong paratermeters,111")
	}
	if os.Args[1] != "xml" && os.Args[1] != "yml" && os.Args[1] != "nginx" {
		printHelp()
		log.Fatal("Wrong paratermeters,222")
	}
	//解析xml
	common.InitSetting()
	// fmt.Println(common.Cfg)
	//获取配置
	if os.Args[1] == "xml" {
		projectname := os.Args[2]
		for _, p := range common.Cfg.Projects {
			if p.ProjectName == projectname {
				outString := p.ProjectName + "|" + p.Javaname + "|" + p.Type + "|" + p.Ports + "|" + p.Domains
				fmt.Println(outString)
			}
		}
	}

	//生成yml并且打印
	if os.Args[1] == "yml" {
		projectname := os.Args[2]
		namespace := common.Cfg.Namespace
		for _, p := range common.Cfg.Projects {
			if p.ProjectName == projectname {

				tmpl, _ := template.New("yml").Parse(common.GetYmlTemplate(p.Type))
				// outString := p.ProjectName + "|" + p.Javaname + "|" + p.Type + "|" + p.Ips + "|" + p.Ports

				// Ips := strings.Split(p.Ips, ",")
				// for i, ip := range Ips {
				pValue := ymlValues{}
				pValue.Namespace = namespace
				pValue.ProjectName = p.ProjectName
				pValue.Javaname = p.Javaname
				pValue.Ports = strings.Split(p.Ports, ",")
				// pValue.Ip = ip
				// if i == 0 {
				fmt.Println("Project name: ", pValue.ProjectName)
				// } else {
				// fmt.Println("Project name: ", pValue.ProjectName+strconv.Itoa(i), ", ip:", ip)
				// }
				tmpl.Execute(os.Stdout, pValue)
				fmt.Println("")
				// }

			}
		}
	}
	//生成nginx配置文件并且打印
	if os.Args[1] == "nginx" {
		projectname := os.Args[2]
		for _, p := range common.Cfg.Projects {
			if p.ProjectName == projectname {
				nValue := nginxValues{}
				nValue.ProjectName = p.ProjectName
				nValue.Domains = strings.Split(p.Domains, ",")

				//获取第一个端口,作为默认端口
				// port_tmp := strings.Split(p.Ports, ",")[0]
				// port := strings.Split(port_tmp, ":")[0]

				// Ips := strings.Split(p.Ips, ",")
				// for _, ip := range Ips {
				nValue.IpPorts = append(nValue.IpPorts, "virtualiIp:"+"k8sSvcPort")
				// }
				tmpl, _ := template.New("nginx").Parse(common.GetNginxTemplate(p.Type))
				tmpl.Execute(os.Stdout, nValue)
				fmt.Println("")
			}
		}
	}
}

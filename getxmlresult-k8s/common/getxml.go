package common

import (
	"encoding/xml"
	"io/ioutil"
	"log"
)

var (
	Cfg     *XMLConfigApp
	Workdir string
	// Logname_today string
	// Nodes_all map[string] string
)

//XMLConfigApp 配置文件结构体
type XMLConfigApp struct {
	Namespace string          `xml:"namespace"`
	Projects  []ProjectsValue `xml:"projects>project"`
}

type ProjectsValue struct {
	Name        xml.Name `xml:"project"`
	ProjectName string   `xml:"name"`
	Javaname    string   `xml:"javaname"`
	Type        string   `xml:"type"`
	Ports       string   `xml:"ports"`
	Domains     string   `xml:"domains"`
}

func loadxmlconfig(xmlpath string) (*XMLConfigApp, error) {
	content, err := ioutil.ReadFile(xmlpath)
	if err != nil {
		return nil, err
	}
	var result XMLConfigApp
	err = xml.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//InitSetting 初始化配置文件
func InitSetting() {
	var err error
	Cfg, err = loadxmlconfig("project_config.xml")
	// Cfg, err = loadxmlconfig("E:\\work\\855cash项目\\系统配置\\project_config.xml")
	// Cfg, err = loadxmlconfig("D:\\go_projects\\src\\getxmlresult-k8s\\project_config.xml")
	if err != nil {
		log.Fatal("Fail to parse 'project_config.xml': ", err)
	}
}

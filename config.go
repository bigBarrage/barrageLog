package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Unknwon/goconfig"
)

var (
	broadcastingStationHost = "localhost"
	broadcastingStationPort = "9981"
	broadcastingStationUri  = "broadcasting"

	mongoAddrList   = []Addrs{Addrs{Host: "localhost", Port: "27017"}}
	mongoUsername   = "root"
	mongoPassword   = ""
	mongoDatabase   = "barrageLog"
	mongoOptions    = make(map[string]string)
	mongoCharset    = "utf8"
	mongoReplicaSet = ""
)

type Addrs struct {
	Host string
	Port string
}

func init() {
	configFile := flag.String("config-file", "./config.ini", "path of config file")
	flag.Parse()

	c, err := goconfig.LoadConfigFile(*configFile)
	if err != nil {
		fmt.Println("读取配置文件失败：", err)
		os.Exit(1)
	}

	//加载广播站配置
	bc, err := c.GetSection("broadcasting")
	if err == nil {
		if b, ok := bc["host"]; ok && b != "" {
			broadcastingStationHost = b
		}

		if b, ok := bc["port"]; ok && b != "" {
			broadcastingStationPort = b
		}

		if b, ok := bc["uri"]; ok && b != "" {
			broadcastingStationUri = b
		}
	}

	//加载mongo配置
	mg, err := c.GetSection("mongo")
	if err == nil {
		if addrStr, ok := mg["addrs"]; ok {
			mongoAddrList = make([]Addrs, 0, 3)
			addSlice := strings.Split(addrStr, " ")
			for _, v := range addSlice {
				a := strings.Split(v, ":")
				if len(a) != 2 {
					continue
				}
				add := Addrs{}
				add.Host = a[0]
				add.Port = a[1]
				mongoAddrList = append(mongoAddrList, add)
			}
		}
		if u, ok := mg["username"]; ok {
			mongoUsername = u
		}
		if pwd, ok := mg["password"]; ok {
			mongoPassword = pwd
		}
		if d, ok := mg["database"]; ok {
			mongoDatabase = d
		}
		if r, ok := mg["replicaset"]; ok {
			mongoReplicaSet = r
		}
		if c, ok := mg["charset"]; ok {
			mongoCharset = c
		}
	}
}

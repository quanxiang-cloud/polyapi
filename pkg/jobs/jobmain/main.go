package main

import (
	"flag"
	"fmt"
	"os"

	// v0.7.3_apidoc
	_ "github.com/quanxiang-cloud/polyapi/pkg/jobs/jobimpl/namespace"
	_ "github.com/quanxiang-cloud/polyapi/pkg/jobs/jobimpl/polyapi"
	_ "github.com/quanxiang-cloud/polyapi/pkg/jobs/jobimpl/rawapi"
	_ "github.com/quanxiang-cloud/polyapi/pkg/jobs/jobimpl/rawpoly"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/version"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/jobs/jobcenter"

	"github.com/quanxiang-cloud/cabin/logger"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
)

const toolVersion = "v0.7.3_apidoc"

var configPath = flag.String("config", "../../../configs/config.yml", "-config 配置文件地址")

func main() {
	wd, err := os.Getwd()
	fmt.Println(wd, version.FullVersion())
	fmt.Println("polyapi data clean, version =", toolVersion)
	flag.Parse()

	jobcenter.ShowList()
	//return

	conf, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("db:", conf.Mysql.Host, conf.Mysql.DB)

	db, err := mysql2.New(conf.Mysql, logger.Logger)
	if err != nil {
		panic(err)
	}

	if err := jobcenter.Run(db); err != nil {
		panic(err)
	}
}

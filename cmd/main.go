package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/quanxiang-cloud/polyapi/api/router"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/version"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/validateinit"
	"github.com/quanxiang-cloud/polyapi/pkg/config"

	"github.com/quanxiang-cloud/cabin/logger"
)

var (
	configPath = flag.String("config", "../configs/config.yml", "-config config-file-path")
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host polyapi
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}

func main() {
	flag.Parse()

	conf, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}

	logger.Logger = logger.New(&conf.Log)

	logger.Logger.Infof("=====polyapi Version=%s starting", version.FullVersion())

	// start router
	router, err := router.NewRouter(conf)
	if err != nil {
		panic(err)
	}

	//NOTE: validate waiting initialize variables
	if err := validateinit.ValidateInit(); err != nil {
		panic(err)
	}

	go router.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			//	router.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

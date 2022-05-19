package router

import (
	"fmt"

	restful "github.com/quanxiang-cloud/polyapi/api/restful"
	"github.com/quanxiang-cloud/polyapi/api/router/register"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/polyhost"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"
	"github.com/quanxiang-cloud/polyapi/pkg/probe"

	"github.com/gin-gonic/gin"
	ginLogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

// Router Router
type Router struct {
	c           *config.Config
	engine      *gin.Engine
	engineInner *gin.Engine // engine for inner service only

	*probe.Probe
}

// NewRouter create a router
func NewRouter(c *config.Config) (*Router, error) {
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}

	v1 := engine.Group("/api/v1/polyapi")

	v1.Any("/ping", restful.PingPong)

	raw, err := restful.NewRawAPI(c)
	if err != nil {
		return nil, err
	}

	// /api/v1/polyapi/raw/doc/*apiPath
	apiPath := fmt.Sprintf("*%s", restful.PathargAPIPath)

	rawGroup := v1.Group("/raw")
	{
		rawGroup.POST(fmt.Sprintf("/upload/%s", apiPath), raw.RegFile)
		rawGroup.POST(fmt.Sprintf("/reg/%s", apiPath), raw.RegSwagger)
		rawGroup.POST(fmt.Sprintf("/delete/%s", apiPath), raw.Del)
		rawGroup.POST(fmt.Sprintf("/query/%s", apiPath), raw.Query)
		//rawGroup.Any(fmt.Sprintf("/request/%s", apiPath), raw.Request)
		//rawGroup.POST(fmt.Sprintf("/doc/%s", apiPath), raw.QueryAPIDoc)
		//rawGroup.GET(fmt.Sprintf("/doc/%s", apiPath), raw.QueryAPIDoc)
		rawGroup.POST(fmt.Sprintf("/list/%s", apiPath), raw.List)                   //
		rawGroup.POST(fmt.Sprintf("/active/%s", apiPath), raw.Active)               //
		rawGroup.POST(fmt.Sprintf("/listInService/%s", apiPath), raw.ListInService) //
		rawGroup.POST(fmt.Sprintf("/search/%s", apiPath), raw.Search)               // TODO
	}

	v1.Any(fmt.Sprintf("/request/%s", apiPath), raw.APIProviderRequest)
	v1.POST(fmt.Sprintf("/doc/%s", apiPath), raw.APIProviderQueryDoc)

	apiOper, err := restful.NewAPIOperator(c)
	if err != nil {
		return nil, err
	}
	v1.POST("/swagger", apiOper.QuerySwagger)

	poly, err := restful.NewPolyAPI(c)
	if err != nil {
		return nil, err
	}

	polyGroup := v1.Group("poly")
	{
		polyGroup.POST(fmt.Sprintf("/create/%s", apiPath), poly.Create)
		polyGroup.POST(fmt.Sprintf("/delete/%s", apiPath), poly.Delete)
		polyGroup.POST(fmt.Sprintf("/save/%s", apiPath), poly.UpdateArrange)
		polyGroup.POST(fmt.Sprintf("/query/%s", apiPath), poly.GetArrange)
		polyGroup.POST(fmt.Sprintf("/build/%s", apiPath), poly.Build)
		//polyGroup.Any(fmt.Sprintf("/request/%s", apiPath), poly.Request)
		//polyGroup.POST(fmt.Sprintf("/doc/%s", apiPath), poly.Request)
		//polyGroup.POST(fmt.Sprintf("/doc/%s", apiPath), poly.QueryAPIDoc) //
		polyGroup.POST(fmt.Sprintf("/active/%s", apiPath), poly.Active) //
		polyGroup.POST(fmt.Sprintf("/list/%s", apiPath), poly.List)     //
		polyGroup.POST(fmt.Sprintf("/search/%s", apiPath), poly.Search)

		polyGroup.POST(fmt.Sprintf("/enums/:%s", restful.PathargType), poly.ShowEnum)
	}

	ns, err := restful.NewAPINamespace(c)
	if err != nil {
		return nil, err
	}
	nsGroup := v1.Group("namespace")
	{
		nsGroup.POST(fmt.Sprintf("/create/%s", apiPath), ns.Create) //
		nsGroup.POST(fmt.Sprintf("/delete/%s", apiPath), ns.Delete) //
		nsGroup.POST(fmt.Sprintf("/update/%s", apiPath), ns.Update) //
		nsGroup.POST(fmt.Sprintf("/active/%s", apiPath), ns.Active) //
		nsGroup.POST(fmt.Sprintf("/list/%s", apiPath), ns.List)     //
		nsGroup.POST(fmt.Sprintf("/query/%s", apiPath), ns.Query)   //
		nsGroup.POST("/appPath", ns.APPPath)                        //
		nsGroup.POST("/initAppPath", ns.InitAPPPath)                //
		nsGroup.POST(fmt.Sprintf("/search/%s", apiPath), ns.Search)
		nsGroup.POST(fmt.Sprintf("/tree/%s", apiPath), ns.Tree)
	}

	svs, err := restful.NewAPIService(c)
	if err != nil {
		return nil, err
	}
	svsGroup := v1.Group("service")
	{
		svsGroup.POST(fmt.Sprintf("/create/%s", apiPath), svs.Create) //
		svsGroup.POST(fmt.Sprintf("/delete/%s", apiPath), svs.Delete) //
		svsGroup.POST(fmt.Sprintf("/update/%s", apiPath), svs.Update) //
		svsGroup.POST(fmt.Sprintf("/active/%s", apiPath), svs.Active) //
		svsGroup.POST(fmt.Sprintf("/list/%s", apiPath), svs.List)     //
		svsGroup.POST(fmt.Sprintf("/query/%s", apiPath), svs.Query)   //
		svsGroup.POST(fmt.Sprintf("/updateProperty/%s", apiPath), svs.UpdateProperty)
	}

	apikey, err := restful.NewAPIKey(c)
	if err != nil {
		return nil, err
	}
	apikeyGroup := v1.Group("apikey")
	{
		apikeyGroup.POST("/create", apikey.Create) //
		apikeyGroup.POST("/delete", apikey.Delete) //
		apikeyGroup.POST("/update", apikey.Update) //
		apikeyGroup.POST("/active", apikey.Active) //
		apikeyGroup.POST("/list", apikey.List)     //
		apikeyGroup.POST("/query", apikey.Query)   //
	}

	holdkey, err := restful.NewAPIKeyHolding(c)
	if err != nil {
		return nil, err
	}
	holdkeydGroup := v1.Group("holdingkey")
	{
		holdkeydGroup.POST("/upload", holdkey.Upload)       //
		holdkeydGroup.POST("/delete", holdkey.Delete)       //
		holdkeydGroup.POST("/update", holdkey.Update)       //
		holdkeydGroup.POST("/active", holdkey.Active)       //
		holdkeydGroup.POST("/list", holdkey.List)           //
		holdkeydGroup.POST("/query", holdkey.Query)         //
		holdkeydGroup.POST("/authTypes", holdkey.AuthTypes) //
		holdkeydGroup.POST("/sample", holdkey.Sample)

	}

	apiSchema, err := restful.NewAPISchema(c)
	if err != nil {
		return nil, err
	}
	schemaGroup := v1.Group("schema")
	{
		schemaGroup.POST(fmt.Sprintf("/create/%s", apiPath), apiSchema.GenSchema)
		schemaGroup.POST(fmt.Sprintf("/list/%s", apiPath), apiSchema.ListSchema)
		schemaGroup.POST(fmt.Sprintf("/query/%s", apiPath), apiSchema.QuerySchema)
		schemaGroup.POST(fmt.Sprintf("/delete/%s", apiPath), apiSchema.DeleteSchema)
		schemaGroup.POST(fmt.Sprintf("/request/%s", apiPath), apiSchema.APISchemaRequest)
	}

	/*
		permit, err := restful.NewAPIPermit(c)
		if err != nil {
			return nil, err
		}
		permitGroup := v1.Group("permit/group")
		{
			permitGroup.POST(fmt.Sprintf("/create/%s", apiPath), permit.CreateGroup) //
			permitGroup.POST(fmt.Sprintf("/delete/%s", apiPath), permit.DeleteGroup) //
			permitGroup.POST(fmt.Sprintf("/update/%s", apiPath), permit.UpdateGroup) //
			permitGroup.POST(fmt.Sprintf("/active/%s", apiPath), permit.ActiveGroup) //
			permitGroup.POST(fmt.Sprintf("/list/%s", apiPath), permit.ListGroup)     //
			permitGroup.POST(fmt.Sprintf("/query/%s", apiPath), permit.QueryGroup)   //
		}

		permitElemGroup := v1.Group("permit/elem")
		{
			permitElemGroup.POST(fmt.Sprintf("/add/%s", apiPath), permit.AddElem)   //
			permitElemGroup.POST("/delete", permit.DeleteElem)                      //
			permitElemGroup.POST("/update", permit.UpdateElem)                      //
			permitElemGroup.POST("/active", permit.ActiveElem)                      //
			permitElemGroup.POST(fmt.Sprintf("/list/%s", apiPath), permit.ListElem) //
			permitElemGroup.POST("/query", permit.QueryElem)                        //
		}

		permitGrantGroup := v1.Group("permit/grant")
		{
			permitGrantGroup.POST(fmt.Sprintf("/grant/%s", apiPath), permit.AddGrant) //
			permitGrantGroup.POST("/delete", permit.DeleteGrant)                      //
			permitGrantGroup.POST("/update", permit.UpdateGrant)                      //
			permitGrantGroup.POST("/active", permit.ActiveGrant)                      //
			permitGrantGroup.POST(fmt.Sprintf("/list/%s", apiPath), permit.ListGrant) //
			permitGrantGroup.POST("/query", permit.QueryGrant)                        //
		}
	*/

	engineInner, err := newRouter(c)
	if err != nil {
		return nil, err
	}

	opt, err := restful.NewOperator(c)
	if err != nil {
		return nil, err
	}

	innerGroup := engineInner.Group("/api/v1/polyapi/inner")
	{

		innerGroup.Any(fmt.Sprintf("/request/%s", apiPath), raw.InnerAPIProviderRequest)       // /request
		innerGroup.POST(fmt.Sprintf("/regSwagger/%s", apiPath), raw.InnerRegSwagger)           // /regSwagger
		innerGroup.POST(fmt.Sprintf("/regSwaggerAlone/%s", apiPath), raw.InnerRegSwaggerAlone) // /regSwaggerAlone
		innerGroup.POST(fmt.Sprintf("/createNamespace/%s", apiPath), ns.InnerCreate)           // /createNamespace
		innerGroup.POST(fmt.Sprintf("/deleteNamespace/%s", apiPath), ns.InnerDelete)           // /deleteNamespace
		innerGroup.POST("/initAppPath", ns.InnerInitAPPPath)                                   // /initAppPath
		innerGroup.POST(fmt.Sprintf("/createService/%s", apiPath), svs.InnerCreate)            // /createService
		innerGroup.POST(fmt.Sprintf("/delApp/%s", apiPath), opt.InnerDelApp)                   // /delApp/:appID
		innerGroup.POST(fmt.Sprintf("/validApp/%s", apiPath), opt.InnerUpdateAppValid)         // /validApp/:appID
		innerGroup.POST(fmt.Sprintf("/exportApp/%s", apiPath), opt.InnerExportApp)             // /exportApp/:appID

		innerGroup.POST("/import", opt.InnerImport)
	}

	register.InitialRoutes(engine, innerGroup) // regist pre defined routes

	// Only init RawPoly for the inner interface
	_, err = restful.CreateRawPoly(c)
	if err != nil {
		return nil, err
	}

	r := &Router{
		c:           c,
		engine:      engine,
		engineInner: engineInner,

		Probe: probe.New(),
	}

	r.probe()
	return r, nil
}

func (r *Router) probe() {
	r.engine.GET("liveness", func(c *gin.Context) {
		r.Probe.LivenessProbe(c.Writer, c.Request)
	})

	r.engine.Any("readiness", func(c *gin.Context) {
		r.Probe.ReadinessProbe(c.Writer, c.Request)
	})
}

func newRouter(c *config.Config) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}

	if err := polyhost.SetSchemaHost(c.MyHostBase); err != nil {
		return nil, err
	}

	gin.SetMode(c.Model)
	engine := gin.New()

	engine.Use(ginLogger.LoggerFunc(),
		ginLogger.RecoveryFunc(), gate.Cors())

	return engine, nil
}

// Run router
func (r *Router) Run() {
	r.Probe.SetRunning()
	go r.engineInner.Run(r.c.PortInner)
	r.engine.Run(r.c.Port)
}

// Close router
func (r *Router) Close() {
}

func (r *Router) router() {

}

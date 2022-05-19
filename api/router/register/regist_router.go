package register

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var regist registGroups

// InitialRoutes initialize pre registered routes
func InitialRoutes(e *gin.Engine, innerGroup *gin.RouterGroup) error {
	return regist.InitialRoutes(e, innerGroup)
}

// RegInnerRoute regist an inner route
func RegInnerRoute(r Route) {
	regist.Inner = append(regist.Inner, r)
}

// RegRouterGroup regist a router group
func RegRouterGroup(r RouteGroup) {
	regist.Groups = append(regist.Groups, r)
}

// InitialRoutes initialize pre registered routes
func (rg *registGroups) InitialRoutes(e *gin.Engine, innerGroup *gin.RouterGroup) error {
	if innerGroup != nil {
		initRoutes(innerGroup, rg.Inner)
	}
	for _, v := range rg.Groups {
		g := e.Group(v.Path)
		initRoutes(g, v.Routes)
	}
	return nil
}

func initRoutes(g *gin.RouterGroup, rs []Route) {
	for _, v := range rs {
		initRoute(g, v)
	}
}

func initRoute(g *gin.RouterGroup, r Route) {
	switch r.Method {
	case "any", "ANY", "*":
		g.Any(r.Path, r.Func)
	case http.MethodGet, http.MethodHead, http.MethodPost,
		http.MethodPut, http.MethodPatch, http.MethodDelete,
		http.MethodOptions, http.MethodConnect, http.MethodTrace:
		g.Handle(r.Method, r.Path, r.Func)
	default:
		panic(fmt.Errorf("unknown http method '%s'", r.Method))
	}
}

// Route is a route pre define
type Route struct {
	Path   string
	Method string
	Func   gin.HandlerFunc
}

// RouteGroup is a route group pre define
type RouteGroup struct {
	Path   string
	Routes []Route
}

type registGroups struct {
	Inner  []Route
	Groups []RouteGroup
}

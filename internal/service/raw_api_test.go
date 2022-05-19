package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RawAPISuite struct {
	suite.Suite
	ctx     context.Context
	conf    *config.Config
	r       *gin.Engine
	rawAPI  RawAPI
	polyAPI PolyAPI
}

func _TestRawAPI(t *testing.T) {
	suite.Run(t, new(RawAPISuite))
}

func (suite *RawAPISuite) SetupSuite() {
	var err error
	suite.ctx = context.TODO()
	suite.conf, err = config.NewConfig("./testdata/local/polyapi.yaml")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.conf)
	suite.conf.Log.Level = 1 // warn
	logger.Logger = logger.New(&suite.conf.Log)
	assert.Nil(suite.T(), err)
	suite.rawAPI, err = CreateRaw(suite.conf)
	suite.polyAPI, err = CreatePoly(suite.conf)

}
func (suite *RawAPISuite) TestCreateRaw() {

	file, _ := os.Open("./testdata/swaggerApi.json")
	all, _ := io.ReadAll(file)
	req := &RegReq{
		Swagger: string(all),
	}
	resp, err := suite.rawAPI.RegSwagger(suite.ctx, req)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), resp)
}

func (suite *RawAPISuite) TestGetRaw() {
	req := &QueryReq{
		APIPath: "",
	}
	if req.APIPath == "" {
		return
	}
	resp, err := suite.rawAPI.Query(suite.ctx, req)
	fmt.Printf("get raw %+v %v\n", resp, err)
	fmt.Printf("%+v\n", resp.Content)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), resp)
}

func (suite *RawAPISuite) TestQuerySwagger() {
	req := &QueryRawSwaggerReq{
		APIPath: []string{"/system/app/bwp2w/raw/customer/dr/rdz90ppf.r"},
	}
	suite.rawAPI.QuerySwagger(suite.ctx, req)
}

func (suite *RawAPISuite) TestQueryInBatches() {
	req := QueryInBatchesReq{
		APIPathList: []string{"/system/app/2zfvz/raw/customer/one_group/api.r", "/system/app/bwp2w/raw/inner/form/custom/A_create.r"},
	}
	_, err := adaptor.GetRawAPIOper().QueryInBatches(suite.ctx, &req)
	suite.Nil(suite.T(), err)
}

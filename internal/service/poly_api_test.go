package service

// import (
// 	"fmt"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

// func (s *RawAPISuite) _TestPolyApis() {
// 	newID := ""
// 	{
// 		req := &PolyCreateReq{}
// 		resp, err := s.polyAPI.Create(s.ctx, req)
// 		fmt.Printf("create %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 		newID = resp.ID
// 	}
// 	{
// 		req := &PolyGetArrangeReq{ID: newID}
// 		resp, err := s.polyAPI.GetArrange(s.ctx, req)
// 		fmt.Printf("get %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	//--------------------------------------------------------------------------
// 	{
// 		time.Sleep(time.Second)
// 		req := &PolyUpdateArrangeReq{
// 			ID:      newID,
// 			Arrange: "arrange V1",
// 		}
// 		resp, err := s.polyAPI.UpdateArrange(s.ctx, req)
// 		fmt.Printf("update V1 %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	{
// 		req := &PolyGetArrangeReq{ID: newID}
// 		resp, err := s.polyAPI.GetArrange(s.ctx, req)
// 		fmt.Printf("get %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	//--------------------------------------------------------------------------
// 	{
// 		time.Sleep(time.Second)
// 		req := &PolyUpdateArrangeReq{
// 			ID:      newID,
// 			Arrange: "arrange V2",
// 		}
// 		resp, err := s.polyAPI.UpdateArrange(s.ctx, req)
// 		fmt.Printf("update V2 %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	{
// 		req := &PolyGetArrangeReq{ID: newID}
// 		resp, err := s.polyAPI.GetArrange(s.ctx, req)
// 		fmt.Printf("get %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	//--------------------------------------------------------------------------
// 	{
// 		time.Sleep(time.Second)
// 		req := &PolyUpdateScriptReq{
// 			ID:     newID,
// 			Script: `"script V1"`,
// 		}
// 		resp, err := s.polyAPI.UpdateScript(s.ctx, req)
// 		fmt.Printf("script V1 %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	{
// 		req := &PolyGetScriptReq{ID: newID}
// 		resp, err := s.polyAPI.GetScript(s.ctx, req)
// 		fmt.Printf("get %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	{
// 		req := &PolyRequestReq{
// 			ID:   newID,
// 			Body: []byte(`{}`),
// 		}
// 		resp, err := s.polyAPI.Request(s.ctx, req)
// 		fmt.Printf("request v1 %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	//--------------------------------------------------------------------------
// 	{
// 		time.Sleep(time.Second)
// 		req := &PolyBuildReq{
// 			ID:      newID,
// 			Arrange: `{"info":"build V3"}`,
// 		}
// 		resp, err := s.polyAPI.Build(s.ctx, req)
// 		fmt.Printf("build V3 %+v %v\n", resp, err)
// 		assert.NotNil(s.T(), err)
// 		assert.Nil(s.T(), resp)
// 	}
// 	{
// 		req := &PolyGetArrangeReq{ID: newID}
// 		resp, err := s.polyAPI.GetArrange(s.ctx, req)
// 		fmt.Printf("get %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	//--------------------------------------------------------------------------
// 	{
// 		time.Sleep(time.Second)
// 		req := &PolyUpdateScriptReq{
// 			ID:      newID,
// 			Script:  `"script V2"`,
// 			Swagger: `"doc v1"`,
// 		}
// 		resp, err := s.polyAPI.UpdateScript(s.ctx, req)
// 		fmt.Printf("script V2 %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	{
// 		req := &PolyGetScriptReq{ID: newID}
// 		resp, err := s.polyAPI.GetScript(s.ctx, req)
// 		fmt.Printf("get %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	{
// 		req := &PolyRequestReq{
// 			ID:   newID,
// 			Body: []byte(`{}`),
// 		}
// 		resp, err := s.polyAPI.Request(s.ctx, req)
// 		fmt.Printf("request v2 %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	//--------------------------------------------------------------------------
// 	if true {
// 		time.Sleep(time.Second)
// 		req := &PolyDeleteReq{ID: newID}
// 		resp, err := s.polyAPI.Delete(s.ctx, req)
// 		fmt.Printf("delete %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// 	{
// 		req := &PolyGetArrangeReq{ID: newID}
// 		resp, err := s.polyAPI.GetArrange(s.ctx, req)
// 		fmt.Printf("get %+v %v\n", resp, err)
// 		assert.Nil(s.T(), err)
// 		assert.NotNil(s.T(), resp)
// 	}
// }

package restful

import (
	"github.com/quanxiang-cloud/polyapi/internal/service"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"

	"github.com/gin-gonic/gin"
)

// APIKey is the route for api key
type APIKey struct {
	s    service.KMSAgent
	gate *gate.Gate
}

// NewAPIKey create a platform api key provider
func NewAPIKey(conf *config.Config) (*APIKey, error) {
	svs, err := service.CreateKMSAgent(conf)
	if err != nil {
		return nil, err
	}
	r := &APIKey{
		gate: gate.NewGate(conf),
		s:    svs,
	}
	return r, nil
}

// Create create a platform api key
func (s *APIKey) Create(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

// Update update a platform api key info
func (s *APIKey) Update(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

// Active update a platform api key active status
func (s *APIKey) Active(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

// Delete delete a platform api key
func (s *APIKey) Delete(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

// List list platform api key
func (s *APIKey) List(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

// Query query platform api key
func (s *APIKey) Query(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

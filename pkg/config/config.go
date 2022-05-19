package config

import (
	"io/ioutil"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/quanxiang-cloud/cabin/logger"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"gopkg.in/yaml.v2"
)

// Conf config
var Conf *Config

// DefaultPath is default config file path
var DefaultPath = "./configs/polyapi.yaml"

// GateLimitRate config
type GateLimitRate struct {
	Enable        bool `yaml:"enable" validate:"required"`
	RatePerSecond int  `yaml:"ratePerSecond" validate:"required,min=1,max=100000"`
}

// GateAPIBlock config
type GateAPIBlock struct {
	Enable        bool  `yaml:"enable" validate:"required"`
	MaxAllowError int64 `yaml:"maxAllowError" validate:"required,min=1,max=100"`
	BlockSeconds  int64 `yaml:"blockSeconds" validate:"required,min=10,max=3600"`
	APITimeoutMS  int64 `yaml:"apiTimeoutMS" validate:"required,min=1"`
}

// GateIPBlock config
type GateIPBlock struct {
	Enable bool     `yaml:"enable" validate:"required"`
	White  []string `yaml:"white"`
	Black  []string `yaml:"black"`
}

// Gate config
type Gate struct {
	LimitRate GateLimitRate `yaml:"limitRate"`
	APIBlock  GateAPIBlock  `yaml:"apiBlock"`
	IPBlock   GateIPBlock   `yaml:"ipBlock"`
}

// Config presents config
type Config struct {
	Port       string          `yaml:"port"`
	PortInner  string          `yaml:"portInner"`
	Model      string          `yaml:"model"`
	MyHostBase string          `yaml:"myHostBase"`
	Gate       Gate            `yaml:"gate"`
	Authorize  AuthorizeConfig `yaml:"authorize"`
	Log        logger.Config   `yaml:"log"`
	Mysql      mysql2.Config   `yaml:"mysql"`
	Redis      redis2.Config   `yaml:"redis"`
}

// AuthorizeConfig presents URLs of auth config
type AuthorizeConfig struct {
	OauthToken OauthConfig `yaml:"oauthToken"`
	OauthKey   OauthConfig `yaml:"oauthKey"`
	Goalie     OauthConfig `yaml:"goalie"`
	FileServer OauthConfig `yaml:"fileServer"`
	AppAccess  OauthConfig `yaml:"appAccess"`
	AppAdmin   OauthConfig `yaml:"appAdmin"`
}

// OauthConfig config url for authorize
type OauthConfig struct {
	Addr         string        `yaml:"addr"`
	MaxIdleConns int           `yaml:"maxIdleConns"`
	Timeout      time.Duration `yaml:"timeout"`
}

// NewConfig load config file
func NewConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		return nil, err
	}

	if err := validator.New().Struct(Conf); err != nil {
		return nil, err
	}

	return Conf, nil
}

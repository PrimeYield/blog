package setting

import (
	"github.com/spf13/viper"
)

type ServerSetting struct {
	// RunMode string `json:"runmode"`
	Port string `json:"httpport"`
}

type DatabaseSetting struct {
	MongodbHost string
	MongodbPort string
	Mongodb_db  string
}

type Setting struct {
	vp *viper.Viper
}

func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("./config")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Setting{vp: vp}, nil
}

func (s *Setting) ReadSection(key string, value interface{}) error {
	err := s.vp.UnmarshalKey(key, value)
	if err != nil {
		return err
	}
	return nil
}

package utils

import (
	"fmt"

	"github.com/spf13/viper"
	log "github.com/sung1011/tk-log"
)

func GetConf(module, key string) interface{} {
	conf := viper.GetViper().Get(module).(map[string]interface{})
	v, exists := conf[key]
	if !exists {
		log.Erro(fmt.Sprintf("conf %s.%s not exists", module, key))
	}
	return v
}

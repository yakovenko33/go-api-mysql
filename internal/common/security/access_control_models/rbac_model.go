package access_control_models

import (
	"fmt"
	"os"
	"sync"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var (
	accessControlModel *casbin.Enforcer
	once               sync.Once
	errorConnect       error
)

func InitAccessControlModel(DB *gorm.DB, logger *zap.Logger) (*casbin.Enforcer, error) {
	once.Do(func() {
		adapter, err := gormadapter.NewAdapterByDB(DB)
		if err != nil {
			logger.Error(fmt.Sprintf("Error create addapter %s", err))
			errorConnect = err
			return
		}

		dir := os.Getenv("CONFIG_CASBIAN_PATH")
		accessControlModel, err = casbin.NewEnforcer(dir, adapter)
		if err != nil {
			logger.Error(fmt.Sprintf("Error load model %s", err))
			errorConnect = err
			return
		}
	})
	return accessControlModel, errorConnect
}

func InitAccessControlModelForConsole(DB *gorm.DB) (*casbin.Enforcer, error) {
	once.Do(func() {
		adapter, err := gormadapter.NewAdapterByDB(DB)
		if err != nil {
			errorConnect = err
			return
		}

		dir := os.Getenv("CONFIG_CASBIAN_PATH")
		accessControlModel, err = casbin.NewEnforcer(dir, adapter)
		if err != nil {
			errorConnect = err
			return
		}
	})
	return accessControlModel, errorConnect
}

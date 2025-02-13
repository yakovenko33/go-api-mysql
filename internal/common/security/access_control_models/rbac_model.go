package access_control_models

import (
	"fmt"
	"sync"

	gormadapter "github.com/casbin/gorm-adapter/v3"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	AccessControlModel *casbin.Enforcer
	once               sync.Once
)

func InitAccessControlModel(DB *gorm.DB, logger *zap.Logger) *casbin.Enforcer {
	once.Do(func() {
		adapter, err := gormadapter.NewAdapterByDB(DB)
		if err != nil {
			logger.Error(fmt.Sprintf("Error create addapter %s", err))
		}

		AccessControlModel, err = casbin.NewEnforcer("model.conf", adapter)
		if err != nil {
			logger.Error(fmt.Sprintf("Error load model %s", err))
		}

		/*err = AccessControlModel.LoadPolicy()
		if err != nil {
			logger.Error(fmt.Sprintf("Error load policy %s", err))
		}*/
	})
	return AccessControlModel
}

/*CREATE TABLE casbin_rule (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ptype VARCHAR(100) NOT NULL,
    v0 VARCHAR(100) DEFAULT NULL,
    v1 VARCHAR(100) DEFAULT NULL,
    v2 VARCHAR(100) DEFAULT NULL,
    v3 VARCHAR(100) DEFAULT NULL,
    v4 VARCHAR(100) DEFAULT NULL,
    v5 VARCHAR(100) DEFAULT NULL
);*/

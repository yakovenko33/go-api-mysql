package access_control_models

import (
	"fmt"
	"sync"

	logging "go-api-docker/internal/common/logging"

	gormadapter "github.com/casbin/gorm-adapter/v3"

	database "go-api-docker/internal/common/database"

	"github.com/casbin/casbin/v2"
)

var (
	AccessControlModel *casbin.Enforcer
	once               sync.Once
)

func InitAccessControlModel() {
	once.Do(func() {
		adapter, err := gormadapter.NewAdapterByDB(database.DB)
		if err != nil {
			logging.Logger.Error(fmt.Sprintf("Error create addapter %s", err))
		}

		AccessControlModel, err = casbin.NewEnforcer("model.conf", adapter)
		if err != nil {
			logging.Logger.Error(fmt.Sprintf("Error load model %s", err))
		}

		err = AccessControlModel.LoadPolicy()
		if err != nil {
			logging.Logger.Error(fmt.Sprintf("Error load policy %s", err))
		}
	})
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

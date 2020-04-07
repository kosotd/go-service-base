package oracle

import (
	"fmt"
	"github.com/kosotd/go-service-base/database/url/domain"
)

func OracleUrlBuilder(params domain.Params) string {
	return fmt.Sprintf("%s/%s@%s:%s/%s", params.User, params.Password, params.Host, params.Port, params.Database)
}

package presto

import (
	"fmt"
	"github.com/kosotd/go-service-base/database/url/domain"
)

func PrestoUrlBuilder(params domain.Params) string {
	return fmt.Sprintf("http://user@%s:%s", params.Host, params.Port)
}

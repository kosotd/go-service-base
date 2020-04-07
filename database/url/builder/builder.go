package builder

import (
	"github.com/kosotd/go-service-base/database/url/builder/oracle"
	"github.com/kosotd/go-service-base/database/url/builder/postgres"
	"github.com/kosotd/go-service-base/database/url/builder/presto"
	"github.com/kosotd/go-service-base/database/url/domain"
	"github.com/kosotd/go-service-base/utils"
	"github.com/pkg/errors"
)

type UrlBuilder func(params domain.Params) string

func GetUrlBuilder(dbType string) (UrlBuilder, error) {
	if utils.Equals(dbType, "oracle") {
		return oracle.OracleUrlBuilder, nil
	} else if utils.Equals(dbType, "postgres") {
		return postgres.PostgresUrlBuilder, nil
	} else if utils.Equals(dbType, "presto") {
		return presto.PrestoUrlBuilder, nil
	} else {
		return func(domain.Params) string { return "" }, errors.Errorf("unknown database type %s", dbType)
	}
}

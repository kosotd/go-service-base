package database

import (
	"regexp"

	ldomain "github.com/kosotd/go-service-base/database/domain"
	"github.com/kosotd/go-service-base/database/url/builder"
	"github.com/kosotd/go-service-base/database/url/domain"
	"github.com/kosotd/go-service-base/utils"
	"github.com/pkg/errors"
)

// !!!
var urlPattern = regexp.MustCompile(`(\w+);(\w+):([\d\w.!-]+)/(\w+)@([\d\w.-]+):(\d+)(/(\w+))?`)

func parseUrl(url string) (ldomain.Connection, error) {
	if !urlPattern.MatchString(url) {
		return ldomain.Connection{}, errors.Errorf("url %s doesn't match with pattern", url)
	}
	submatch := urlPattern.FindStringSubmatch(url)
	dbType := submatch[2]
	params := domain.Params{
		Host:     submatch[5],
		Port:     submatch[6],
		User:     submatch[3],
		Password: submatch[4],
		Database: submatch[8],
	}
	driver, err := getDriverName(dbType)
	if err != nil {
		return ldomain.Connection{}, errors.Wrapf(err, "database.parseUrl -> getDriverName(%s)", dbType)
	}

	urlBuilder, err := builder.GetUrlBuilder(dbType)
	if err != nil {
		return ldomain.Connection{}, errors.Wrapf(err, "database.parseUrl -> builder.GetUrlBuilder(%s)", dbType)
	}

	return ldomain.Connection{
		Name:   submatch[1],
		DbType: dbType,
		Driver: driver,
		Url:    urlBuilder(params),
	}, nil
}

func getDriverName(dbType string) (string, error) {
	if utils.Equals(dbType, "oracle") {
		return "goracle", nil
	} else if utils.Equals(dbType, "postgres") {
		return "postgres", nil
	} else if utils.Equals(dbType, "presto") {
		return "presto", nil
	} else {
		return "", errors.Errorf("unknown database type %s", dbType)
	}
}

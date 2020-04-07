package postgres

import (
	"fmt"
	"github.com/kosotd/go-service-base/database/url/domain"
)

func PostgresUrlBuilder(params domain.Params) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		params.Host, params.Port, params.User, params.Password, params.Database)
}

package params_type

import (
	"github.com/pkg/errors"
	"strings"
)

const (
	None = iota
	FormUrlencoded
	PathParams
	RawQuery
	JsonBody
)

func FromString(paramsType string) (int, error) {
	switch strings.Trim(strings.ToLower(paramsType), " ") {
	case "":
		return None, nil
	case "form-urlencoded":
		return FormUrlencoded, nil
	case "json-body":
		return JsonBody, nil
	case "path-params":
		return PathParams, nil
	case "query-params":
		return RawQuery, nil
	default:
		return -1, errors.New("unknown params type")
	}
}

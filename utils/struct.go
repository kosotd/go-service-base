package utils

import (
	"encoding/json"
	"github.com/pkg/errors"
)

func StructToMap(object interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, errors.Wrapf(err, "utils.StructToMap -> json.Marshal")
	}
	res := make(map[string]interface{})
	if err = json.Unmarshal(bytes, &res); err != nil {
		return nil, errors.Wrapf(err, "utils.StructToMap -> json.Unmarshal")
	}
	return res, nil
}

func MustStructToMap(object interface{}) map[string]interface{} {
	if res, err := StructToMap(object); err != nil {
		panic(err)
	} else {
		return res
	}
}

func StructSliceToMap(objects interface{}) ([]map[string]interface{}, error) {
	bytes, err := json.Marshal(objects)
	if err != nil {
		return nil, errors.Wrapf(err, "utils.StructSliceToMap -> json.Marshal")
	}
	res := make([]map[string]interface{}, 0)
	if err = json.Unmarshal(bytes, &res); err != nil {
		return nil, errors.Wrapf(err, "utils.StructSliceToMap -> json.Unmarshal")
	}
	return res, nil
}

func MustStructSliceToMap(objects interface{}) []map[string]interface{} {
	if res, err := StructSliceToMap(objects); err != nil {
		panic(err)
	} else {
		return res
	}
}

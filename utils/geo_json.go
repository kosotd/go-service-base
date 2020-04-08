package utils

import (
	"encoding/json"
	"github.com/go-spatial/geom/encoding/geojson"
	"github.com/go-spatial/geom/encoding/wkt"
	"github.com/pkg/errors"
)

func BuildFeatureCollection(props []map[string]interface{}) ([]byte, error) {
	collection := geojson.FeatureCollection{Features: make([]geojson.Feature, 0)}
	for _, prop := range props {
		geom, err := getStringProperty(prop, "geometry")
		if err != nil {
			LogError(errors.Wrapf(err, "utils.BuildFeatureCollection -> getStringProperty").Error())
			continue
		}
		geometry, err := wkt.DecodeString(geom)
		if err != nil {
			LogError(errors.Wrapf(err, "utils.BuildFeatureCollection -> wkt.DecodeString(%s)", geom).Error())
			continue
		}

		delete(prop, "geometry")
		feature := geojson.Feature{
			Geometry:   geojson.Geometry{Geometry: geometry},
			Properties: prop,
		}
		collection.Features = append(collection.Features, feature)
	}

	res, err := json.Marshal(collection)
	if err != nil {
		return nil, errors.Wrapf(err, "utils.BuildFeatureCollection -> json.Marshal")
	}

	return res, nil
}

func getStringProperty(prop map[string]interface{}, key string) (string, error) {
	val, ok := prop[key]
	if !ok {
		return "", errors.Errorf("utils.getStringProperty -> prop has no key: %s", key)
	}

	valStr, ok := val.(string)
	if !ok {
		return "", errors.Errorf("utils.getStringProperty -> prop[%s] is not string", key)
	}

	return valStr, nil
}

package utils

import (
	"encoding/json"
	"github.com/go-spatial/geom/encoding/geojson"
	"gotest.tools/assert"
	"testing"
)

func TestBuildFeatureCollection(t *testing.T) {
	var geom *string
	nnGeom := "POINT(20 20)"
	var pnnGeom = &nnGeom
	bytes, err := BuildFeatureCollection([]map[string]interface{}{
		{
			"geometry": "POINT(10 10)",
			"col1":     "val1",
			"col2":     1,
		},
		{
			"geometry": "LINESTRING(10 10, 20 20, 30 30)",
			"col3":     "val2",
			"col4":     2,
		},
		{
			"geometry": "POLYGON((10 10, 10 20, 20 20, 10 10))",
			"col1":     "val3",
			"col2":     3,
		},
		{
			"col1": "val4",
			"col2": 4,
		},
		{
			"geometry": geom,
			"col1":     "val4",
			"col2":     4,
		},
		{
			"geometry": pnnGeom,
			"col1":     "val5",
			"col2":     5,
		},
	})
	assert.NilError(t, err)

	var collection geojson.FeatureCollection
	err = json.Unmarshal(bytes, &collection)
	assert.NilError(t, err)

	assert.Equal(t, len(collection.Features), 4)
	assert.Equal(t, collection.Features[0].Properties["col1"], "val1")
	assert.Equal(t, collection.Features[0].Properties["col2"], 1.0)
	assert.Equal(t, collection.Features[1].Properties["col3"], "val2")
	assert.Equal(t, collection.Features[1].Properties["col4"], 2.0)
	assert.Equal(t, collection.Features[2].Properties["col1"], "val3")
	assert.Equal(t, collection.Features[2].Properties["col2"], 3.0)
	assert.Equal(t, collection.Features[3].Properties["col1"], "val5")
	assert.Equal(t, collection.Features[3].Properties["col2"], 5.0)

	bytes = MustBuildFeatureCollection([]map[string]interface{}{
		{
			"geometry": "POINT(10 10)",
			"col1":     "val1",
			"col2":     1,
		},
		{
			"geometry": "LINESTRING(10 10, 20 20, 30 30)",
			"col3":     "val2",
			"col4":     2,
		},
		{
			"geometry": "POLYGON((10 10, 10 20, 20 20, 10 10))",
			"col1":     "val3",
			"col2":     3,
		},
		{
			"col1": "val4",
			"col2": 4,
		},
	})
	err = json.Unmarshal(bytes, &collection)
	assert.NilError(t, err)
	assert.Equal(t, len(collection.Features), 3)
}

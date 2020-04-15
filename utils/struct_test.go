package utils

import (
	"gotest.tools/assert"
	"testing"
	"time"
)

type Object struct {
	Id      int64     `json:"id"`
	Text    string    `json:"text"`
	TextPtr *string   `json:"text_ptr"`
	Date    time.Time `json:"date"`
}

func TestStructToMap(t *testing.T) {
	tm, err := time.Parse("2006-01-02", "2016-06-22")
	assert.NilError(t, err)
	text := "textPtr1"
	obj := Object{
		Id:      1,
		Text:    "text1",
		TextPtr: &text,
		Date:    tm,
	}

	objMap, err := StructToMap(obj)
	assert.NilError(t, err)
	assert.Equal(t, objMap["id"].(float64), float64(1))
	assert.Equal(t, objMap["text"].(string), "text1")
	assert.Equal(t, objMap["text_ptr"].(string), "textPtr1")
	assert.Equal(t, objMap["date"].(string), "2016-06-22T00:00:00Z")

	objMap = MustStructToMap(obj)
	assert.Equal(t, objMap["id"].(float64), float64(1))
	assert.Equal(t, objMap["text"].(string), "text1")
	assert.Equal(t, objMap["text_ptr"].(string), "textPtr1")
	assert.Equal(t, objMap["date"].(string), "2016-06-22T00:00:00Z")
}

func TestStructSliceToMap(t *testing.T) {
	tm, err := time.Parse("2006-01-02", "2016-06-22")
	assert.NilError(t, err)
	text := "textPtr1"
	slice := []Object{
		{
			Id:      1,
			Text:    "text1",
			TextPtr: &text,
			Date:    tm,
		},
		{
			Id:      2,
			Text:    "text2",
			TextPtr: nil,
			Date:    tm,
		},
	}

	objMap, err := StructSliceToMap(slice)
	assert.NilError(t, err)
	assert.Equal(t, objMap[0]["id"].(float64), float64(1))
	assert.Equal(t, objMap[0]["text"].(string), "text1")
	assert.Equal(t, objMap[0]["text_ptr"].(string), "textPtr1")
	assert.Equal(t, objMap[0]["date"].(string), "2016-06-22T00:00:00Z")

	objMap = MustStructSliceToMap(slice)
	assert.Equal(t, objMap[1]["id"].(float64), float64(2))
	assert.Equal(t, objMap[1]["text"].(string), "text2")
	assert.Equal(t, objMap[1]["text_ptr"], nil)
	assert.Equal(t, objMap[1]["date"].(string), "2016-06-22T00:00:00Z")
}

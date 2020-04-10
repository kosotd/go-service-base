package utils

import (
	"bytes"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/kosotd/go-service-base/utils/domain"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestExcelCol(t *testing.T) {
	assert.Equal(t, col(1), "A")
	assert.Equal(t, col(2), "B")
	assert.Equal(t, col(26), "Z")
	assert.Equal(t, col(27), "AA")
	assert.Equal(t, col(28), "AB")
	assert.Equal(t, col(52), "AZ")
	assert.Equal(t, col(53), "BA")
	assert.Equal(t, col(54), "BB")
}

func TestBuildExcelFile1(t *testing.T) {
	data := domain.ExcelData{
		Columns: []string{"col1 asdasdasdasd asdasdasdasd asdasdasdasdasdas dasdasdasdasd asdasd asdads", "col2", "col3"},
		Values: [][]interface{}{
			{"q", 1, 10 * time.Second},
			{"w", 2, 11 * time.Second},
			{"e", 3, 12 * time.Second},
			{"r", 4, 13 * time.Second},
		},
	}
	excelData, err := BuildExcelFile(data)
	assert.NilError(t, err)

	file, err := excelize.OpenReader(bytes.NewReader(excelData))
	assert.NilError(t, err)

	a1, err := file.GetCellValue("Sheet1", "A1")
	assert.NilError(t, err)
	assert.Equal(t, a1, "col1 asdasdasdasd asdasdasdasd asdasdasdasdasdas dasdasdasdasd asdasd asdads")

	b1, err := file.GetCellValue("Sheet1", "B1")
	assert.NilError(t, err)
	assert.Equal(t, b1, "col2")

	c1, err := file.GetCellValue("Sheet1", "C1")
	assert.NilError(t, err)
	assert.Equal(t, c1, "col3")

	a2, err := file.GetCellValue("Sheet1", "A2")
	assert.NilError(t, err)
	assert.Equal(t, a2, "q")

	b2, err := file.GetCellValue("Sheet1", "B2")
	assert.NilError(t, err)
	assert.Equal(t, b2, "1")

	c2, err := file.GetCellValue("Sheet1", "C2")
	assert.NilError(t, err)
	assert.Equal(t, c2, fmt.Sprintf("0.00011574074"))

	excelData = MustBuildExcelFile(data)
	file, err = excelize.OpenReader(bytes.NewReader(excelData))
	assert.NilError(t, err)

	a1, err = file.GetCellValue("Sheet1", "A1")
	assert.NilError(t, err)
	assert.Equal(t, a1, "col1 asdasdasdasd asdasdasdasd asdasdasdasdasdas dasdasdasdasd asdasd asdads")
}

func TestBuildExcelFile2(t *testing.T) {
	data := domain.ExcelData{
		Columns: []string{},
		Values:  [][]interface{}{},
	}
	_, err := BuildExcelFile(data)
	assert.NilError(t, err)
}

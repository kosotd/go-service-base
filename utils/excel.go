package utils

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/kosotd/go-service-base/utils/domain"
	"github.com/pkg/errors"
	"golang.org/x/text/width"
)

const alphabetSize = 26

var alphabet = make([]rune, alphabetSize)

func init() {
	for i := 0; i < alphabetSize; i++ {
		alphabet[i] = rune(i + 65)
	}
}

func BuildExcelFile(data domain.ExcelData) ([]byte, error) {
	file := excelize.NewFile()

	styleText := `{"alignment":{"wrap_text":true,"horizontal":"center","vertical":"center"}}`
	valueStyle, err := file.NewStyle(styleText)
	if err != nil {
		return nil, errors.Wrapf(err, "utils.BuildExcelFile -> file.NewStyle(%s)", styleText)
	}

	err = setHeaders(file, data.Columns)
	if err != nil {
		return nil, errors.Wrapf(err, "utils.BuildExcelFile -> setHeaders")
	}

	for rowNum := range data.Values {
		values := data.Values[rowNum]
		err := setRowValues(file, rowNum+2, values, valueStyle)
		if err != nil {
			return nil, errors.Wrapf(err, "utils.BuildExcelFile -> setRowValues")
		}
	}

	for i := range data.Columns {
		column := fmt.Sprintf("%s", col(i+1))
		columnWidth, err := getColumnWidth(file, len(data.Values)+1, i+1)
		if err != nil {
			return nil, errors.Wrapf(err, "utils.BuildExcelFile -> getColumnWidth")
		}
		err = file.SetColWidth("Sheet1", column, column, columnWidth)
		if err != nil {
			return nil, errors.Wrapf(err, "utils.BuildExcelFile -> file.SetColWidth")
		}
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, errors.Wrapf(err, "utils.BuildExcelFile -> file.WriteToBuffer")
	}

	return buffer.Bytes(), nil
}

func MustBuildExcelFile(data domain.ExcelData) []byte {
	if res, err := BuildExcelFile(data); err != nil {
		panic(err)
	} else {
		return res
	}
}

func setHeaders(file *excelize.File, columns []string) error {
	for i := range columns {
		axis := fmt.Sprintf("%s%d", col(i+1), 1)
		err := file.SetCellValue("Sheet1", axis, columns[i])
		if err != nil {
			return errors.Wrapf(err, "utils.setHeaders -> file.SetCellValue(%s, %s)", axis, columns[i])
		}
	}
	return nil
}

func setRowValues(file *excelize.File, row int, values []interface{}, style int) error {
	for i := range values {
		axis := fmt.Sprintf("%s%d", col(i+1), row)
		err := file.SetCellValue("Sheet1", axis, values[i])
		if err != nil {
			return errors.Wrapf(err, "utils.setRowValues -> file.SetCellValue(%s, %v)", axis, values[i])
		}
	}
	firstAxis := fmt.Sprintf("%s%d", col(1), row)
	lastAxis := fmt.Sprintf("%s%d", col(len(values)), row)
	err := file.SetCellStyle("Sheet1", firstAxis, lastAxis, style)
	if err != nil {
		return errors.Wrapf(err, "utils.setRowValues -> file.SetCellStyle(%s, %s)", firstAxis, lastAxis)
	}
	return nil
}

func getColumnWidth(file *excelize.File, rowCount int, column int) (float64, error) {
	max := -1.0
	for i := 1; i <= rowCount; i++ {
		axis := fmt.Sprintf("%s%d", col(column), i)
		value, err := file.GetCellValue("Sheet1", axis)
		if err != nil {
			return 0, errors.Wrapf(err, "utils.getColumnWidth -> file.GetCellValue")
		}

		if len(value) < 1 {
			continue
		}
		_, size := width.LookupString(value)
		curr := float64(size * len(value))
		if curr > max {
			max = curr
		}
	}
	return max, nil
}

func col(ind int) string {
	ind -= 1
	quotient := ind / alphabetSize
	if quotient > 0 {
		return col(quotient) + string(alphabet[ind%alphabetSize])
	} else {
		return string(alphabet[ind%alphabetSize])
	}
}

package util

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/dinesh-g1/csv-utility/consts"
	"github.com/dinesh-g1/csv-utility/types"
	"net/http"
	"strconv"
)

func GetCSVContent(r *http.Request) ([][]string, error) {
	if r.Method != http.MethodPost {
		return nil, &types.ApiError{
			Cause:      nil,
			Message:    "Method not allowed",
			StatusCode: http.StatusMethodNotAllowed,
		}
	}
	return ParseCSV(r)
}

func ParseCSV(r *http.Request) ([][]string, error) {
	records, err := getCSVContent(r)
	if err != nil {
		return nil, &types.ApiError{
			Cause:      err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}
	return records, nil
}

func getCSVContent(r *http.Request) ([][]string, error) {
	file, _, err := r.FormFile(consts.CSV_FILE_KEY)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}
	if err = ValidateCSV(records); err != nil {
		return nil, err
	}
	return records, nil
}

func ValidateCSV(records [][]string) error {
	if err := ValidateCSVOnRowColSize(records); err != nil {
		return err
	}
	if err := ValidateCSVCellValue(records); err != nil {
		return err
	}
	return nil
}

func ValidateCSVOnRowColSize(records [][]string) error {
	if len(records) < 1 {
		emptyCSVFile := fmt.Sprintf("given file has no content")
		return errors.New(emptyCSVFile)
	}
	if len(records) != len(records[0]) {
		sizeNotEqual := fmt.Sprintf("no of rows %d is not equal to no of columns %d", len(records), len(records[0]))
		return errors.New(sizeNotEqual)
	}

	return nil
}

func ValidateCSVCellValue(records [][]string) error {
	for rIdx, record := range records {
		for cIdx, val := range record {
			_, err := strconv.Atoi(val)
			if err != nil {
				cellValueNotIntErrMsg := fmt.Sprintf("cell(%d, %d) value is not an integer", rIdx, cIdx)
				return errors.New(cellValueNotIntErrMsg)
			}
		}
	}
	return nil
}

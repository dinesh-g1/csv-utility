package util

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/dinesh-g1/csv-utility/consts"
	"net/http"
	"strconv"
)

func GetCSVContentFromRequest(r *http.Request) ([][]string, error) {
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
	if len(records) != len(records[0]) {
		sizeNotEqualErrMsg := fmt.Sprintf("no of rows %d is not equal to no of columns %d", len(records), len(records[0]))
		return errors.New(sizeNotEqualErrMsg)
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

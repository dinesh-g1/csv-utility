package api

import (
	"fmt"
	"github.com/dinesh-g1/csv-utility/types"
	"github.com/dinesh-g1/csv-utility/util"
	"net/http"
	"strconv"
	"strings"
)

type RootHandler func(http.ResponseWriter, *http.Request) *types.ApiError

func (fn RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}
	//Handle different error cases
}

func Echo(w http.ResponseWriter, r *http.Request) *types.ApiError {
	records, err := util.GetCSVContentFromRequest(r)
	if err != nil {
		return &types.ApiError{
			Message:   "",
			Body:      nil,
			ErrorCode: 0,
		}
	}

	var response string
	for _, row := range records {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
	}
	fmt.Fprint(w, response)
	return nil
}

func Sum(w http.ResponseWriter, r *http.Request) *types.ApiError {
	records, err := util.GetCSVContentFromRequest(r)
	if err != nil {
		return &types.ApiError{
			Message:   "",
			Body:      nil,
			ErrorCode: 0,
		}
	}
	var sum int32
	for _, record := range records {
		for _, num := range record {
			intNum, err := strconv.Atoi(num)
			if err != nil {
				return &types.ApiError{
					Message:   "",
					Body:      nil,
					ErrorCode: 0,
				}
			}
			sum += int32(intNum)
		}
	}
	fmt.Fprint(w, "Sum of all records in the given CSV file is : ", sum)
	return nil
}

func Multiply(w http.ResponseWriter, r *http.Request) *types.ApiError {
	records, err := util.GetCSVContentFromRequest(r)
	if err != nil {
		return &types.ApiError{
			Message:   "",
			Body:      nil,
			ErrorCode: 0,
		}
	}
	var totalMul int64
	for _, record := range records {
		for _, num := range record {
			intNum, err := strconv.Atoi(num)
			if err != nil {
				return &types.ApiError{
					Message:   "",
					Body:      nil,
					ErrorCode: 0,
				}
			}
			totalMul *= int64(intNum)
		}
	}
	fmt.Fprint(w, "Multiplication of all records in the given CSV file is : ", totalMul)
	return nil
}

func Invert(w http.ResponseWriter, r *http.Request) *types.ApiError {
	records, err := util.GetCSVContentFromRequest(r)
	if err != nil {
		return &types.ApiError{
			Message:   "",
			Body:      nil,
			ErrorCode: 0,
		}
	}

	for i := 0; i < len(records); i++ {
		for j := 0; j < len(records[0]); j++ {
			if i < j {
				temp := records[i][j]
				records[i][j] = records[j][i]
				records[j][i] = temp
			}
		}
	}

	var response string
	for _, row := range records {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
	}
	fmt.Fprint(w, response)
	return nil
}

func Flatten(w http.ResponseWriter, r *http.Request) *types.ApiError {
	records, err := util.GetCSVContentFromRequest(r)
	if err != nil {
		return &types.ApiError{
			Message:   "",
			Body:      nil,
			ErrorCode: 0,
		}
	}
	var responseBuilder strings.Builder
	for _, row := range records {
		for _, val := range row {
			responseBuilder.WriteString(val)
			responseBuilder.WriteString(",")
		}
	}
	response := strings.TrimSuffix(responseBuilder.String(), ",")
	fmt.Fprint(w, response)
	return nil
}

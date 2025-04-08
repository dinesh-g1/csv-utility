package api

import (
	"errors"
	"fmt"
	"github.com/dinesh-g1/csv-utility/types"
	"github.com/dinesh-g1/csv-utility/util"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RootHandler func(http.ResponseWriter, *http.Request) error

func (fn RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}
	//Logging for application debugging
	log.Printf("error occured while processing the endpoint %s: %+v", r.RequestURI, err)

	var apiErr types.Error
	ok := errors.As(err, &apiErr)
	if ok {
		body, jErr := apiErr.ErrorMessage()
		if jErr != nil {
			log.Printf("error while marshalling the respone %v", jErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		statusCode, resHeaders := apiErr.ErrorStatusCode()
		for k, v := range resHeaders {
			w.Header().Set(k, v)
		}
		w.WriteHeader(statusCode)
		_, err = w.Write(body)
		if err != nil {
			log.Printf("error while writing the response %v", err)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("error while writing the response %v", err)
			return
		}
	}
}

func Echo(w http.ResponseWriter, r *http.Request) error {
	records, err := util.GetCSVContent(r)
	if err != nil {
		return err
	}

	var echoStr string
	for _, row := range records {
		echoStr = fmt.Sprintf("%s%s\n", echoStr, strings.Join(row, ","))
	}

	response := types.SuccessResponse{
		Value:      echoStr,
		StatusCode: http.StatusOK,
	}

	err = util.SendResponse(w, &response)
	if err != nil {
		return err
	}
	return nil
}

func Sum(w http.ResponseWriter, r *http.Request) error {
	records, err := util.GetCSVContent(r)
	if err != nil {
		return err
	}
	var sum int32
	for _, record := range records {
		for _, num := range record {
			intNum, err := strconv.Atoi(num)
			if err != nil {
				return &types.ApiError{
					Cause:      err,
					Message:    err.Error(),
					StatusCode: http.StatusInternalServerError,
				}
			}
			sum += int32(intNum)
		}
	}

	response := types.SuccessResponse{
		Value:      fmt.Sprintf("%d", sum),
		StatusCode: http.StatusOK,
	}

	err = util.SendResponse(w, &response)
	if err != nil {
		return err
	}
	return nil
}

func Multiply(w http.ResponseWriter, r *http.Request) error {
	records, err := util.GetCSVContent(r)
	if err != nil {
		return err
	}
	var totalMul int64 = 1
	for _, record := range records {
		for _, num := range record {
			intNum, err := strconv.Atoi(num)
			if err != nil {
				return &types.ApiError{
					Cause:      err,
					Message:    err.Error(),
					StatusCode: http.StatusInternalServerError,
				}
			}
			totalMul *= int64(intNum)
		}
	}

	response := types.SuccessResponse{
		Value:      fmt.Sprintf("%d", totalMul),
		StatusCode: http.StatusOK,
	}

	err = util.SendResponse(w, &response)
	if err != nil {
		return err
	}
	return nil
}

func Invert(w http.ResponseWriter, r *http.Request) error {
	records, err := util.GetCSVContent(r)
	if err != nil {
		return err
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

	var inverted string
	for _, row := range records {
		inverted = fmt.Sprintf("%s%s\n", inverted, strings.Join(row, ","))
	}

	response := types.SuccessResponse{
		Value:      inverted,
		StatusCode: http.StatusOK,
	}

	err = util.SendResponse(w, &response)
	if err != nil {
		return err
	}
	return nil
}

func Flatten(w http.ResponseWriter, r *http.Request) error {
	records, err := util.GetCSVContent(r)
	if err != nil {
		return err
	}
	var responseBuilder strings.Builder
	for _, row := range records {
		for _, val := range row {
			responseBuilder.WriteString(val)
			responseBuilder.WriteString(",")
		}
	}
	flattened := strings.TrimSuffix(responseBuilder.String(), ",")

	response := types.SuccessResponse{
		Value:      flattened,
		StatusCode: http.StatusOK,
	}

	err = util.SendResponse(w, &response)
	if err != nil {
		return err
	}
	return nil
}

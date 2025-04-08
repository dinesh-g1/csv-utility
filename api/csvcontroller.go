package api

import (
	"errors"
	"fmt"
	"github.com/dinesh-g1/csv-utility/types"
	"github.com/dinesh-g1/csv-utility/util"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

// RootHandler wraps all the other handlers to provide centralized error handling capabilities
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

	var response string
	for _, row := range records {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
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

	var limit big.Int
	limit.Exp(big.NewInt(10), big.NewInt(99), nil)

	sum := big.NewInt(0)
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

			if sum.Cmp(&limit) < 0 {
				sum.Add(sum, big.NewInt(int64(intNum)))
			} else {
				return &types.ApiError{
					Cause:      nil,
					Message:    fmt.Sprintf("sum is larger than %s", limit.String()),
					StatusCode: http.StatusInternalServerError,
				}
			}
		}
	}

	sumStr := sum.String()

	err = util.SendResponse(w, &sumStr)
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

	var limit big.Int
	limit.Exp(big.NewInt(10), big.NewInt(99), nil)

	multiply := big.NewInt(1)
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
			if multiply.Cmp(&limit) < 0 {
				multiply.Mul(multiply, big.NewInt(int64(intNum)))
			} else {
				return &types.ApiError{
					Cause:      nil,
					Message:    fmt.Sprintf("multiplicated value is larger than %s", limit.String()),
					StatusCode: http.StatusInternalServerError,
				}
			}

		}
	}
	totalMulStr := multiply.String()

	err = util.SendResponse(w, &totalMulStr)
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

	err = util.SendResponse(w, &inverted)
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

	err = util.SendResponse(w, &flattened)
	if err != nil {
		return err
	}
	return nil
}

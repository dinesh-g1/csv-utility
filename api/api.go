package api

import (
	"errors"
	"fmt"
	"github.com/dinesh-g1/csv-utility/types"
	"github.com/dinesh-g1/csv-utility/util"
	"log"
	"math/big"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// CSVOperationsHandler wraps all the other handlers to provide centralized error handling capabilities
type CSVOperationsHandler func(http.ResponseWriter, *http.Request) error

func (fn CSVOperationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

// Echo takes csv file records as input and simply returns the echo of input
func Echo(w http.ResponseWriter, r *http.Request) error {
	records, err := util.GetCSVContent(r)
	if err != nil {
		return err
	}

	var response string
	for _, row := range records {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
	}

	util.SendResponse(w, &response)
	return nil
}

// Sum receives the records from csv and returns the total sum of all values present in the file
func Sum(w http.ResponseWriter, r *http.Request) error {
	records, err := util.GetCSVContent(r)
	if err != nil {
		return err
	}

	sum, err := csvSum(records)
	if err != nil {
		return err
	}
	sumStr := sum.String()

	util.SendResponse(w, &sumStr)
	return nil
}

func csvSum(records [][]string) (*big.Int, error) {
	if len(records) < 50 {
		return sequentialSum(records)
	} else {
		return parallelSum(records)
	}
}

func parallelSum(records [][]string) (*big.Int, error) {
	start := time.Now()
	cpuCores := runtime.NumCPU()
	recordsLen := len(records)
	sumChunks := make(chan *big.Int)

	for i := 0; i < cpuCores; i++ {
		from := (i * recordsLen) / cpuCores
		to := (i + 1) * recordsLen / cpuCores
		go func(records [][]string) {
			sum := big.NewInt(0)
			for _, record := range records {
				for _, num := range record {
					intNum, _ := strconv.Atoi(num)
					sum.Add(sum, big.NewInt(int64(intNum)))
				}
			}
			sumChunks <- sum
		}(records[from:to])
	}

	totalSum := big.NewInt(0)
	for i := 0; i < cpuCores; i++ {
		totalSum.Add(totalSum, <-sumChunks)
	}
	end := time.Now()
	log.Printf("time elapsed: %v", end.Sub(start))
	return totalSum, nil
}

func sequentialSum(records [][]string) (*big.Int, error) {
	start := time.Now()

	sum := big.NewInt(0)
	for _, record := range records {
		for _, num := range record {
			intNum, _ := strconv.Atoi(num)
			sum.Add(sum, big.NewInt(int64(intNum)))
		}
	}
	end := time.Now()
	log.Printf("time elapsed: %v", end.Sub(start))
	return sum, nil
}

// Multiply receives the records of csv and returns the multiplication of all values present in the file
// Since file can have many values, and multiplication result can be huge, the result can be upper bounded using limit
func Multiply(w http.ResponseWriter, r *http.Request) error {
	records, err := util.GetCSVContent(r)
	if err != nil {
		return err
	}

	//var limit big.Int
	//limit.Exp(big.NewInt(10), big.NewInt(99), nil)

	multiply := big.NewInt(1)
	for _, record := range records {
		for _, num := range record {
			intNum, _ := strconv.Atoi(num)
			multiply.Mul(multiply, big.NewInt(int64(intNum)))
		}
	}
	totalMulStr := multiply.String()

	util.SendResponse(w, &totalMulStr)
	return nil
}

// Invert inverts the received csv matrix
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

	util.SendResponse(w, &inverted)
	return nil
}

// Flatten takes the multiple rows of csv file and flattens (appends each row after another) the file.
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

	util.SendResponse(w, &flattened)
	return nil
}

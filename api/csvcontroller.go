package api

import (
	"fmt"
	"github.com/dinesh-g1/csv-utility/util"
	"net/http"
	"strconv"
	"strings"
)

func Echo(w http.ResponseWriter, r *http.Request) {
	records, err := util.GetCSVContentFromRequest(r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
		return
	}

	var response string
	for _, row := range records {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
	}
	fmt.Fprint(w, response)
}

func Sum(w http.ResponseWriter, r *http.Request) {
	records, err := util.GetCSVContentFromRequest(r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
		return
	}
	var sum int32
	for _, record := range records {
		for _, num := range record {
			intNum, err := strconv.Atoi(num)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
				return
			}
			sum += int32(intNum)
		}
	}
	fmt.Fprint(w, "Sum of all records in the given CSV file is : ", sum)
}

func Multiply(w http.ResponseWriter, r *http.Request) {
	records, err := util.GetCSVContentFromRequest(r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
		return
	}
	var totalMul int64
	for _, record := range records {
		for _, num := range record {
			intNum, err := strconv.Atoi(num)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
				return
			}
			totalMul *= int64(intNum)
		}
	}
	fmt.Fprint(w, "Multiplication of all records in the given CSV file is : ", totalMul)
}

func Invert(w http.ResponseWriter, r *http.Request) {
	records, err := util.GetCSVContentFromRequest(r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
		return
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
}

func Flatten(w http.ResponseWriter, r *http.Request) {
	records, err := util.GetCSVContentFromRequest(r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
		return
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
}

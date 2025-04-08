package api

import (
	"bytes"
	"encoding/json"
	"github.com/dinesh-g1/csv-utility/consts"
	"github.com/dinesh-g1/csv-utility/types"
	mux "github.com/gorilla/mux"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var r *mux.Router

func TestMain(m *testing.M) {
	r = mux.NewRouter()

	//Versioning the api
	v1Api := r.PathPrefix(consts.ApiV1).Subrouter()

	//Registering all the api endpoints
	v1Api.Handle(consts.EndpointEcho, CSVOperationsHandler(Echo)).Methods(http.MethodPost)
	v1Api.Handle(consts.EndpointInvert, CSVOperationsHandler(Invert)).Methods(http.MethodPost)
	v1Api.Handle(consts.EndpointFlatten, CSVOperationsHandler(Flatten)).Methods(http.MethodPost)
	v1Api.Handle(consts.EndpointSum, CSVOperationsHandler(Sum)).Methods(http.MethodPost)
	v1Api.Handle(consts.EndpointMultiply, CSVOperationsHandler(Multiply)).Methods(http.MethodPost)

	code := m.Run()

	r = nil
	os.Exit(code)
}

func TestEchoSuccess(t *testing.T) {
	t.Run("it returns the echo of given CSV", func(t *testing.T) {
		request := newRequest(consts.SampleCsvFileName, consts.EndpointEcho)
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBody(t, response.Body.String(), "1,2,3\n4,5,6\n7,8,9\n")
	})
}

func TestInvertSuccess(t *testing.T) {
	t.Run("it returns inverted csv matrix", func(t *testing.T) {
		request := newRequest(consts.SampleCsvFileName, consts.EndpointInvert)
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBody(t, response.Body.String(), "1,4,7\n2,5,8\n3,6,9\n")
	})
}

func TestFlattenSuccess(t *testing.T) {
	t.Run("it returns flattened csv records", func(t *testing.T) {
		request := newRequest(consts.SampleCsvFileName, consts.EndpointFlatten)
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBody(t, response.Body.String(), "1,2,3,4,5,6,7,8,9")
	})
}
func TestSumSequentialSuccess(t *testing.T) {
	t.Run("it returns sum of csv values", func(t *testing.T) {
		request := newRequest(consts.SampleCsvFileName, consts.EndpointSum)
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBody(t, response.Body.String(), "45")
	})
}

func TestSumSequentialError(t *testing.T) {
	t.Run("it returns sum of csv values", func(t *testing.T) {
		request := newRequest(consts.SampleCsvFileName, consts.EndpointSum)
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBody(t, response.Body.String(), "45")
	})
}

func TestSumParallelSuccess(t *testing.T) {
	t.Run("it returns sum of csv values", func(t *testing.T) {
		request := newRequest(consts.BigCsvFileName, consts.EndpointSum)
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBody(t, response.Body.String(), "505000")
	})
}

func TestMultiplySuccess(t *testing.T) {
	t.Run("it returns multiplication of csv values", func(t *testing.T) {
		request := newRequest(consts.SampleCsvFileName, consts.EndpointMultiply)
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBody(t, response.Body.String(), "362880")
	})
}

func TestEmptyCsvError(t *testing.T) {
	t.Run("it returns bad request as csv is empty", func(t *testing.T) {
		request := newRequest(consts.EmptyCsvFileName, consts.EndpointEcho)
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
		var errResponse types.ApiError
		err := json.NewDecoder(response.Body).Decode(&errResponse)
		if err != nil {
			log.Fatal(err)
		}
		assertBody(t, errResponse.Message, "given file has no content")
	})
}
func TestMethodNotAllowedError(t *testing.T) {
	t.Run("it returns not found error as http method is not post", func(t *testing.T) {
		request := newRequest(consts.EmptyCsvFileName, consts.EndpointEcho)
		request.Method = http.MethodGet
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func newRequest(csvFileName string, endPointEcho string) *http.Request {
	csv := new(bytes.Buffer)
	writer := multipart.NewWriter(csv)
	fWriter, err := writer.CreateFormFile(consts.CsvFileKey, filepath.Base(consts.CsvFileBasePath+csvFileName))
	if err != nil {
		log.Fatal(err)
	}

	fd, err := os.Open(consts.CsvFileBasePath + csvFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(fd *os.File) {
		err := fd.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(fd)
	_, err = io.Copy(fWriter, fd)
	if err != nil {
		log.Fatal(err)
	}
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, consts.ApiV1+endPointEcho, csv)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func assertStatus(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
}

func assertBody(t *testing.T, got, want string) {
	if !strings.EqualFold(got, want) {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

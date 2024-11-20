package rest

import (
	"encoding/json"
	"github.com/DanKo-code/TODO-list/pkg/logger"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	if err != nil {
		logger.ErrorLogger.Printf("Failed decode request body: %v", err)
		return err
	}
	return nil
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	encoder.Encode(response)
}

func WriteErrToResponseBody(w http.ResponseWriter, err error, status int) {
	w.Header().Add("Content-Type", "application/json;")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

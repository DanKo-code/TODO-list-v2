package rest

import (
	"encoding/json"
	"github.com/DanKo-code/TODO-list/pkg/logger"
	"net/http"
	"regexp"
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

func isValidUUID(uuid string) bool {
	regex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[1-5][a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)
	return regex.MatchString(uuid)
}

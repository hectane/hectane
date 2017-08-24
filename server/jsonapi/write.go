package jsonapi

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
)

func writeJson(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		// TODO: log the error
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// WriteError generates an error response with the specified message.
func WriteError(w http.ResponseWriter, e string) {
	writeJson(w, map[string]string{
		"error": e,
	})
}

// WriteData generates a JSON API response for a slice of the specified type.
func WriteData(w http.ResponseWriter, model string, v interface{}) {
	var (
		sliceVal = reflect.ValueOf(v)
		data     []interface{}
	)
	for i := 0; i < sliceVal.Len(); i++ {
		var (
			itemVal = reflect.Indirect(sliceVal.Index(i))
		)
		data = append(data, map[string]interface{}{
			"type":       model,
			"id":         itemVal.Field(0).Interface(),
			"attributes": itemVal.Interface(),
		})
	}
	writeJson(w, map[string]interface{}{
		"data": data,
	})
}

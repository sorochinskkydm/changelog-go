package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func JSON(writer http.ResponseWriter, status int, data any) {
	writer.Header().Set("Content-Type", "Application/json")
	writer.WriteHeader(status)

	if err := json.NewEncoder(writer).Encode(data); err != nil {
		slog.Error("write json response", "error", err)
	}
}
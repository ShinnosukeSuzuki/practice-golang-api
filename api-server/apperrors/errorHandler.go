package apperrors

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ShinnosukeSuzuki/practice-golang-api/api/middlewares"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	// エラーの種類を判別して、適切なHTTPステータスコードを返却
	var appErr *MyAppError
	if !errors.As(err, &appErr) {
		appErr = &MyAppError{
			ErrCode: Unknown,
			Message: "internal porcess failed",
			Err:     err,
		}
	}

	traceID := middlewares.GetTraceID(r.Context())
	log.Printf("[%d] ERROR: %+v\n", traceID, appErr)

	var statusCode int
	switch appErr.ErrCode {
	case NAData:
		statusCode = http.StatusNotFound
	case NoTargetData, ReqBodyDecodeFailed, BadParam:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(appErr)
}

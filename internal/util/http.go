package util

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// //////////////////////////////////////////////////
// extract

func ExtractPathParameter(ctx context.Context, name string) string {
	params := httprouter.ParamsFromContext(ctx)
	return params.ByName(name)
}

func ExtractParameter(r *http.Request, name string) string {
	return r.FormValue(name)
}

func ExtractPathIntParameter(ctx context.Context, name string) int {
	return ToInt(ExtractPathParameter(ctx, name))
}

func ExtractIntParameter(r *http.Request, name string) int {
	return ToInt(ExtractParameter(r, name))
}

func ExtractPathBoolParameter(ctx context.Context, name string) bool {
	return ToBool(ExtractPathParameter(ctx, name))
}

func ExtractBoolParameter(r *http.Request, name string) bool {
	return ToBool(ExtractParameter(r, name))
}

// //////////////////////////////////////////////////
// encode error

func EncodeJsonErrorResponse(resp http.ResponseWriter, err any) {
	resp.Header().Set("Content-Type", "application/json")
	if o, ok := err.(interface{ Status() int }); ok {
		resp.WriteHeader(o.Status())
	} else {
		resp.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(resp).Encode(ToJsonErrorResponse(err))
}

func ToJsonErrorResponse(err any) *JsonErrorResponse {
	return &JsonErrorResponse{
		Error: ToJsonError(err),
	}
}

func ToJsonError(err any) *JsonError {
	jsonError := &JsonError{}
	if o, ok := err.(interface{ ErrorCode() int }); ok {
		jsonError.Code = o.ErrorCode()
	}
	if o, ok := err.(interface{ ErrorMessage() string }); ok {
		jsonError.Message = o.ErrorMessage()
	} else if o, ok := err.(interface{ Error() string }); ok {
		jsonError.Message = o.Error()
	} else {
		jsonError.Message = "internal error"
	}
	if o, ok := err.(interface{ ErrorCause() any }); ok {
		jsonError.Cause = ToJsonError(o.ErrorCause())
	}
	return jsonError
}

type JsonErrorResponse struct {
	Error *JsonError `json:"error"`
}

type JsonError struct {
	Code    int        `json:"code,omitempty"`
	Message string     `json:"message"`
	Cause   *JsonError `json:"cause,omitempty"`
}

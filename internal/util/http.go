package util

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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

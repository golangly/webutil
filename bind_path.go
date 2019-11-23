package webutil

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi"
	"github.com/golangly/errors"
)

func bindRequestPathTag(r *http.Request, fieldName string, fieldValue reflect.Value, tag string) error {
	parameterName := ""
	required := false
	for _, token := range strings.Split(tag, ",") {
		if token == "required" {
			required = true
		} else if parameterName != "" {
			return errors.New("invalid tag: only one parameter name allowed per field").
				AddTag("fieldName", fieldName).
				AddTag("tag", tag)
		} else {
			parameterName = token
		}
	}

	var values []string
	value := chi.URLParam(r, parameterName)
	if value == "" {
		values = []string{}
	} else {
		values = []string{value}
	}

	if err := bindValues(values, fieldValue, required); err != nil {
		return errors.Wrap(err, "path bind error").
			AddTag("fieldName", fieldName).
			AddTag("tag", tag)
	} else {
		return nil
	}
}

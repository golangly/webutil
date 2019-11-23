package webutil

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/golangly/errors"
)

func bindRequestQueryTag(r *http.Request, fieldName string, fieldValue reflect.Value, tag string) error {
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

	if err := r.ParseForm(); err != nil {
		return errors.Wrap(err, "form parse error").
			AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest).
			AddTag("fieldName", fieldName).
			AddTag("tag", tag)
	}

	values := r.Form[parameterName]
	if values == nil {
		values = make([]string, 0)
	}

	if err := bindValues(values, fieldValue, required); err != nil {
		return errors.Wrap(err, "query bind error").
			AddTag("fieldName", fieldName).
			AddTag("tag", tag)
	} else {
		return nil
	}
}

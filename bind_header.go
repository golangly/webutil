package webutil

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/golangly/errors"
)

func bindRequestHeaderTag(r *http.Request, fieldName string, fieldValue reflect.Value, tag string) error {
	headerName := ""
	required := false
	for _, token := range strings.Split(tag, ",") {
		if token == "required" {
			required = true
		} else if headerName != "" {
			return errors.New("invalid tag: only one header name allowed per field").
				AddTag("fieldName", fieldName).
				AddTag("tag", tag)
		} else {
			headerName = token
		}
	}

	if err := bindValues(r.Header[headerName], fieldValue, required); err != nil {
		return errors.Wrap(err, "bind error").
			AddTag("fieldName", fieldName).
			AddTag("tag", tag)
	} else {
		return nil
	}
}

package webutil

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/golangly/errors"
)

func bindRequestCookieTag(r *http.Request, fieldName string, fieldValue reflect.Value, tag string) error {
	cookieName := ""
	required := false
	for _, token := range strings.Split(tag, ",") {
		if token == "required" {
			required = true
		} else if cookieName != "" {
			return errors.New("invalid tag: only one cookie name allowed per field").
				AddTag("fieldName", fieldName).
				AddTag("tag", tag)
		} else {
			cookieName = token
		}
	}

	if cookie, err := r.Cookie(cookieName); err != nil {
		return errors.Wrap(err, "cookie parse error").
			AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest).
			AddTag("fieldName", fieldName).
			AddTag("tag", tag)
	} else if err := bindValues([]string{cookie.Value}, fieldValue, required); err != nil {
		return errors.Wrap(err, "cookie bind error").
			AddTag("fieldName", fieldName).
			AddTag("tag", tag)
	} else {
		return nil
	}
}

package webutil

import (
	"net/http"
	"reflect"

	"github.com/golangly/errors"
)

func Bind(r *http.Request, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return errors.New("bind target must be a struct pointer").AddTag("type", targetValue.Type())
	}
	targetValue = targetValue.Elem()
	targetType := targetValue.Type()
	for i := 0; i < targetValue.NumField(); i++ {
		fieldValue := targetValue.Field(i)
		fieldType := targetType.Field(i)
		if headerSpec := fieldType.Tag.Get("header"); headerSpec != "" {
			if err := bindRequestHeaderTag(r, fieldType.Name, fieldValue, headerSpec); err != nil {
				return err
			}
		} else if pathSpec := fieldType.Tag.Get("path"); pathSpec != "" {
			if err := bindRequestPathTag(r, fieldType.Name, fieldValue, pathSpec); err != nil {
				return err
			}
		} else if querySpec := fieldType.Tag.Get("query"); querySpec != "" {
			if err := bindRequestQueryTag(r, fieldType.Name, fieldValue, querySpec); err != nil {
				return err
			}
		} else if cookieSpec := fieldType.Tag.Get("cookie"); cookieSpec != "" {
			if err := bindRequestCookieTag(r, fieldType.Name, fieldValue, cookieSpec); err != nil {
				return err
			}
		} else if bodySpec, ok := fieldType.Tag.Lookup("body"); ok {
			if err := bindRequestBodyTag(r, fieldType.Name, fieldValue, bodySpec); err != nil {
				return err
			}
		}
	}

	return validateTarget(r, target)
}

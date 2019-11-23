package webutil

import (
	"encoding/json"
	"encoding/xml"
	"mime"
	"net/http"
	"reflect"

	"github.com/golangly/errors"
	"gopkg.in/yaml.v2"
)

func bindRequestBodyTag(r *http.Request, fieldName string, fieldValue reflect.Value, tag string) error {
	if tag != "" {
		return errors.New("invalid tag: body tag must be empty").
			AddTag("fieldName", fieldName).
			AddTag("tag", tag)
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return errors.Wrap(err, "form parse error").
			AddTag(ErrTagHTTPStatusCodeKey, http.StatusUnsupportedMediaType).
			AddTag("fieldName", fieldName).
			AddTag("contentType", r.Header.Get("Content-Type"))
	}

	switch mediaType {
	case "application/x-yaml", "application/yaml", "text/yaml":
		decoder := yaml.NewDecoder(r.Body)
		decoder.SetStrict(false)
		if err := decoder.Decode(fieldValue); err != nil {
			return errors.Wrap(err, "decoding failed").
				AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest).
				AddTag("fieldName", fieldName).
				AddTag("contentType", r.Header.Get("Content-Type"))
		}

	case "application/json", "text/json":
		decoder := json.NewDecoder(r.Body)
		decoder.UseNumber()
		decoder.DisallowUnknownFields()
		vptr := reflect.New(fieldValue.Type())
		if err := decoder.Decode(vptr.Interface()); err != nil {
			return errors.Wrap(err, "decoding failed").
				AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest).
				AddTag("fieldName", fieldName).
				AddTag("contentType", r.Header.Get("Content-Type"))
		} else {
			fieldValue.Set(vptr.Elem())
		}

	case "application/xml", "text/xml":
		decoder := xml.NewDecoder(r.Body)
		if err := decoder.Decode(fieldValue); err != nil {
			return errors.Wrap(err, "decoding failed").
				AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest).
				AddTag("fieldName", fieldName).
				AddTag("contentType", r.Header.Get("Content-Type"))
		}

	default:
		return errors.New("decoding failed").
			AddTag(ErrTagHTTPStatusCodeKey, http.StatusUnsupportedMediaType).
			AddTag("fieldName", fieldName).
			AddTag("contentType", r.Header.Get("Content-Type"))
	}
	return nil
}

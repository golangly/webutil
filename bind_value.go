package webutil

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/golangly/errors"
)

func bindValues(values []string, v reflect.Value, required bool) error {
	switch v.Kind() {
	case reflect.Bool:
		if len(values) == 0 {
			if required {
				return errors.New("required value missing").
					AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest)
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else if b, err := strconv.ParseBool(values[0]); err != nil {
			return errors.Wrap(err, "bad boolean value").
				AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest)
		} else {
			v.SetBool(b)
		}
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		if len(values) == 0 {
			if required {
				return errors.New("required value missing").
					AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest)
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else if i, err := strconv.ParseInt(values[0], 10, 64); err != nil {
			return errors.Wrap(err, "bad int value").
				AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest)
		} else {
			v.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		if len(values) == 0 {
			if required {
				return errors.New("required value missing").
					AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest)
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else if i, err := strconv.ParseUint(values[0], 10, 64); err != nil {
			return errors.Wrap(err, "bad uint value").
				AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest)
		} else {
			v.SetUint(i)
		}
	case reflect.Float32, reflect.Float64:
		if len(values) == 0 {
			if required {
				return errors.New("required value missing").
					AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest)
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else if f, err := strconv.ParseFloat(values[0], 64); err != nil {
			return errors.Wrap(err, "bad float value").
				AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest)
		} else {
			v.SetFloat(f)
		}
	case reflect.Ptr:
		if len(values) == 0 {
			if required {
				return errors.New("required value missing").
					AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest)
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else {
			targetType := v.Type().Elem()
			targetValue := reflect.New(targetType)
			if err := bindValues(values, targetValue, true); err != nil {
				return err
			} else {
				v.Set(targetValue)
			}
		}
	case reflect.String, reflect.Interface:
		if len(values) == 0 {
			if required {
				return errors.New("required value missing").
					AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest)
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else {
			v.SetString(values[0])
		}
	case reflect.Slice:
		v.SetLen(len(values))
		v.SetCap(len(values))
		for i, value := range values {
			vi := v.Index(i)
			if err := bindValues([]string{value}, vi, true); err != nil {
				return err
			}
		}
	default:
		panic("not implemented")
	}
	return nil
}

package webutil

import (
	"net/http"

	"github.com/golangly/errors"
	"gopkg.in/go-playground/validator.v9"
)

func validateTarget(r *http.Request, target interface{}) error {
	v := r.Context().Value("validator").(*validator.Validate)
	if err := v.StructCtx(r.Context(), target); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			// TODO: send back validation errors to client
			return errors.Wrap(err, "validation failed").
				AddTag(ErrTagHTTPStatusCodeKey, http.StatusBadRequest).
				AddTag("validationError", ve.Error())
		} else {
			return errors.Wrap(err, "validation error").
				AddTag("validationError", ve.Error())
		}
	}
	return nil
}

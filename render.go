package webutil

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/golangly/errors"
	"github.com/golangly/log"
	"gopkg.in/yaml.v2"
)

const ErrTagHTTPStatusCodeKey = "StatusCode"
const ErrTypeHTTPPublic = "Public"

var (
	OfferedContentTypes = []string{
		"application/x-yaml",
		"application/yaml",
		"text/yaml",
		"application/json",
		"text/json",
		"application/xml",
		"text/xml",
		"text/html",
		"text/plain",
	}
)

func RenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, v interface{}) {
	if err, ok := v.(error); ok {
		if errorStatusCode, ok := errors.LookupTag(err, ErrTagHTTPStatusCodeKey).(int); ok {
			statusCode = errorStatusCode
		}
		if statusCode == http.StatusInternalServerError {
			log.WithErr(err).Error("Internal server error rendered to HTTP client")
		}
		if errors.HasType(err, ErrTypeHTTPPublic) {
			v = err.Error()
		} else {
			v = http.StatusText(statusCode)
		}
	}
	w.WriteHeader(statusCode)
	if v != nil {
		Render(w, r, v)
	}
}

func Render(w http.ResponseWriter, r *http.Request, v interface{}) {
	switch acceptedMimeType := NegotiateContentType(r.Header.Get("Accept"), OfferedContentTypes, ""); acceptedMimeType {
	case "application/x-yaml", "application/yaml", "text/yaml":
		w.Header().Set("Content-Type", acceptedMimeType)
		encoder := yaml.NewEncoder(w)
		defer encoder.Close()
		if err := encoder.Encode(v); err != nil {
			log.With("mimeType", acceptedMimeType).
				With("data", v).
				Warn("Failed encoding data to response")
		}

	case "application/json", "text/json":
		w.Header().Set("Content-Type", acceptedMimeType)
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(v); err != nil {
			log.With("mimeType", acceptedMimeType).
				With("data", v).
				Warn("Failed encoding data to response")
		}

	case "application/xml", "text/xml":
		w.Header().Set("Content-Type", acceptedMimeType)
		encoder := xml.NewEncoder(w)
		if err := encoder.Encode(v); err != nil {
			log.With("mimeType", acceptedMimeType).
				With("data", v).
				Warn("Failed encoding data to response")
		} else if err := encoder.Flush(); err != nil {
			log.With("mimeType", acceptedMimeType).
				With("data", v).
				Warn("Failed encoding data to response")
		}

	case "text/html":
		w.WriteHeader(http.StatusOK) // send HTTP 200 since this is most probably a browser
		w.Header().Set("Content-Type", acceptedMimeType)
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		if _, err := w.Write([]byte("<!DOCTYPE html><html><body><pre>")); err != nil {
			log.With("mimeType", acceptedMimeType).
				With("data", v).
				Warn("Failed encoding data to response")
		} else if err := encoder.Encode(v); err != nil {
			log.With("mimeType", acceptedMimeType).
				With("data", v).
				Warn("Failed encoding data to response")
		} else if _, err := w.Write([]byte("</pre></body></html>")); err != nil {
			log.With("mimeType", acceptedMimeType).
				With("data", v).
				Warn("Failed encoding data to response")
		}

	case "text/plain":
		w.Header().Set("Content-Type", acceptedMimeType)
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(v); err != nil {
			log.With("mimeType", acceptedMimeType).
				With("data", v).
				Warn("Failed encoding data to response")
		}

	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

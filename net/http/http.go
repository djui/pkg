package http

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

// HTTPError holds an error and HTTP status.
type HTTPError struct {
	error
	Status int
}

// ToHTTPError takes an error and returns a best guess of what the corresponding
// HTTP error should be.
func ToHTTPError(err error) (msg string, httpStatus int) {
	if httpErr, ok := err.(*HTTPError); ok {
		return httpErr.Error(), httpErr.Status
	}
	if os.IsNotExist(err) {
		return fmt.Sprintf("%d %s", http.StatusNotFound,
			http.StatusText(http.StatusNotFound)), http.StatusNotFound
	}
	if os.IsPermission(err) {
		return fmt.Sprintf("%d %s", http.StatusForbidden,
			http.StatusText(http.StatusForbidden)), http.StatusForbidden
	}
	return err.Error(), http.StatusInternalServerError
}

// ContentDisposition constructs a content-disposition header.
func ContentDisposition(filepath string) string {
	return fmt.Sprint("attachment; filename=", strconv.Quote(path.Base(filepath)))
}

// DeleteCookie effectively delete a cookie.
func DeleteCookie(w http.ResponseWriter, n string) {
	c := &http.Cookie{
		Name:    n,
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Time{},
	}
	http.SetCookie(w, c)
}

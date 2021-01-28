// Package request contains methods for handling http.Request
package request

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// ParseBody parses the request body into output variable
func ParseBody(r *http.Request, output interface{}) error {
	// Get the contentType for comparisons
	ct := r.Header.Get("Content-Type")

	if strings.Contains(ct, "application/json") {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &output)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("header content-type is different from application/json")
}

package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func DecodeQuery(dst interface{}, r *http.Request) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(dst, r.URL.Query()); err != nil {
		return fmt.Errorf(`query parse error: %w`, err)
	}

	return nil
}

func DecodeBody(v interface{}, r *http.Request) error {
	if r.Body == nil {
		return nil
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		if err == io.EOF {
			return nil
		}
		return fmt.Errorf(`query parse error: %w`, err)
	}

	return nil
}

func GetVarString(r *http.Request, Key string) (string, error) {
	vars := mux.Vars(r)
	if res, ok := vars[Key]; ok {
		return res, nil
	}

	return "", errors.New(`request parse error: bad request`)
}

func GetVarInt(r *http.Request, Key string) (int64, error) {
	str, err := GetVarString(r, Key)
	if err != nil {
		return 0, err
	}

	if res, err := strconv.ParseInt(str, 10, 64); err == nil {
		return res, nil
	}

	return 0, errors.New(`request parse error: bad request`)
}

// CreateGUID create new GUID
func CreateGUID() string {
	id, _ := uuid.NewUUID()
	return strings.ToUpper(strings.Replace(id.String(), "-", "", -1))
}

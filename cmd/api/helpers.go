package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/olzzhas/narxozer/internal/validator"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) writeJSONRedis(data any) (response any, err error) {
	js, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	js = append(js, '\n')

	return js, nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()
}

func (app *application) parseTimeStringLegacy(input string) (time.Time, error) {
	if len(input) != len("2006-01-02-15-04-05") {
		return time.Time{}, errors.New("incorrect time format")
	}

	year, err := strconv.Atoi(input[:4])
	if err != nil {
		return time.Time{}, errors.New("year parsing error")
	}

	month, err := strconv.Atoi(input[5:7])
	if err != nil {
		return time.Time{}, errors.New("month parsing error")
	}

	day, err := strconv.Atoi(input[8:10])
	if err != nil {
		return time.Time{}, errors.New("day parsing error")
	}

	hour, err := strconv.Atoi(input[11:13])
	if err != nil {
		return time.Time{}, errors.New("hour parsing error")
	}

	minute, err := strconv.Atoi(input[14:16])
	if err != nil {
		return time.Time{}, errors.New("minute parsing error")
	}

	second, err := strconv.Atoi(input[17:19])
	if err != nil {
		return time.Time{}, errors.New("second parsing error")
	}

	parsedTime := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)

	if parsedTime.Year() != year || parsedTime.Month() != time.Month(month) || parsedTime.Day() != day {
		return time.Time{}, errors.New("incorrect date")
	}

	return parsedTime, nil
}

func (app *application) parseTimeString(input string) (time.Time, error) {
	if len(input) != len("2006-01-02-15-04-05") {
		return time.Time{}, errors.New("incorrect time format, expect YYYY-MM-DD-HH-mm-ss")
	}

	// Specify the location for UTC+5 here:
	loc, err := time.LoadLocation("Asia/Aqtobe") // Adjust the time zone as needed
	if err != nil {
		return time.Time{}, errors.New("failed to load time zone")
	}

	layout := "2006-01-02-15-04-05" // Specific layout for your string format
	parsedTime, err := time.ParseInLocation(layout, input, loc)
	if err != nil {
		return time.Time{}, errors.New("time parsing error")
	}

	return parsedTime, nil
}

func (app *application) objNameFromURL(imageURL string) (string, error) {
	if imageURL == "" {
		objID, _ := uuid.NewRandom()
		return objID.String(), nil
	}

	urlPath, err := url.Parse(imageURL)

	if err != nil {
		app.logger.PrintError(fmt.Errorf("failed to parse objectName from imageURL: %v\n", imageURL), nil)
		return "", err
	}

	return path.Base(urlPath.Path), nil
}

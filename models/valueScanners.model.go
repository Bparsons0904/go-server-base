package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type JSONB map[string]interface{}

func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, a)
}

func CastJSON[T any](object T) (result JSONB) {
	data, _ := json.Marshal(object)
	json.Unmarshal(data, &result)
	return
}

type TSRANGE struct {
	Lower time.Time
	Upper *time.Time // Using pointer so it can be nil (unbounded)
}

// Implement Scanner and Valuer for custom type for proper DB interaction
func (r *TSRANGE) Scan(value interface{}) error {
	// Convert value from DB format to TSRANGE.
	return nil
}

func (r TSRANGE) Value() (driver.Value, error) {
	if r.Upper == nil {
		return fmt.Sprintf("[%s,)", r.Lower.Format(time.RFC3339)), nil
	}
	return fmt.Sprintf("[%s,%s)", r.Lower.Format(time.RFC3339), r.Upper.Format(time.RFC3339)), nil
}

type TSTZRANGE struct {
	Lower time.Time
	Upper *time.Time // Using pointer so it can be nil (unbounded)
}

const pgLayout = "2006-01-02 15:04:05-07"

func (r *TSTZRANGE) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected value to be a string")
	}

	str = strings.TrimSpace(str)
	str = strings.TrimLeft(str, "[")
	str = strings.TrimRight(str, ")")

	parts := strings.Split(str, ",")

	if str != "," {
		var err error
		if parts[0] != "" {
			r.Lower, err = time.Parse(pgLayout, strings.Trim(parts[0], "\""))
			if err != nil {
				return fmt.Errorf("error parsing lower bound: %w", err)
			}
		}
		if parts[1] != "" {
			r.Upper = new(time.Time)
			*r.Upper, err = time.Parse(pgLayout, strings.Trim(parts[1], "\""))
			if err != nil {
				return fmt.Errorf("error parsing upper bound: %w", err)
			}
		}
	}

	return nil
}

func (r TSTZRANGE) Value() (driver.Value, error) {
	if r.Upper == nil {
		return fmt.Sprintf("[%s,)", r.Lower.Format(pgLayout)), nil
	}
	return fmt.Sprintf("[%s,%s)", r.Lower.Format(pgLayout), r.Upper.Format(pgLayout)), nil
}

package cfg

import (
	"encoding"
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/frozzare/go/structs"
)

var (
	// ErrNotFound is the field not found error.
	ErrNotFound = errors.New("Field not found")
)

// Config represents a config struct.
type Config struct {
	target interface{}
}

// New creates a new configration.
func New(s interface{}, option ...Option) (*Config, error) {
	c := &Config{
		target: s,
	}

	if err := c.Extend(option...); err != nil {
		return nil, err
	}

	return c, nil
}

// Extend extends the existing configuration object with more values.
func (c *Config) Extend(option ...Option) error {
	for _, o := range option {
		if err := o(c); err != nil {
			return err
		}
	}

	return nil
}

var timeType = reflect.TypeOf(time.Time{})
var urlType = reflect.TypeOf(url.URL{})

func (c *Config) field(name string) (*structs.StructField, error) {
	p := strings.Split(name, ".")

	if len(p) == 0 {
		return nil, ErrNotFound
	}

	f, err := structs.Field(c.target, p[0])
	for _, r := range p[1:] {
		if f == nil {
			return nil, ErrNotFound
		}

		f, err = f.Field(r)

		if err != nil {
			return nil, err
		}
	}

	if f == nil {
		return nil, ErrNotFound
	}

	if f.Kind() == reflect.Ptr && f.IsZero() {
		v := f.ReflectValue()
		e := v.Type().Elem()

		if e.Kind() == reflect.Ptr {
			return nil, nil
		}

		z := reflect.New(e)
		v.Set(z)
	}

	return f, nil
}

func (c *Config) assign(field structs.StructField, val string) error {
	// Bail if field cannot be set.
	if !field.CanSet() {
		return errors.New("Field cannot be set")
	}

	v := reflect.ValueOf(val)
	f := field.ReflectValue()

	// Direct assignment.
	if v.Type().AssignableTo(f.Type()) {
		return field.Set(val)
	}

	t := f.Type()
	p := false
	if f.Kind() == reflect.Ptr {
		t = f.Elem().Type()
		p = true
	}

	// Assign special fields.
	switch t {
	case timeType:
		pt, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return err
		}

		if p {
			return field.Set(&pt)
		}

		return field.Set(pt)
	case urlType:
		u, err := url.Parse(val)
		if err != nil {
			return err
		}

		if p {
			return field.Set(u)
		}

		return field.Set(*u)
	}

	// Support TextUnmarshaler.
	iface := f.Addr().Interface()
	if tx, ok := iface.(encoding.TextUnmarshaler); ok {
		if err := tx.UnmarshalText([]byte(val)); err != nil {
			return err
		}

		return nil
	}

	// JSON Unmarshal fallback.
	return json.Unmarshal([]byte(val), iface)
}

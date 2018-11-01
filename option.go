package cfg

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/frozzare/go/env"
)

// Option configuration the configuration.
type Option func(*Config) error

// UnmarshalFunc is an unmarshal function.
type UnmarshalFunc func([]byte, interface{}) error

// WithContent unmarshal and sets the file with the configuration.
func WithContent(b []byte, f ...UnmarshalFunc) Option {
	return func(c *Config) error {
		if len(f) == 0 {
			f = []UnmarshalFunc{json.Unmarshal}
		}

		return f[0](b, c.target)
	}
}

// WithData sets a key value map with the configuration.
func WithData(v map[string]interface{}) Option {
	return func(c *Config) error {
		for name, val := range v {
			f, err := c.field(name)

			if f == nil && err == ErrNotFound {
				continue
			}

			if err != nil {
				return err
			}

			if f.Kind() != reflect.String && reflect.ValueOf(val).Kind() == reflect.String {
				if err := c.assign(*f, val.(string)); err != nil {
					return err
				}
			} else if err := f.Set(val); err != nil {
				return err
			}

		}

		return nil
	}
}

// WithFile reads, unmarshal and sets the file with the configuration.
func WithFile(p string, f ...UnmarshalFunc) Option {
	return func(c *Config) error {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		return WithContent(b, f...)(c)
	}
}

// WithEnvironment reads and sets environment values with the configuration.
func WithEnvironment(v map[string]string) Option {
	return func(c *Config) error {
		for name, envName := range v {
			f, err := c.field(name)

			if f == nil && err == ErrNotFound {
				continue
			}

			if err != nil {
				return err
			}

			envVal := env.Get(envName)

			if len(envVal) == 0 {
				continue
			}

			if err := c.assign(*f, envVal); err != nil {
				return err
			}
		}

		return nil
	}
}

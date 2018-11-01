package cfg

import (
	"net/url"
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
)

type CustomChild struct {
	Name string
}

type CustomConfig struct {
	Name     string
	Num      int
	URL      url.URL
	URL2     *url.URL
	Time     time.Time
	Time2    *time.Time
	Duration time.Duration
	Child    *CustomChild
	Child2   CustomChild
}

func TestWithContent(t *testing.T) {
	cc := &CustomConfig{}

	_, err := New(&cc, WithContent([]byte(`{"Name": "Fredrik", "Num": 1}`)))

	if err != nil {
		t.Fatal(err)
	}

	if cc.Name != "Fredrik" {
		t.Errorf("Unexpected name value: %q", cc.Name)
	}

	if cc.Num != 1 {
		t.Errorf("Unexpected num value: %q", cc.Num)
	}
}

func TestWithData(t *testing.T) {
	cc := &CustomConfig{}

	_, err := New(&cc, WithData(map[string]interface{}{
		"Name":           "Fredrik",
		"Num":            "1",
		"URL":            "https://www.golang.org",
		"URL2":           "https://www.golang.org",
		"Time":           "2012-11-01T22:08:41+00:00",
		"Time2":          "2012-11-01T22:08:41+00:00",
		"Duration":       56 * time.Second,
		"Child.Name":     "Fredrik",
		"Child2.Name":    "Fredrik2",
		"NotFound":       "NotFound",
		"NotFound.Child": "NotFound.Child",
	}))

	if err != nil {
		t.Fatal(err)
	}

	if cc.Name != "Fredrik" {
		t.Errorf("Unexpected name value: %q", cc.Name)
	}

	if cc.Num != 1 {
		t.Errorf("Unexpected num value: %q", cc.Num)
	}

	if cc.URL.String() != "https://www.golang.org" {
		t.Errorf("Unexpected url value: %q", cc.URL.String())
	}

	if cc.URL2.String() != "https://www.golang.org" {
		t.Errorf("Unexpected url value: %q", cc.URL2.String())
	}

	if cc.Time.String() != "2012-11-01 22:08:41 +0000 +0000" {
		t.Errorf("Unexpected time value: %q", cc.Time.String())
	}

	if cc.Time2.String() != "2012-11-01 22:08:41 +0000 +0000" {
		t.Errorf("Unexpected time value: %q", cc.Time2.String())
	}

	if cc.Duration.String() != "56s" {
		t.Errorf("Unexpected duration value: %q", cc.Duration)
	}

	if cc.Child.Name != "Fredrik" {
		t.Errorf("Unexpected name value: %q", cc.Child.Name)
	}

	if cc.Child2.Name != "Fredrik2" {
		t.Errorf("Unexpected name value: %q", cc.Child2.Name)
	}
}

func TestExtendConfig(t *testing.T) {
	cc := &CustomConfig{}

	cfg, err := New(&cc, WithData(map[string]interface{}{
		"Name": "Fredrik",
	}))

	if err != nil {
		t.Fatal(err)
	}

	if cc.Name != "Fredrik" {
		t.Errorf("Unexpected name value: %q", cc.Name)
	}

	err = cfg.Extend(WithData(map[string]interface{}{
		"URL": "https://www.golang.org",
	}))

	if err != nil {
		t.Fatal(err)
	}

	if cc.URL.String() != "https://www.golang.org" {
		t.Errorf("Unexpected url value: %q", cc.URL.String())
	}
}

func TestWithFileJson(t *testing.T) {
	cc := &CustomConfig{}

	_, err := New(&cc, WithFile("testdata/test.json"))

	if err != nil {
		t.Fatal(err)
	}

	if cc.Name != "Fredrik" {
		t.Errorf("Unexpected name value: %q", cc.Name)
	}

	if cc.Num != 1 {
		t.Errorf("Unexpected num value: %q", cc.Num)
	}
}

func TestWithFileYaml(t *testing.T) {
	cc := &CustomConfig{}

	_, err := New(&cc, WithFile("testdata/test.yml", yaml.Unmarshal))

	if err != nil {
		t.Fatal(err)
	}

	if cc.Name != "Fredrik" {
		t.Errorf("Unexpected name value: %q", cc.Name)
	}

	if cc.Num != 1 {
		t.Errorf("Unexpected num value: %q", cc.Num)
	}
}

func TestWithEnvironment(t *testing.T) {
	cc := &CustomConfig{}

	os.Setenv("TEST_NAME", "Fredrik")
	os.Setenv("TEST_NUM", "1")
	os.Setenv("TEST_URL", "https://www.golang.org")
	os.Setenv("TEST_TIME", "2012-11-01T22:08:41+00:00")
	os.Setenv("TEST_DURATION", "56")

	_, err := New(&cc, WithEnvironment(map[string]string{
		"Name":     "TEST_NAME",
		"Num":      "TEST_NUM",
		"URL":      "TEST_URL",
		"Time":     "TEST_TIME",
		"Duration": "TEST_DURATION",
	}))

	if err != nil {
		t.Fatal(err)
	}

	if cc.Name != "Fredrik" {
		t.Errorf("Unexpected name value: %q", cc.Name)
	}

	if cc.Num != 1 {
		t.Errorf("Unexpected num value: %q", cc.Num)
	}

	if cc.URL.String() != "https://www.golang.org" {
		t.Errorf("Unexpected url value: %q", cc.URL.String())
	}

	if cc.Time.String() != "2012-11-01 22:08:41 +0000 +0000" {
		t.Errorf("Unexpected time value: %q", cc.Time.String())
	}

	if cc.Duration.String() != "56ns" {
		t.Errorf("Unexpected duration value: %q", cc.Duration)
	}
}

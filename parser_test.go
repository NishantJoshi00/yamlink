package yamlink_test

import (
	"strings"
	"testing"

	"github.com/NishantJoshi00/yamlink"
)

func TestDepth1Parse(t *testing.T) {
	data := `
  key1: value1
  key2: value2
  `

	m, err := yamlink.ParseYaml(strings.NewReader(data))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	var nm map[string]interface{} = m.(map[string]interface{})

	if len(nm) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(nm))
	}

	if nm["key1"] != "value1" {
		t.Errorf("Expected value1, got %v", nm["key1"])
	}
}

func TestPathDepth1(t *testing.T) {
	table := `
  key1: example.com
  `

	m, err := yamlink.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	val, err := yamlink.PathLookup("key1", m)

	if err != nil {
		t.Errorf("Error looking up path: %v", err)
	}

	if val != "example.com" {
		t.Errorf("Expected example.com, got %v", val)
	}
}

func TestPathDepth2Map(t *testing.T) {
	table := `
  key1:
    key2: example.com
  `

	m, err := yamlink.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	val, err := yamlink.PathLookup("key1/key2", m)

	if err != nil {
		t.Errorf("Error looking up path: %v", err)
	}

	if val != "example.com" {
		t.Errorf("Expected example.com, got %v", val)
	}
}

func TestPathDepth2Array(t *testing.T) {
	table := `
  key1:
    - example.com
  `

	m, err := yamlink.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	val, err := yamlink.PathLookup("key1/0", m)

	if err != nil {
		t.Errorf("Error looking up path: %v", err)
	}

	if val != "example.com" {
		t.Errorf("Expected example.com, got %v", val)
	}
}

func TestPathDepth3ArrayMap(t *testing.T) {
	table := `
  key1:
    - magic.com
    - key2: example.com
  `

	m, err := yamlink.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	val, err := yamlink.PathLookup("key1/1/key2", m)

	if err != nil {
		t.Errorf("Error looking up path: %v", err)
	}

	if val != "example.com" {
		t.Errorf("Expected example.com, got %v", val)
	}
}

func TestUnknownPathDepth1(t *testing.T) {
	table := `
  key1:
    - magic.com
    - key2: example.com
  `

	m, err := yamlink.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	_, err = yamlink.PathLookup("key1/1/key3", m)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUnknownPathDepth2(t *testing.T) {
	table := `
  key1:
    - magic.com
    - key2: example.com
  `

	m, err := yamlink.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	_, err = yamlink.PathLookup("key1/2/key2", m)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

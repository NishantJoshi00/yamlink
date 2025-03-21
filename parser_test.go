package waypoint_test

import (
	"os"
	"strings"
	"testing"

	"github.com/NishantJoshi00/waypoint"
)

func TestDepth1Parse(t *testing.T) {
	data := `
  key1: value1
  key2: value2
  `

	m, err := waypoint.ParseYaml(strings.NewReader(data))

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

	m, err := waypoint.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	val, err := waypoint.PathLookup("key1", m)

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

	m, err := waypoint.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	val, err := waypoint.PathLookup("key1/key2", m)

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

	m, err := waypoint.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	val, err := waypoint.PathLookup("key1/0", m)

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

	m, err := waypoint.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	val, err := waypoint.PathLookup("key1/1/key2", m)

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

	m, err := waypoint.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	_, err = waypoint.PathLookup("key1/1/key3", m)

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

	m, err := waypoint.ParseYaml(strings.NewReader(table))

	if err != nil {
		t.Errorf("Error parsing yaml: %v", err)
	}

	_, err = waypoint.PathLookup("key1/2/key2", m)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseYamlError(t *testing.T) {
	// Test invalid YAML causing a decoder error
	invalidYaml := `
    key1: value1
    key2: [invalid
    `
	_, err := waypoint.ParseYaml(strings.NewReader(invalidYaml))
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}
func TestReadFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write some YAML to the file
	yamlContent := "key1: value1\nkey2: value2\n"
	if _, err := tmpfile.Write([]byte(yamlContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test reading the file
	m, err := waypoint.ReadFile(tmpfile.Name())
	if err != nil {
		t.Errorf("ReadFile returned error: %v", err)
	}

	// Verify contents
	mapResult, ok := m.(map[string]interface{})
	if !ok {
		t.Errorf("Expected map, got different type")
	}
	if mapResult["key1"] != "value1" {
		t.Errorf("Expected value1, got %v", mapResult["key1"])
	}

	// Test file not found
	_, err = waypoint.ReadFile("nonexistent-file.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestPathLookupWithTrailingSlash(t *testing.T) {
	yaml := "key1: value1"
	m, err := waypoint.ParseYaml(strings.NewReader(yaml))
	if err != nil {
		t.Fatal(err)
	}

	// Test with trailing slash
	val, err := waypoint.PathLookup("/key1/", m)

	if err != nil {
		t.Errorf("Error looking up path: %v", err)
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}
}

func TestPathLookupEmptyPath(t *testing.T) {
	yaml := "key1: value1"
	m, err := waypoint.ParseYaml(strings.NewReader(yaml))
	if err != nil {
		t.Fatal(err)
	}

	// Test with empty path
	_, err = waypoint.PathLookup("", m)
	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}
}

func TestPathLookupInvalidType(t *testing.T) {
	yaml := `
    array:
      - item1
      - item2
    map:
      key: value
    `
	m, err := waypoint.ParseYaml(strings.NewReader(yaml))
	if err != nil {
		t.Fatal(err)
	}

	// Test accessing array with string key
	_, err = waypoint.PathLookup("array/key", m)
	if err == nil {
		t.Error("Expected error when accessing array with string key, got nil")
	}

	// Test accessing map with numeric key
	_, err = waypoint.PathLookup("map/0", m)
	if err == nil {
		t.Error("Expected error when accessing map with numeric key, got nil")
	}
}

func TestReadFileOpenError(t *testing.T) {
	// Create a directory that can't be opened as a file
	dirName := "test_directory"
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirName)

	// Try to open the directory as a file (should fail)
	_, err = waypoint.ReadFile(dirName)
	if err == nil {
		t.Error("Expected error when opening directory as file, got nil")
	}
}

func TestIndexLookupInvalidIndex(t *testing.T) {
	yaml := `
    array:
      - item1
      - item2
    `
	m, err := waypoint.ParseYaml(strings.NewReader(yaml))
	if err != nil {
		t.Fatal(err)
	}

	// This is a bit of a trick - we're creating a scenario where an internal function
	// might be called with a non-numeric string that bypasses the normal checks
	// Note: In practice, this might not actually reach the uncovered code path
	// due to how the PathLookup function is implemented
	_, err = waypoint.PathLookup("array/invalid", m)
	if err == nil {
		t.Error("Expected error when accessing array with invalid index, got nil")
	}
}

package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListEnvs_Basic(t *testing.T) {
	data_dir := os.Getenv("DATA_DIR")

	ds := NewDataStore(data_dir)
	got, err := ds.listEnvs("appname-1")
	if err != nil {
		t.Fatalf("listEnvs returned error: %v", err)
	}

	want := []EnvNameType{"uat", "prod"}
	if !envSlicesEqualUnordered(got, want) {
		t.Fatalf("unexpected env list\ngot:  %v\nwant: %v", got, want)
	}
}

func TestListEnvs_CacheBehavior(t *testing.T) {
	data_dir := os.Getenv("DATA_DIR")

	ds := NewDataStore(data_dir)
	got, err := ds.listEnvs("appname-1")
	if err != nil {
		t.Fatalf("first listEnvs returned error: %v", err)
	}

	want := []EnvNameType{"uat", "prod"}
	if !envSlicesEqualUnordered(got, want) {
		t.Fatalf("unexpected first env list: %v", got)
	}

	// add another env file after the first call; since ds should cache,
	// second call should return same result, and not see the new file
	dummyFilePath := filepath.Join(data_dir, "appname-1", "dummy.yaml")
	if err := os.WriteFile(dummyFilePath, []byte("name: dummy"), 0o644); err != nil {
		t.Fatalf("write dummy.yaml: %v", err)
	}
	defer os.Remove(dummyFilePath)

	got, err = ds.listEnvs("appname-1")
	if err != nil {
		t.Fatalf("second listEnvs returned error: %v", err)
	}
	// expect cached result (only uat and prod and not dummy)
	if !envSlicesEqualUnordered(got, want) {
		t.Fatalf("cache did not hold. got: %v want: %v", got, want)
	}

	// a fresh datastore should see dummy as well
	ds2 := NewDataStore(data_dir)
	got2, err := ds2.listEnvs("appname-1")
	if err != nil {
		t.Fatalf("fresh datastore listEnvs error: %v", err)
	}
	want2 := []EnvNameType{"uat", "prod", "dummy"}
	if !envSlicesEqualUnordered(got2, want2) {
		t.Fatalf("fresh datastore unexpected env list\ngot:  %v\nwant: %v",
			got2, want2)
	}
}

func TestListEnvs_NonExistentApp(t *testing.T) {
	data_dir := os.Getenv("DATA_DIR")

	ds := NewDataStore(data_dir)
	_, err := ds.listEnvs("does-not-exist")
	if err == nil {
		t.Fatalf("expected error for non-existent app dir, got nil")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got: %v", err)
	}
}

// helper: compare two slices of EnvNameType ignoring order
func envSlicesEqualUnordered(a, b []EnvNameType) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int)
	for _, v := range a {
		m[string(v)]++
	}
	for _, v := range b {
		s := string(v)
		if m[s] == 0 {
			return false
		}
		m[s]--
	}
	for _, cnt := range m {
		if cnt != 0 {
			return false
		}
	}
	return true
}

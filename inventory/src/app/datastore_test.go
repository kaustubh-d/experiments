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

func TestLoadAppEnv_BasicAndCache(t *testing.T) {
	data_dir := os.Getenv("DATA_DIR")
	app := "appname-1"
	env := "prod"

	ds := NewDataStore(data_dir)

	// first load should succeed
	d1, err := ds.loadAppEnv(app, env)
	if err != nil {
		t.Fatalf("first loadAppEnv returned error: %v", err)
	}
	if d1 == nil {
		t.Fatalf("first loadAppEnv returned nil data")
	}

	// second load from same datastore should return cached pointer
	d2, err := ds.loadAppEnv(app, env)
	if err != nil {
		t.Fatalf("second loadAppEnv returned error: %v", err)
	}
	if d2 == nil {
		t.Fatalf("second loadAppEnv returned nil data")
	}
	if d1 != d2 {
		t.Fatalf("expected cached pointer on second load, got different pointers: %p vs %p", d1, d2)
	}
}

func TestLoadAppEnv_NonExistentApp(t *testing.T) {
	data_dir := os.Getenv("DATA_DIR")
	ds := NewDataStore(data_dir)
	_, err := ds.loadAppEnv("this-app-does-not-exist", "prod")
	if err == nil {
		t.Fatalf("expected error for non-existent env file, got nil")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got: %v", err)
	}
}

func TestLoadAppEnv_NonExistentEnv(t *testing.T) {
	data_dir := os.Getenv("DATA_DIR")
	ds := NewDataStore(data_dir)
	_, err := ds.loadAppEnv("appname-1", "this-env-does-not-exist")
	if err == nil {
		t.Fatalf("expected error for non-existent env file, got nil")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got: %v", err)
	}
}

func TestLoadEnabledApps_BasicAndCache(t *testing.T) {
	data_dir := os.Getenv("DATA_DIR")

	ds := NewDataStore(data_dir)

	apps1, err := ds.loadEnabledApps()
	if err != nil {
		t.Fatalf("first loadEnabledApps returned error: %v", err)
	}
	if apps1 == nil {
		t.Fatalf("first loadEnabledApps returned nil")
	}
	ua, ok := (*apps1)["appname-1"]
	if !ok {
		t.Fatalf("expected appname-1 in enabled apps, got keys: %v", apps1)
	}
	if len(ua.Users) != 1 {
		t.Fatalf("unexpected users for appname-1: %v", ua.Users)
	}

	// Another load should return cached pointer
	apps2, err := ds.loadEnabledApps()
	if err != nil {
		t.Fatalf("second loadEnabledApps returned error: %v", err)
	}
	if apps1 != apps2 {
		t.Fatalf("expected cached pointer on second load, got different pointers: %p vs %p", apps1, apps2)
	}

}

func TestLoadEnabledApps_NonExistent(t *testing.T) {
	data_dir := os.Getenv("DATA_DIR")

	ds := NewDataStore(filepath.Join(data_dir, "empty-data-store"))
	_, err := ds.loadEnabledApps()
	if err == nil {
		t.Fatalf("expected error for missing enabled-app-list.yaml, got nil")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got: %v", err)
	}
}

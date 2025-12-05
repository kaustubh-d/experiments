package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type AppNameType string
type EnvNameType string

type DataStore struct {
	dataDir string

	mu sync.RWMutex

	enabledAppsLoaded bool
	enabledApps       *EnabledApps

	// app -> env -> data
	appEnvCache map[AppNameType]map[EnvNameType]*AppEnvData

	// app -> list of env names
	envListCache map[AppNameType][]EnvNameType
}

func NewDataStore(dataDir string) *DataStore {
	return &DataStore{
		dataDir:      dataDir,
		appEnvCache:  make(map[AppNameType]map[EnvNameType]*AppEnvData),
		envListCache: make(map[AppNameType][]EnvNameType),
	}
}

// loadEnabledApps loads and returns the cached list of enabled applications.
// It is safe for concurrent use: a read lock is acquired first to return the
// already-cached value if present, otherwise a write lock is taken to load the
// data from disk and populate the cache (double-checked locking).
//
// The function reads the YAML file "enabled-app-list.yaml" from ds.dataDir,
// unmarshals it into an EnabledApps value, stores it in ds.enabledApps, sets
// ds.enabledAppsLoaded to true, and returns the value.
//
// Parameters:
//   - ds *DataStore: receiver providing the data directory, cache fields, and locks.
//
// Returns:
//   - EnabledApps: the parsed list of enabled applications (from cache or disk).
//   - error: non-nil if reading the file or unmarshalling YAML fails.
func (ds *DataStore) loadEnabledApps() (*EnabledApps, error) {
	ds.mu.RLock()
	if ds.enabledAppsLoaded {
		c := ds.enabledApps
		ds.mu.RUnlock()
		return c, nil
	}
	ds.mu.RUnlock()

	ds.mu.Lock()
	defer ds.mu.Unlock()
	// double-check
	if ds.enabledAppsLoaded {
		return ds.enabledApps, nil
	}
	path := filepath.Join(ds.dataDir, "enabled-app-list.yaml")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var apps EnabledApps
	if err := yaml.Unmarshal(b, &apps); err != nil {
		return nil, err
	}
	ds.enabledApps = &apps
	ds.enabledAppsLoaded = true
	return ds.enabledApps, nil
}

// loadAppEnv loads and returns the AppEnvData for the specified
// application and environment. Cache is used to avoid reloading from disk
// if already loaded. It is safe for concurrent use: a read lock is acquired
// first to check the cache, and if not found, a write lock is taken to load
// from disk and update the cache (double-checked locking).
func (ds *DataStore) loadAppEnv(app, env string) (*AppEnvData, error) {
	appName := AppNameType(app)
	envName := EnvNameType(env)

	// check cache first and return if found
	ds.mu.RLock()
	if envs, ok := ds.appEnvCache[appName]; ok {
		if data, ok2 := envs[envName]; ok2 {
			ds.mu.RUnlock()
			return data, nil
		}
	}
	ds.mu.RUnlock()

	// get lock to load from file and add to cache
	ds.mu.Lock()
	defer ds.mu.Unlock()
	// ensure map exists
	if _, ok := ds.appEnvCache[appName]; !ok {
		ds.appEnvCache[appName] = make(map[EnvNameType]*AppEnvData)
	}
	// double-check presence in cache
	if data, ok := ds.appEnvCache[appName][envName]; ok {
		return data, nil
	}

	// not in cache, load from file
	// data/appname/env.yaml
	path := filepath.Join(ds.dataDir, app, env+".yaml")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data AppEnvData
	if err := yaml.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	ds.appEnvCache[appName][envName] = &data
	return &data, nil
}

// listEnvs lists the available environment names for the specified application.
// It reads the filenames in the directory data/<app>/, extracts the names
// (without .yaml/.yml extensions), and returns them as a slice of EnvNameType.
// Cache is used to avoid re-reading the directory if already loaded.
func (ds *DataStore) listEnvs(app string) ([]EnvNameType, error) {
	appName := AppNameType(app)

	// check cache first and return if found
	log.Println("Listing envs for app:", appName)
	ds.mu.RLock()
	if list, ok := ds.envListCache[appName]; ok {
		log.Println("Found env list in cache for app:", appName)
		ds.mu.RUnlock()
		return list, nil
	}
	ds.mu.RUnlock()

	log.Println("Env list not in cache for app:", appName)
	// get exclusive lock to load from file and add to cache
	ds.mu.Lock()
	defer ds.mu.Unlock()
	log.Println("Acquired write lock to load env list for app:", appName)
	// double-check
	if list, ok := ds.envListCache[appName]; ok {
		log.Println("Found env list in cache for app (after double-check):", appName)
		return list, nil
	}

	// read directory contents, each file is an env
	dir := filepath.Join(ds.dataDir, app)
	log.Println("Reading env directory:", dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var envs []EnvNameType
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			env := EnvNameType(strings.TrimSuffix(name, filepath.Ext(name)))
			log.Println("Found env:", env)
			envs = append(envs, env)
		}
	}
	log.Println("Caching env list for app:", appName)
	ds.envListCache[appName] = envs
	return envs, nil
}

func (ds *DataStore) GetUserList(appName string) ([]string, error) {

	// Add access control list here.
	apps, err := ds.loadEnabledApps()
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("Cache load failed. app name: %s", appName)
	}

	users, ok := (*apps)[appName]
	if !ok {
		return nil, fmt.Errorf("%s not found", appName)
	}
	return users.Users, nil
}

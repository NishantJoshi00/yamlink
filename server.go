package waypoint

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

var Logger *slog.Logger = nil

type Config struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Map_file        string `yaml:"map_file"`
	RefreshInterval int    `yaml:"refresh_interval"`
}

type State struct {
	PathFinder   interface{}
	lastModified time.Time
}

func Init(config *Config) (*http.ServeMux, error) {

	state := State{}

	state.load(config.Map_file)

	Logger.Debug("State loaded successfully")

	go state.refresh(config.Map_file, config.RefreshInterval)

	mux := http.NewServeMux()
	mux.HandleFunc("/", state.magicHandler)
	mux.HandleFunc("/health", state.healthHandler)

	return mux, nil
}

func (state *State) healthHandler(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("Health check")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Health is good!"))
}

func (state *State) magicHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		state.healthHandler(w, r)
		return
	}

	Logger.Info("Request received", "query", r.URL.Path)

	resolve, err := PathLookup(r.URL.Path, state.PathFinder)

	if err != nil {
		Logger.Error("Failed to resolve path", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	Logger.Debug("Resolved path:", "answer", resolve)

	// redirect to the URL
	http.Redirect(w, r, resolve, http.StatusFound)
	return
}

func (state *State) refresh(map_file string, refresh_interval int) {
	Logger.Debug("Starting mapping refresh loop")

	for {
		Logger.Debug("Checking for file changes")
		time.Sleep(time.Duration(refresh_interval) * time.Second)

		changed, err := hasFileChanged(map_file, state.lastModified)
		if err != nil {
			Logger.Error("Failed to check file changes", "error", err)
			continue
		}

		if changed {
			Logger.Info("File changed, reloading")
			err := state.load(map_file)
			if err != nil {
				Logger.Error("Failed to load file", "error", err)
				continue
			}
		}
	}
}

func (state *State) load(path string) error {
	pathfinder, err := ReadFile(path)

	if err != nil {
		return err
	}
	state.PathFinder = pathfinder
	state.lastModified = time.Now()

	return nil
}

func hasFileChanged(path string, lastModified time.Time) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.ModTime().After(lastModified), nil
}

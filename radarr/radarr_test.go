package radarr

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
)

// Get info of the current directory of the executed file
var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func initTestEnvironnement(t *testing.T) {
	err := godotenv.Load(filepath.Join(basepath, "../.env"))
	if err != nil {
		t.Fatalf("Error while loading env file.\nErr: %s", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Error while setting test root path.\nErr: %s", err)
	}
}

func TestSendTmdbIDsToRadarr(t *testing.T) {
	initTestEnvironnement(t)

	tmdbIds := []string{"207", "4951", "995746", "9297"}
	SendTmdbIDsToRadarr(tmdbIds)
}

func TestGetRadarrState(t *testing.T) {
	initTestEnvironnement(t)

	fmt.Printf("%+v\n", GetRadarrState("9297"))
}

func TestAddToRadarrDownload(t *testing.T) {
	initTestEnvironnement(t)

	radarrState := RadarrStatus{
		HasFile:        true,
		Monitored:      true,
		Title:          "Monster House",
		TmdbId:         "9297",
		ProductionYear: 2006,
		IsAnimation:    true,
	}

	AddToRadarrDownload(radarrState)
}

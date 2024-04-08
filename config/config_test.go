package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		t.Fatalf("Failed to change directory")
	}

	conf := LoadConfiguration()
	fmt.Println(conf)
}

package radarr

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/mock"

	f "diikstra.fr/letterboxd-jellyfin-go/fetch"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) FetchData(fp f.FetcherParams) ([]byte, error) {
	args := m.Called(fp.Url)
	return args.Get(0).([]byte), args.Error(1)
}

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

func TestGetRadarrState(t *testing.T) {
	initTestEnvironnement(t)

	testData1 := []RadarrMovieLookupResp{{
		MovieFile: MovieLookupFile{
			RelativePath: "folder/to/alien",
		},
		Monitored: true,
		Title:     "Alien",
		TmdbId:    123,
		Year:      1979,
		Genres:    []string{"Horror", "Science Fiction"},
	}}
	byteTestData_1, _ := json.Marshal(testData1)

	testData2 := []RadarrMovieLookupResp{{
		MovieFile: MovieLookupFile{},
		Monitored: false,
		Title:     "Cloudy with a Chance of Meatballs",
		TmdbId:    22794,
		Year:      2009,
		Genres:    []string{"Animation", "Comedy", "Family"},
	}}
	byteTestData_2, _ := json.Marshal(testData2)

	type args struct {
		tmdbId string
	}
	tests := []struct {
		name           string
		args           args
		want           RadarrStatus
		clientResponse []byte
	}{
		{
			name: "Test Standard Movie State in Jellyfin",
			args: args{
				tmdbId: "123",
			},
			want: RadarrStatus{
				HasFile:        true,
				Monitored:      true,
				Title:          "Alien",
				TmdbId:         "123",
				ProductionYear: 1979,
				IsAnimation:    false,
			},
			clientResponse: byteTestData_1,
		},
		{
			name: "Test Animation Movie State not in Jellyfin",
			args: args{
				tmdbId: "22794",
			},
			want: RadarrStatus{
				HasFile:        false,
				Monitored:      false,
				Title:          "Cloudy with a Chance of Meatballs",
				TmdbId:         "22794",
				ProductionYear: 2009,
				IsAnimation:    true,
			},
			clientResponse: byteTestData_2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockClient)
			mockClient.On("FetchData", mock.Anything).Return(tt.clientResponse, nil)
			got, err := GetRadarrState(mockClient, tt.args.tmdbId)
			if err != nil {
				t.Errorf("GetRadarrState() returned error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRadarrState() = %v, want %v", got, tt.want)
			}
		})
	}
}

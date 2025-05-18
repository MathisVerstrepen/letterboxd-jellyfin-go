package jellyfin

import (
	"encoding/json"
	"path/filepath"
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
}

func TestGetUserId(t *testing.T) {
	initTestEnvironnement(t)

	byteTestData_1, _ := json.Marshal(
		[]User{{
			Name: "Jean Bon",
			Id:   "123",
		}, {
			Name: "Michel Sapin",
			Id:   "456",
		}},
	)

	type args struct {
		userName string
	}
	tests := []struct {
		name           string
		args           args
		want           string
		wantErr        bool
		clientResponse []byte
	}{
		{
			name: "Test existing user",
			args: args{
				userName: "Michel Sapin",
			},
			want:           "456",
			wantErr:        false,
			clientResponse: byteTestData_1,
		},
		{
			name: "Test unknown user",
			args: args{
				userName: "Wrong User",
			},
			want:           "",
			wantErr:        true,
			clientResponse: byteTestData_1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockClient)
			mockClient.On("FetchData", mock.Anything).Return(tt.clientResponse, nil)
			got, err := GetUserId(mockClient, tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserViews(t *testing.T) {
	initTestEnvironnement(t)

	testData1 := []UserView{{
		Name: "testMovie",
		Id:   "testId",
		UserData: UserData{
			Played: false,
		},
	}, {
		Name: "testMovie2",
		Id:   "testId2",
		UserData: UserData{
			Played: true,
		},
	}}
	byteTestData_1, _ := json.Marshal(ReqUserViewWrapper{
		Items: testData1})

	type args struct {
		userId           string
		userCollectionId string
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		want           []UserView
		clientResponse []byte
	}{
		{
			name: "Test with one played movie",
			args: args{
				userId:           "exampleUserId",
				userCollectionId: "exampleUserCollectionId",
			},
			wantErr:        false,
			want:           testData1,
			clientResponse: byteTestData_1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockClient)
			mockClient.On("FetchData", mock.Anything).Return(tt.clientResponse, nil)
			got, err := GetUserViews(mockClient, tt.args.userId, tt.args.userCollectionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserViews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for index, gotItem := range got {
				if gotItem != tt.want[index] {
					t.Errorf("GetUserViews() error = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_removeSeenMoviesFromUserCollection(t *testing.T) {
	mockClient := new(MockClient)
	initTestEnvironnement(t)

	byteTestData_1, _ := json.Marshal(ReqUserViewWrapper{
		Items: []UserView{{
			Name: "testMovie",
			Id:   "testId",
			UserData: UserData{
				Played: false,
			},
		}, {
			Name: "testMovie2",
			Id:   "testId2",
			UserData: UserData{
				Played: true,
			},
		}},
	})

	type args struct {
		userId           string
		userCollectionId string
	}
	tests := []struct {
		name           string
		args           args
		want           int
		clientResponse []byte
	}{
		{
			name: "Test with one played movie",
			args: args{
				userId:           "exampleUserId",
				userCollectionId: "exampleUserCollectionId",
			},
			want:           1,
			clientResponse: byteTestData_1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("FetchData", mock.Anything).Return(tt.clientResponse, nil)
			if got := RemoveSeenMoviesFromUserCollection(mockClient, tt.args.userId, tt.args.userCollectionId); got != tt.want {
				t.Errorf("removeSeenMoviesFromUserCollection() error = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func (m *MockClient) FetchData(fp f.FetcherParams) []byte {
	args := m.Called(fp.Url)
	return args.Get(0).([]byte)
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

func TestGetUsers(t *testing.T) {
	initTestEnvironnement(t)

	tests := []struct {
		name string
	}{
		{
			name: "test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUsers()
			if len(got) == 0 {
				t.Fatalf(`len(got) = 0, want > 0, error`)
			}
		})
	}
}

func TestGetUserId(t *testing.T) {
	initTestEnvironnement(t)

	type args struct {
		userName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test user test",
			args: args{
				userName: "test",
			},
			want:    "e64cd0aadacf4db3b05aca48aa8ef644",
			wantErr: false,
		},
		{
			name: "test unknown user",
			args: args{
				userName: "helloworld",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserId(tt.args.userName)
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
	mockClient := new(MockClient)
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

	byteTestData_2, _ := json.Marshal(struct {
		FalseField string
	}{
		FalseField: "testFalseJsonFormat",
	})

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
		{
			name: "Test with error from wrong json response",
			args: args{
				userId:           "exampleUserId",
				userCollectionId: "exampleUserCollectionId",
			},
			wantErr:        true,
			want:           testData1,
			clientResponse: byteTestData_2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			if got := removeSeenMoviesFromUserCollection(mockClient, tt.args.userId, tt.args.userCollectionId); got != tt.want {
				t.Errorf("removeSeenMoviesFromUserCollection() error = %v, want %v", got, tt.want)
			}
		})
	}
}

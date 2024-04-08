package jellyfin

import (
	"fmt"
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
	initTestEnvironnement(t)

	type args struct {
		userId           string
		userCollectionId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				userId:           "642e94dd14cb4045b31cf67a36ce998f",
				userCollectionId: "e88a35853226fa5bc6af39a1c29bfc09",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserViews(tt.args.userId, tt.args.userCollectionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserViews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}

func Test_removeSeenMoviesFromUserCollection(t *testing.T) {
	initTestEnvironnement(t)

	type args struct {
		userId           string
		userCollectionId string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test1",
			args: args{
				userId:           "642e94dd14cb4045b31cf67a36ce998f",
				userCollectionId: "e88a35853226fa5bc6af39a1c29bfc09",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeSeenMoviesFromUserCollection(tt.args.userId, tt.args.userCollectionId); got != tt.want {
				t.Errorf("removeSeenMoviesFromUserCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

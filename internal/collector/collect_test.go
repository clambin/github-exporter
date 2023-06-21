package collector

import (
	"github.com/clambin/github-exporter/internal/github"
	"testing"
)

func Test_analyzeRepo(t *testing.T) {
	type args struct {
		repo github.Repo
	}
	tests := []struct {
		name         string
		args         args
		wantUser     string
		wantArchived string
		wantFork     string
		wantPrivate  string
	}{
		{
			name: "base",
			args: args{
				repo: github.Repo{
					Name:     "bar",
					FullName: "foo/bar",
					Archived: true,
					Private:  true,
					Fork:     true,
				},
			},
			wantUser:     "foo",
			wantArchived: "true",
			wantPrivate:  "true",
			wantFork:     "true",
		},
		{
			name: "invalid",
			args: args{
				repo: github.Repo{
					Name:     "bar",
					FullName: "bar",
					Archived: true,
					Private:  true,
					Fork:     false,
				},
			},
			wantUser:     "unknown",
			wantArchived: "true",
			wantPrivate:  "true",
			wantFork:     "false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, gotArchived, gotFork, gotPrivate := analyzeRepo(tt.args.repo)
			if gotUser != tt.wantUser {
				t.Errorf("analyzeRepo() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
			if gotArchived != tt.wantArchived {
				t.Errorf("analyzeRepo() gotArchived = %v, want %v", gotArchived, tt.wantArchived)
			}
			if gotFork != tt.wantFork {
				t.Errorf("analyzeRepo() gotFork = %v, want %v", gotFork, tt.wantFork)
			}
			if gotPrivate != tt.wantPrivate {
				t.Errorf("analyzeRepo() gotPrivate = %v, want %v", gotPrivate, tt.wantPrivate)
			}
		})
	}
}

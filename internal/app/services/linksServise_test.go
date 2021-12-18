package services

import (
	"github.com/malyg1n/shortener/internal/app/storage"
	"testing"
)

func TestDefaultLinksService_GetLink(t *testing.T) {
	type fields struct {
		storage storage.LinksStorage
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DefaultLinksService{
				storage: tt.fields.storage,
			}
			got, err := s.GetLink(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultLinksService_SetLink(t *testing.T) {
	type fields struct {
		storage storage.LinksStorage
	}
	type args struct {
		link string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DefaultLinksService{
				storage: tt.fields.storage,
			}
			got, err := s.SetLink(tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_randomString(t *testing.T) {
	type args struct {
		n uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randomString(tt.args.n); got != tt.want {
				t.Errorf("randomString() = %v, want %v", got, tt.want)
			}
		})
	}
}

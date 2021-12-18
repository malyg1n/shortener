package storage

import (
	"reflect"
	"testing"
)

func TestLinksStorageMap_GetLink(t *testing.T) {
	type fields struct {
		links map[string]string
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
			s := &LinksStorageMap{
				links: tt.fields.links,
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

func TestLinksStorageMap_SetLink(t *testing.T) {
	type fields struct {
		links map[string]string
	}
	type args struct {
		id   string
		link string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LinksStorageMap{
				links: tt.fields.links,
			}
		})
	}
}

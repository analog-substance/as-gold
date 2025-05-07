package types

import (
	"reflect"
	"testing"
)

func TestHuman_AddURL(t *testing.T) {

	var h = NewHuman()
	h.AddURL("https://example.com/username")

	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Username already exists",
			args: args{url: "https://example.com/username"},
			want: []string{"https://example.com/username"},
		},
		{
			name: "Username doesn't exist",
			args: args{url: "https://example2.com/username"},
			want: []string{"https://example.com/username", "https://example2.com/username"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h.AddURL(tt.args.url)

			if got := h.URLs; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddURL() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestHuman_AddUsername(t *testing.T) {

	var h = NewHuman()
	h.AddUsername("nobody")

	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Username already exists",
			args: args{username: "nobody"},
			want: []string{"nobody"},
		},
		{
			name: "Username doesn't exist",
			args: args{username: "nobodycares"},
			want: []string{"nobody", "nobodycares"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h.AddUsername(tt.args.username)

			if got := h.Usernames; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddUsername() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestHuman_AddEmail(t *testing.T) {

	var h = NewHuman()
	h.AddEmail("test@localhost")

	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Email already exists",
			args: args{email: "test@localhost"},
			want: []string{"test@localhost"},
		},
		{
			name: "Email doesn't exist",
			args: args{email: "test2@localhost"},
			want: []string{"test2@localhost", "test@localhost"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h.AddEmail(tt.args.email)

			if got := h.Emails; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddEmail() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_addToSlice(t *testing.T) {
	type args struct {
		slice []string
		item  string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Does not add duplicates",
			args: args{
				slice: []string{"test@localhost"},
				item:  "test@localhost",
			},
			want: []string{"test@localhost"},
		},
		{
			name: "Does not add empty",
			args: args{
				slice: []string{"test@localhost"},
				item:  "   ",
			},
			want: []string{"test@localhost"},
		},
		{
			name: "Adds unique",
			args: args{
				slice: []string{"test2@localhost"},
				item:  "test@localhost",
			},
			want: []string{"test2@localhost", "test@localhost"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addToSlice(tt.args.slice, tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

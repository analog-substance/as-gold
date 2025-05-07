package types

import (
	"reflect"
	"testing"
)

func TestGroup_FindOrCreateHumanByUsername(t *testing.T) {

	var bob = NewHuman()
	bob.AddUsername("bob")

	var bobby = NewHuman()
	bobby.AddUsername("bobby")

	type fields struct {
		Members []*Human
	}
	type args struct {
		username string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Human
	}{
		{
			name: "find human by username",
			fields: fields{
				Members: []*Human{bob},
			},
			args: args{
				username: "bob",
			},
			want: bob,
		},

		{
			name: "find human by username",
			fields: fields{
				[]*Human{bob, bobby},
			},
			args: args{
				username: "bobby",
			},
			want: bobby,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				Members: tt.fields.Members,
			}
			if got := g.FindOrCreateHumanByUsername(tt.args.username); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOrCreateHumanByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

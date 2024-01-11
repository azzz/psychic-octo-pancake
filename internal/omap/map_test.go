package omap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOMap_Set(t *testing.T) {
	m := newMap()
	assert.False(t, m.Set("alice", 33))
	assert.False(t, m.Set("bob", 42))
	assert.True(t, m.Set("alice", 30))

	assert.Equal(t, []string{"alice", "bob"}, m.Keys())
}

func TestOMap_Get(t *testing.T) {
	nonEmpty := newMap()
	nonEmpty.Set("alice", 42)
	nonEmpty.Set("bob", 33)

	type args struct {
		key string
	}

	type testCase struct {
		name  string
		m     *OMap[string, int]
		args  args
		want  int
		want1 bool
	}
	tests := []testCase{
		{
			"empty map",
			newMap(),
			args{"alice"},
			0,
			false,
		},

		{
			"key exist",
			nonEmpty,
			args{"alice"},
			42,
			true,
		},

		{
			"key does not exist",
			nonEmpty,
			args{"john"},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.m.Get(tt.args.key)
			assert.Equalf(t, tt.want, got, "Get(%v)", tt.args.key)
			assert.Equalf(t, tt.want1, got1, "Get(%v)", tt.args.key)
		})
	}
}

func TestOMap_Nth(t *testing.T) {
	nonEmpty := newMap()
	nonEmpty.Set("alice", 42)
	nonEmpty.Set("bob", 33)

	type args struct {
		n int
	}
	type testCase struct {
		name  string
		m     *OMap[string, int]
		args  args
		want  int
		want1 bool
	}
	tests := []testCase{
		{
			"empty map",
			newMap(),
			args{3},
			0,
			false,
		},

		{
			"key exists",
			nonEmpty,
			args{1},
			33, // bob
			true,
		},

		{
			"key exists",
			nonEmpty,
			args{10},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.m.Nth(tt.args.n)
			assert.Equalf(t, tt.want, got, "Nth(%v)", tt.args.n)
			assert.Equalf(t, tt.want1, got1, "Nth(%v)", tt.args.n)
		})
	}
}

func TestOMap_Keys(t *testing.T) {
	nonEmpty := newMap()
	nonEmpty.Set("alice", 42)
	nonEmpty.Set("bob", 33)

	type testCase struct {
		name string
		m    *OMap[string, int]
		want []string
	}
	tests := []testCase{
		{
			"empty map",
			newMap(),
			nil,
		},

		{
			"non empty map",
			nonEmpty,
			[]string{"alice", "bob"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.m.Keys(), "Keys()")
		})
	}
}

func newMap() *OMap[string, int] {
	return New[string, int]()
}

func TestOMap_Pairs(t *testing.T) {
	nonEmpty := newMap()
	nonEmpty.Set("alice", 42)
	nonEmpty.Set("bob", 33)

	type testCase struct {
		name string
		m    *OMap[string, int]
		want []Pair[string, int]
	}
	tests := []testCase{
		{
			"empty map",
			newMap(),
			[]Pair[string, int]{},
		},

		{
			"non empty map",
			nonEmpty,
			[]Pair[string, int]{
				Pair[string, int]{"alice", 42},
				Pair[string, int]{"bob", 33},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.m.Pairs(), "Pairs()")
		})
	}
}

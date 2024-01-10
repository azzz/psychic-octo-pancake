package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_Push(t *testing.T) {
	l := New[int]()
	l.Push(4)
	l.Push(8)
	l.Push(15)

	assert.Equal(t, []int{4, 8, 15}, l.Values())
}

func TestList_Unshift(t *testing.T) {
	l := New[int]()
	l.Unshift(4)
	l.Unshift(8)
	l.Unshift(15)

	assert.Equal(t, []int{15, 8, 4}, l.Values())
}

func TestList_Nth(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := New[int]()
		assert.Nil(t, l.Nth(0))
		assert.Nil(t, l.Nth(33))
	})

	t.Run("n is out of the list boundaries", func(t *testing.T) {
		l := New[int](4)
		assert.Nil(t, l.Nth(33))
		assert.Nil(t, l.Nth(-33))
	})

	t.Run("valid case", func(t *testing.T) {
		l := New[int](4, 8, 15)

		assert.Equal(t, 4, l.Nth(0).Value())
		assert.Equal(t, 8, l.Nth(1).Value())
		assert.Equal(t, 15, l.Nth(2).Value())
	})

}

func TestList_Values(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := New[int]()

		assert.Zero(t, l.Values())
	})

	t.Run("single element", func(t *testing.T) {
		l := New[int](4)

		expected := []int{4}
		assert.Equal(t, expected, l.Values())
	})

	t.Run("many elenments", func(t *testing.T) {
		l := New[int](4, 8, 15, 16, 23, 42)

		expected := []int{4, 8, 15, 16, 23, 42}
		assert.Equal(t, expected, l.Values())
	})
}

func TestList_Delete(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := New[int]()
		assert.False(t, l.Delete(42))
		assert.Zero(t, l.Values())
	})

	t.Run("n is bigger then list size", func(t *testing.T) {
		l := New[int](4, 8, 15, 16, 23, 42)
		assert.False(t, l.Delete(33))
		assert.Equal(t, []int{4, 8, 15, 16, 23, 42}, l.Values())
	})

	t.Run("n is negative", func(t *testing.T) {
		l := New[int](4, 8, 15, 16, 23, 42)
		assert.False(t, l.Delete(-33))
		assert.Equal(t, []int{4, 8, 15, 16, 23, 42}, l.Values())
	})

	t.Run("delete first", func(t *testing.T) {
		l := New[int](4, 8, 15, 16, 23, 42)
		assert.True(t, l.Delete(0))
		assert.Equal(t, []int{8, 15, 16, 23, 42}, l.Values())
	})

	t.Run("delete last", func(t *testing.T) {
		l := New[int](4, 8, 15, 16, 23, 42)
		assert.True(t, l.Delete(5))
		assert.Equal(t, []int{4, 8, 15, 16, 23}, l.Values())
	})

	t.Run("delete in the middle", func(t *testing.T) {
		l := New[int](4, 8, 15, 16, 23, 42)
		assert.True(t, l.Delete(2))
		assert.Equal(t, []int{4, 8, 16, 23, 42}, l.Values())
	})
}

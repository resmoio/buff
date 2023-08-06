package buff

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type collectedItems[T any] struct {
	pushes *[][]T
	when   *[]time.Time
	fn     func(items []T)
}

func (c *collectedItems[T]) Equals(t *testing.T, items [][]T) {
	assert.Equal(t, items, *c.pushes)
}

func getFn[T any]() *collectedItems[T] {
	pushes := make([][]T, 0)
	when := make([]time.Time, 0)
	return &collectedItems[T]{
		pushes: &pushes,
		when:   &when,
		fn: func(items []T) {
			t := time.Now()
			log.Println("Flush", items)
			pushes = append(pushes, items)
			when = append(when, t)
		},
	}
}

func TestBuffer_Init(t *testing.T) {
	fn := func(items []string) {
		log.Println("Flush", items)
	}

	b := NewBuffer(fn, 5, time.Millisecond*100)
	b.Start()
	defer b.Close()
}

func TestBuffer_FewItems1(t *testing.T) {
	str := getFn[string]()

	b := NewBuffer(str.fn, 5, time.Millisecond*100)
	b.Start()

	b.Push("1")
	b.Push("2")
	b.Push("3")
	b.Push("4")
	time.Sleep(time.Millisecond * 150)
	b.Push("5")
	time.Sleep(time.Millisecond * 150)
	b.Push("6")
	b.Close()

	var expected = [][]string{
		{"1", "2", "3", "4"},
		{"5"},
		{"6"},
	}
	str.Equals(t, expected)
}

func TestBuffer_FewItems2(t *testing.T) {
	str := getFn[string]()

	b := NewBuffer(str.fn, 2, time.Millisecond*100)
	b.Start()

	b.Push("1")
	b.Push("2")
	b.Push("3")
	b.Push("4")
	time.Sleep(time.Millisecond * 150)
	b.Push("5")
	time.Sleep(time.Millisecond * 100)
	b.Push("6")
	b.Close()

	var expected = [][]string{
		{"1", "2"},
		{"3", "4"},
		{"5"},
		{"6"},
	}
	str.Equals(t, expected)
}

func TestBuffer_FewItems3(t *testing.T) {
	str := getFn[string]()

	b := NewBuffer(str.fn, 10, time.Millisecond*100)
	b.Start()

	b.Push("1")
	time.Sleep(time.Millisecond * 150)
	b.Push("2")
	b.Push("3")
	b.Push("4")
	time.Sleep(time.Millisecond * 150)
	b.Push("5")
	b.Push("6")
	b.Push("7")
	time.Sleep(time.Millisecond * 150)
	b.Push("8")
	time.Sleep(time.Millisecond * 150)
	b.Push("9")
	b.Close()

	var expected = [][]string{
		{"1"},
		{"2", "3", "4"},
		{"5", "6", "7"},
		{"8"},
		{"9"},
	}
	str.Equals(t, expected)
}

func TestBuffer_FewItems5(t *testing.T) {
	str := getFn[string]()

	b := NewBuffer(str.fn, 2, time.Millisecond*100)
	b.Start()

	b.Push("1")
	time.Sleep(time.Millisecond * 150)
	b.Push("2")
	b.Push("3")
	b.Push("4")
	time.Sleep(time.Millisecond * 150)
	b.Push("5")
	b.Push("6")
	b.Push("7")
	time.Sleep(time.Millisecond * 150)
	b.Push("8")
	time.Sleep(time.Millisecond * 150)
	b.Push("9")
	b.Close()

	var expected = [][]string{
		{"1"},
		{"2", "3"},
		{"4"},
		{"5", "6"},
		{"7"},
		{"8"},
		{"9"},
	}
	str.Equals(t, expected)
}

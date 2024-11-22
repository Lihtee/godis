package storage

import (
	"slices"
	"strconv"
	"testing"
	"time"
)

var (
	now       = time.Now()
	nowPlus10 = now.Add(time.Second * 10)
	nowPlus20 = now.Add(time.Second * 20)
)

func TestSortTtl(t *testing.T) {
	cases := []struct {
		input    []*time.Time
		expected []*time.Time
	}{
		{[]*time.Time{}, []*time.Time{}},
		{[]*time.Time{nil, nil, nil}, []*time.Time{nil, nil, nil}},
		{[]*time.Time{&now, &nowPlus20, &nowPlus10}, []*time.Time{&now, &nowPlus10, &nowPlus20}},
		{[]*time.Time{&now, &nowPlus20, &nowPlus10, nil, nil, nil}, []*time.Time{nil, nil, nil, &now, &nowPlus10, &nowPlus20}},
	}

	for _, testCase := range cases {
		storage := &GodisStorage{}
		storage.ttl = make([]GodisTtl, len(testCase.input))
		for i, ttl := range testCase.input {
			storage.ttl[i].ttl = ttl
			storage.ttl[i].key = string(rune(i))
		}

		storage.sortTtl()
		for i, ttl := range storage.ttl {
			if ttl.ttl != testCase.expected[i] {
				t.Errorf("ttl.ttl = %v, expected: %v", ttl.ttl, testCase.expected[i])
			}
		}
	}
}

func TestRemoveStaleKeys(t *testing.T) {
	inPast := now.Add(-time.Second)
	inFuture := now.Add(time.Hour)
	cases := []struct {
		input    []*time.Time
		expected []int
	}{
		{[]*time.Time{}, []int{}},
		{[]*time.Time{&inPast}, []int{}},
		{[]*time.Time{&inFuture}, []int{0}},
		{[]*time.Time{nil, nil, nil}, []int{0, 1, 2}},
		{[]*time.Time{&inFuture, &inFuture, &inPast, &inPast, nil}, []int{0, 1, 4}},
	}

	for _, testCase := range cases {
		storage := New(true)
		for i, ttl := range testCase.input {
			key := strconv.Itoa(i)
			value := GodisValue{
				stringValue: key,
				ttl:         ttl,
			}
			storage.data[key] = &value
			storage.ttl = append(storage.ttl, GodisTtl{key: key, ttl: value.ttl})
		}
		storage.removeStaleKeys()

		for _, i := range testCase.expected {
			key := strconv.Itoa(i)
			if _, ok := storage.data[key]; !ok {
				t.Errorf("key %s should not be deleted", key)
			}
		}

		for key := range storage.data {
			intKey, err := strconv.Atoi(key)
			if err != nil {
				t.Errorf("test is broken: key %s should be a number", key)
				continue
			}
			if !slices.Contains(testCase.expected, intKey) {
				t.Errorf("key %s was not removed", key)
			}
		}
	}
}

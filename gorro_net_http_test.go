package gorro

import (
	"reflect"
	"testing"
)

func Test_getFlattened(t *testing.T) {
	var regex = []rune(`(/blog/(\d+)/post/(\d+))|(/blog/(\d+))|(other)`)
	positions, i, r := getFlattened(regex, 0, 0)
	if i != 23 {
		t.Fatalf(`iteration pos: Expected %d and got %d`, 23, i)
	}
	if r != 2 {
		t.Fatalf(`iteration pos: Expected %d and got %d`, 2, r)
	}
	expected := []int{1, 2}
	if !reflect.DeepEqual(positions, expected) {
		t.Fatalf(`iteration pos: Expected %v and got %v`, expected, positions)
	}
}

func Test_getFlattenedIgnoresNonCapturingGroup(t *testing.T) {
	var regex = []rune(`(/blog/(\d+)/(?:post)/(\d+))|(/blog/(\d+))|(other)`)
	positions, i, r := getFlattened(regex, 0, 0)
	if i != 27 {
		t.Fatalf(`iteration pos: Expected %d and got %d`, 27, i)
	}
	if r != 2 {
		t.Fatalf(`iteration pos: Expected %d and got %d`, 2, r)
	}
	expected := []int{1, 2}
	if !reflect.DeepEqual(positions, expected) {
		t.Fatalf(`iteration pos: Expected %v and got %v`, expected, positions)
	}
}

func Test_getFlattenedReturnsEmptyWhenNoGroups(t *testing.T) {
	var regex = []rune(`(/blog/\d+/post/\d+)|(/blog/(\d+))|(other)`)
	positions, i, r := getFlattened(regex, 0, 0)
	if i != 19 {
		t.Fatalf(`iteration pos: Expected %d and got %d`, 19, i)
	}
	if r != 0 {
		t.Fatalf(`iteration pos: Expected %d and got %d`, 0, r)
	}
	expected := []int{}
	if !reflect.DeepEqual(positions, expected) {
		t.Fatalf(`iteration pos: Expected %v and got %v`, expected, positions)
	}
}

func Test_getPatterns(t *testing.T) {
	var regex = []rune(`(/blog/(\d+)/post/(\d+))|(/blog/(\d+))|(other)`)
	patterns := getPatterns(regex)

	expected := make(map[int]queryRegex)
	expected[1] = queryRegex{0, []int{2, 3}}
	expected[4] = queryRegex{1, []int{5}}
	expected[6] = queryRegex{2, []int{}}
	if !reflect.DeepEqual(patterns, expected) {
		t.Fatalf(`iteration pos: Expected %v and got %v`, expected, patterns)
	}
}

func Test_getPatternsMultiple(t *testing.T) {
	var regex = []rune(`(/blog/(\d+)/post/(\d+))|(/blog/(\d+))|(o(t)h(e)(?:r))`)
	patterns := getPatterns(regex)

	expected := make(map[int]queryRegex)
	expected[1] = queryRegex{0, []int{2, 3}}
	expected[4] = queryRegex{1, []int{5}}
	expected[6] = queryRegex{2, []int{7, 8}}
	if !reflect.DeepEqual(patterns, expected) {
		t.Fatalf(`iteration pos: Expected %v and got %v`, expected, patterns)
	}
}

func Test_getPatternsComplex(t *testing.T) {
	var regex = []rune(`(/test/(\d+)/p/(?P<pArg>[^/]+)/\d+/(?P<rest>.*)) | (/blog/(\d+)/post/(?P<post>\d+)(?P<rest>.*)) | (.*)`)
	patterns := getPatterns(regex)

	expected := make(map[int]queryRegex)
	expected[1] = queryRegex{0, []int{2, 3, 4}}
	expected[5] = queryRegex{1, []int{6, 7, 8}}
	expected[9] = queryRegex{2, []int{}}
	if !reflect.DeepEqual(patterns, expected) {
		t.Fatalf(`iteration pos: Expected %v and got %v`, expected, patterns)
	}
}

package main

import (
	"reflect"
	"testing"
)

func checkEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestStatTracker_NilGet(t *testing.T) {
	// Test that get() works on a nil StatTracker
	var st *StatTracker
	expected := &StatResponse{}
	checkEqual(t, expected, st.get())
}

func TestStatTracker_AddGet(t *testing.T) {
	st := NewStatTracker()

	stat1 := &Stat{
		Tags: []string{"foo:bar", "baz:blarg"},
	}

	stat2 := &Stat{
		Tags: []string{"foo:bar", "baz:bang"},
	}

	// Add a single stat
	st.add(stat1)

	// Check the output
	expected := &StatResponse{
		Count: 1,
		Tags: map[string][]string{
			"foo": {"bar"},
			"baz": {"blarg"},
		},
		TagSets: map[string]int{
			"baz:blarg,foo:bar": 1,
		},
	}
	checkEqual(t, expected, st.get())

	// Add another stat, plus the first one again
	st.add(stat2)
	st.add(stat1)

	// Check the output
	expected = &StatResponse{
		Count: 3,
		Tags: map[string][]string{
			"foo": {"bar"},
			"baz": {"bang", "blarg"},
		},
		TagSets: map[string]int{
			"baz:blarg,foo:bar": 2,
			"baz:bang,foo:bar":  1,
		},
	}
	checkEqual(t, expected, st.get())
}

func TestSummary_Empty(t *testing.T) {
	s := NewSummary()

	// Try to get a metric that doesn't exist
	expected1 := &StatResponse{}
	r := s.get("non.existent.metric")
	checkEqual(t, expected1, r)

	expected2 := map[string]int{}
	count := s.getAllCount()
	checkEqual(t, expected2, count)

	expected3 := map[string]*StatResponse{}
	details := s.getAllDetails()
	checkEqual(t, expected3, details)
}

func TestSummary_WithMetrics(t *testing.T) {
	s := NewSummary()

	stat1 := &Stat{
		Name: "first.metric",
		Tags: []string{"foo:bar", "baz:blarg"},
	}

	stat2 := &Stat{
		Name: "second.metric",
		Tags: []string{"foo:bar", "baz:bang"},
	}

	// Add stat1 once, stat2 twice
	s.add(stat1)
	s.add(stat2)
	s.add(stat2)

	// Get one of the metrics
	expectedFirst := &StatResponse{
		Count: 1,
		Tags: map[string][]string{
			"foo": {"bar"},
			"baz": {"blarg"},
		},
		TagSets: map[string]int{
			"baz:blarg,foo:bar": 1,
		},
	}
	r := s.get(stat1.Name)
	checkEqual(t, expectedFirst, r)

	// Get a count of all metrics
	expected2 := map[string]int{
		"first.metric":  1,
		"second.metric": 2,
	}
	count := s.getAllCount()
	checkEqual(t, expected2, count)

	// Get details about all metrics
	expectedSecond := &StatResponse{
		Count: 2,
		Tags: map[string][]string{
			"foo": {"bar"},
			"baz": {"bang"},
		},
		TagSets: map[string]int{
			"baz:bang,foo:bar": 2,
		},
	}
	details := s.getAllDetails()
	checkEqual(t, 2, len(details))
	checkEqual(t, expectedFirst, details["first.metric"])
	checkEqual(t, expectedSecond, details["second.metric"])
}

func TestSummary_Reset(t *testing.T) {
	s := NewSummary()

	stat1 := &Stat{
		Name: "first.metric",
		Tags: []string{"foo:bar", "baz:blarg"},
	}

	stat2 := &Stat{
		Name: "second.metric",
		Tags: []string{"foo:bar", "baz:bang"},
	}

	// Add a few metrics
	s.add(stat1)
	s.add(stat2)

	// Check that it worked
	expected1 := map[string]int{
		"first.metric":  1,
		"second.metric": 1,
	}
	count1 := s.getAllCount()
	checkEqual(t, expected1, count1)

	// Reset the summary
	s.reset()

	// Make sure the count was reset
	expected2 := map[string]int{}
	count2 := s.getAllCount()
	checkEqual(t, expected2, count2)

	// Details should also be reset
	expected3 := map[string]*StatResponse{}
	details := s.getAllDetails()
	checkEqual(t, expected3, details)
}

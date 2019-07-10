package main

import (
	"sort"
	"strings"
	"sync"
)

var empty struct{}

// StatTracker holds a summary for a single metric. Note that it does NOT handle RW locking.
type StatTracker struct {
	count int
	// Map of tag name to tag value (to empty struct, such that it is effectively a set of values)
	tags map[string]map[string]struct{}
	// Map of each unique combination of tags to the number of times that group appeared.
	// This is useful for estimating cardinality, as each set of tags corresponds to a unique timeseries
	tagSets map[string]int
}

// StatResponse is used to json-encode a StatTracker
type StatResponse struct {
	Count   int                 `json:"count"`
	Tags    map[string][]string `json:"tags"`
	TagSets map[string]int      `json:"tag_sets"`
}

// NewStatTracker initializes a StatTracker
func NewStatTracker() *StatTracker {
	return &StatTracker{
		tags:    make(map[string]map[string]struct{}),
		tagSets: make(map[string]int),
	}
}

func (t *StatTracker) add(stat *Stat) {
	t.count++

	// Sort tags before joining them
	sort.Strings(stat.Tags)
	t.tagSets[strings.Join(stat.Tags, ",")]++

	for _, tag := range stat.Tags {
		parts := strings.Split(tag, ":")
		if len(parts) != 2 {
			log.Warn("failed parsing a tag: unexpected format", "tag", tag)
			continue
		}
		if _, ok := t.tags[parts[0]]; !ok {
			t.tags[parts[0]] = make(map[string]struct{})
		}
		t.tags[parts[0]][parts[1]] = empty
	}
}

func (t *StatTracker) get() *StatResponse {
	r := &StatResponse{}

	if t == nil {
		return r
	}

	r.Count = t.count
	r.Tags = make(map[string][]string, len(t.tags))
	for tag, vals := range t.tags {
		r.Tags[tag] = make([]string, len(vals))
		i := 0
		for v := range vals {
			r.Tags[tag][i] = v
			i++
		}
		sort.Strings(r.Tags[tag])
	}

	r.TagSets = t.tagSets
	return r
}

// ------------------------------------------------------------------------

// Summary manages locking for a map of metric names to StatTrackers
type Summary struct {
	metrics map[string]*StatTracker

	sync.RWMutex
}

// NewSummary initializes a Summary
func NewSummary() *Summary {
	return &Summary{
		metrics: make(map[string]*StatTracker),
	}
}

// add a single instance of a reported Stat to that metric's Summary
func (s *Summary) add(stat *Stat) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.metrics[stat.Name]; !ok {
		s.metrics[stat.Name] = NewStatTracker()
	}
	s.metrics[stat.Name].add(stat)
}

// get a StatResponse for a single metric
func (s *Summary) get(name string) *StatResponse {
	s.RLock()
	defer s.RUnlock()
	return s.metrics[name].get()
}

func (s *Summary) getAllCount() map[string]int {
	s.RLock()
	defer s.RUnlock()

	r := make(map[string]int)
	for m, tracker := range s.metrics {
		r[m] = tracker.count
	}
	return r
}

// get a map of all metric names to corresponding StatResponses
func (s *Summary) getAllDetails() map[string]*StatResponse {
	s.RLock()
	defer s.RUnlock()

	r := make(map[string]*StatResponse)
	for m := range s.metrics {
		r[m] = s.get(m)
	}
	return r
}

// reset clears the tracked metrics
func (s *Summary) reset() {
	s.Lock()
	defer s.Unlock()

	s.metrics = make(map[string]*StatTracker)
}

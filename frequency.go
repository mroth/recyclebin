package main

import "sort"

type termCounter struct {
	_storage map[string]int
}

func NewTermCounter() *termCounter {
	tc := new(termCounter)
	tc._storage = make(map[string]int)
	return tc
}

func (c termCounter) Increment(key string) {
	c._storage[key] = c._storage[key] + 1
}

func (c termCounter) Scores() ScoreList {
	sl := make(ScoreList, len(c._storage))
	i := 0
	for k, v := range c._storage {
		sl[i] = Score{k, v}
		i++
	}
	return sl
}

func (c termCounter) SortedScores() ScoreList {
	sl := c.Scores()
	sort.Sort(c.Scores())
	return sl
}

type Score struct {
	Key   string
	Value int
}

type ScoreList []Score

func (p ScoreList) Len() int           { return len(p) }
func (p ScoreList) Less(i, j int) bool { return p[i].Value > p[j].Value }
func (p ScoreList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// First returns the first n Scores in a ScoreList p.
//
// If p.Len() < n, returns the first p.Len() Scores instead to prevent a slice panic.
func (p ScoreList) First(n int) ScoreList {
	if n <= 0 {
		return []Score{}
	}
	if p.Len() < n {
		return p
	}
	return p[:(n)]
}

// GreaterThan returns a new ScoreList with all Scores from p with Score.Value > n
func (p ScoreList) GreaterThan(n int) ScoreList {
	results := make(ScoreList, 0)
	for _, s := range p {
		if s.Value > n {
			results = append(results, s)
		}
	}
	return results
}

// Sorted returns a sorted *copy* of ScoreList p
//
// Useful for chaining when you don't want to sort in place.
func (p ScoreList) Sorted() ScoreList {
	dst := make(ScoreList, p.Len())
	copy(dst, p)
	sort.Sort(dst)
	return dst
}

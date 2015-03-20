package main

import (
	"fmt"
	"sync"
	"time"
)

type Log struct {
	when time.Time
	what []byte
}

func (a *Log) Copy() *Log {
	cp := make([]byte, len(a.what))
	copy(cp, a.what)
	return &Log{
		when: a.when,
		what: cp,
	}
}

type HistoryLog struct {
	numAbs          int
	numGen          int
	generateHistory []*Log
	absorbHistory   []*Log
	name            string
	mut             sync.Mutex
}

func (r *HistoryLog) GetHistory() *HistoryLog {
	return r.DeepCopy()
}

func NewHistoryLog(name string) *HistoryLog {
	r := &HistoryLog{
		generateHistory: make([]*Log, 0),
		absorbHistory:   make([]*Log, 0),
		name:            name,
	}
	return r
}

func (r *HistoryLog) DeepCopy() *HistoryLog {
	r.mut.Lock()
	defer r.mut.Unlock()

	s := &HistoryLog{
		generateHistory: make([]*Log, len(r.generateHistory)),
		absorbHistory:   make([]*Log, len(r.absorbHistory)),
		name:            r.name,
	}
	for _, v := range r.generateHistory {
		s.generateHistory = append(s.generateHistory, v.Copy())
	}

	return s
}

func (s *HistoryLog) RecordGen(what []byte) {
	s.mut.Lock()
	defer s.mut.Unlock()

	cp := make([]byte, len(what))
	copy(cp, what)
	s.generateHistory = append(s.generateHistory, &Log{when: time.Now(), what: cp})
	s.absorbHistory = append(s.absorbHistory, &Log{}) // make spacing apparent
	s.numGen++
}

func (s *HistoryLog) RecordAbs(what []byte) {
	s.mut.Lock()
	defer s.mut.Unlock()

	cp := make([]byte, len(what))
	copy(cp, what)
	s.absorbHistory = append(s.absorbHistory, &Log{when: time.Now(), what: cp})
	s.generateHistory = append(s.generateHistory, &Log{})
	s.numAbs++
}

func (s *HistoryLog) ShowHistory() {
	s.mut.Lock()
	defer s.mut.Unlock()

	fmt.Printf("%s history:\n", s.name)
	for i := 0; i < len(s.absorbHistory); i++ {
		if s.absorbHistory[i].when.IsZero() {

		} else {
			fmt.Printf("Abs @ %v: '%s'\n",
				s.absorbHistory[i].when,
				string(s.absorbHistory[i].what))
		}

		if s.generateHistory[i].when.IsZero() {

		} else {
			fmt.Printf("Gen @ %v:                  '%s'\n",
				s.generateHistory[i].when,
				string(s.generateHistory[i].what))
		}
	}
}

func (s *HistoryLog) CountAbsorbs() int {
	s.mut.Lock()
	defer s.mut.Unlock()
	return s.numAbs
}

func (s *HistoryLog) CountGenerates() int {
	s.mut.Lock()
	defer s.mut.Unlock()
	return s.numGen
}
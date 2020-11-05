package multiconfig

import (
	"testing"
)

func TestDefaultValues(t *testing.T) {
	m := &TagLoader{}
	s := new(Server)
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	if s.Port != getDefaultServer().Port {
		t.Errorf("Port value is wrong: %d, want: %d", s.Port, getDefaultServer().Port)
	}

	if s.Postgres.DBName != getDefaultServer().Postgres.DBName {
		t.Errorf("Postgres DBName value is wrong: %s, want: %s", s.Postgres.DBName, getDefaultServer().Postgres.DBName)
	}
}

func TestNestedValues(t *testing.T) {
	m := &TagLoader{}
	s := &Server{
		Name:    "hello",
		Labels:  []int{0, 1},
		Threads: []Thread{{}},
	}
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	if s.Threads[0].Name != "foobar" {
		t.Errorf("threads[0].name value is wrong: %s, want: %s", s.Threads[0].Name, "foobar")
	}
}

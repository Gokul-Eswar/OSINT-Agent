package core

// Collector is the interface that all passive collection plugins must implement.
type Collector interface {
	Name() string
	Description() string
	Collect(caseID string, target string) ([]Evidence, error)
}

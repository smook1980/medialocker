package medialocker

// Report
type Report struct {
	Errors    int32
	Completed int32
	Pending   int32
}

// MediaPathFinder
type MediaPathFinder interface {
	AddPath(string)
	EachImagePath(func(MediaPath) error)
	EachVideoPath(func(MediaPath) error)
	Report() Report
	Wait()
}

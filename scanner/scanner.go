package scanner

import (
	"os"
	"sync"

	"path/filepath"

	"github.com/smook1980/medialocker"
	. "github.com/smook1980/medialocker/types"
	"github.com/smook1980/medialocker/util"
)

// EachFunc is a function called for each MediaPath object found by Scanner
type EachFunc func(mp MediaPath)

// Scanner walks a base path for types.MediaPath objects.
type Scanner struct {
	basePath string
	errors   []error
	Log      medialocker.Log

	files    map[string]MediaPath
	in       chan MediaPath
	outs     []chan MediaPath
	outsWg   sync.WaitGroup
	workerWg sync.WaitGroup
}

// NewScanner returns a new Scanner for the given path with workers number of
// go routines processing file data.
func NewScanner(path string, workers int, log *medialocker.Logger) *Scanner {
	var bp string
	var err error

	if bp, err = filepath.Abs(path); err != nil {
		bp = path
	}

	logger := log.
		WithField("prefix", "scanner").
		WithField("root", path).
		WithField("worker_count", workers)

	s := &Scanner{
		basePath: bp,
		in:       make(chan MediaPath, workers+1),
		files:    map[string]MediaPath{},
		Log:      logger,
	}

	for x := 0; x < workers; x++ {
		go s.doWork()
	}

	return s
}

func (s *Scanner) doWork() {
	s.workerWg.Add(1)
	defer s.workerWg.Done()
	for mp := range s.in {
		if err := mp.Update(); err != nil {
			s.Log.Errorf("Failed to read details for %+v. %s", mp, err)
		} else {
			switch mp.Type {
			case Video, Image:
				for _, out := range s.outs {
					out <- mp
				}
			case Archive:
				// process archive here..
				s.Log.Debugf("doWork: mp is of type archive: %v\n", mp)
			}
		}
	}
}

// Each calls the given function for each MediaPath item found. Call multiple
// times to register multiple functions to be called for each MediaPath.
func (s *Scanner) Each(cb EachFunc) {
	out := make(chan MediaPath, 1)
	s.outs = append(s.outs, out)
	go func(cb EachFunc, out chan MediaPath, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		for mp := range out {
			cb(mp)
		}
	}(cb, out, &s.outsWg)
}

// Run begins the scan.
func (s *Scanner) Run() error {
	var errs []error

	wfn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errs = append(errs, err)
			return nil
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		if abs, err := filepath.Abs(path); err != nil {
			s.Log.Errorln("Error getting realpath for: ", path, " error: ", err)
		} else {
			path = abs
		}

		if mp, exist := s.files[path]; !exist {
			mp.FileInfo = info
			mp.Realpath = path
			s.files[path] = mp
			s.in <- mp
		}

		return nil
	}

	if err := filepath.Walk(s.basePath, wfn); err != nil {
		s.Log.Errorln("Failed walking path: ", s.basePath, ". ERROR: ", err)
		errs = append(errs, err)
	}

	close(s.in)
	s.Log.Infoln("Scanner done walking FS! Waiting for workers to finish...")
	s.workerWg.Wait()
	s.Log.Infoln("All scanner workers finished. Closing outs and waiting...")

	for _, out := range s.outs {
		close(out)
	}

	s.outsWg.Wait()

	s.Log.Infoln("Scanner all done!")

	if len(errs) > 0 {
		return util.MultiError(errs...)
	}

	return nil
}

func (s *Scanner) Module(app *medialocker.App) error {
	return s.Run()
}

package acceptance

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

// Every requirement test starts fitness-install, about 400 MB with its
// runner and checks. Several agents running `go test ./...` at once, each
// test binary running as many tests as there are cores, once put 147
// installers and 54 GB on a 36 GB Mac and froze it until the kernel
// restarted it. These limits hold the whole machine to machineSlots test
// binaries of perBinary tests each, whoever starts them: 24 installers,
// about 9 GB.
const (
	// machineSlots is how many requirement test binaries run at once on
	// one machine; FITNESS_TEST_SLOTS overrides it.
	machineSlots = 2
	// perBinary is the most tests one binary runs at once, even when the
	// command line asks for more with -parallel.
	perBinary = 12
)

// slotDir is shared by every user on the machine, so agents running as
// different users still wait for each other.
var slotDir = filepath.Join("/tmp", "fitness-test-slots")

// TakeTurn waits for one of the machine's test slots and caps how many tests
// this binary runs at once. Each folder's TestMain calls it before m.Run. The
// slot is a file lock, so it frees itself when the process exits, even on a
// crash.
func TakeTurn() {
	limitParallel()
	slots := machineSlots
	if n, err := strconv.Atoi(os.Getenv("FITNESS_TEST_SLOTS")); err == nil && n > 0 {
		slots = n
	}
	waited := false
	for {
		if _, ok := trySlot(slotDir, slots); ok {
			return
		}
		if !waited {
			fmt.Fprintf(os.Stderr, "acceptance: all %d machine test slots in %s are busy; waiting for one\n", slots, slotDir)
			waited = true
		}
		time.Sleep(time.Second)
	}
}

// limitParallel lowers -test.parallel to perBinary when the command line,
// or its default of the core count, asks for more.
func limitParallel() {
	if !flag.Parsed() {
		flag.Parse()
	}
	asked := runtime.GOMAXPROCS(0)
	if f := flag.Lookup("test.parallel"); f != nil {
		if n, err := strconv.Atoi(f.Value.String()); err == nil {
			asked = n
		}
	}
	_ = flag.Set("test.parallel", strconv.Itoa(min(perBinary, asked)))
}

// trySlot locks the first free one of slots lock files in dir without
// waiting, and returns the open file that holds the lock. It reports false
// when every slot is held.
func trySlot(dir string, slots int) (*os.File, bool) {
	_ = os.MkdirAll(dir, 0o777)
	_ = os.Chmod(dir, 0o777|os.ModeSticky)
	for i := range slots {
		path := filepath.Join(dir, fmt.Sprintf("slot-%d", i))
		f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0o666)
		if err != nil {
			continue
		}
		_ = os.Chmod(path, 0o666)
		if syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB) == nil {
			held = append(held, f)
			return f, true
		}
		_ = f.Close()
	}
	return nil, false
}

// held keeps the locked slot files open, and so locked, for the life of the
// process.
var held []*os.File

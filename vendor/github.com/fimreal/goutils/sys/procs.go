package sys

import (
	"runtime"

	"github.com/fimreal/goutils/ezap"
)

// Allow as many threads as we have cores unless the user specified a value.
func SetMaxProcs(maxProcs int) {
	var numProcs int
	if maxProcs < 1 {
		numProcs = runtime.NumCPU()
	} else {
		numProcs = maxProcs
	}
	runtime.GOMAXPROCS(numProcs)

	// Check if the setting was successful.
	actualNumProcs := runtime.GOMAXPROCS(0)
	if actualNumProcs != numProcs {
		ezap.Warnf("Specified max procs of %d but using %d", numProcs, actualNumProcs)
	}
}

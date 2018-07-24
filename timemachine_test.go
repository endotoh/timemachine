package timemachine

import (
	"testing"
	"time"
)

var SLEEP = time.Microsecond

func TestInit(t *testing.T) {
	fatalIf(IsFrozen(), t, "Was not expecting time to be frozen")
}

func TestRealNowIncreases(t *testing.T) {
	n := Now()
	diff := Since(n)
	errorIf(diff <= 0, t, "Now() should be monotonically increasing")
	fatalIf(IsFrozen(), t, "Was not expecting Now() to be frozen")
}

func TestRealNowWithRealSleep(t *testing.T) {
	n := Now()
	time.Sleep(SLEEP)
	diff := Since(n)
	fatalIf(diff <= 0, t, "time.Sleep()ing should move a clock under normal conditions")
}

func TestRealNowWithFakeSleep(t *testing.T) {
	n := Now()
	Sleep(SLEEP)
	diff := Since(n)
	fatalIf(diff <= 0, t, "timemachine.Sleep()ing should also move a clock under normal conditions")
}

func TestFrozenNow(t *testing.T) {
	fatalIf(IsFrozen(), t, "Time should not be frozen")

	var n time.Time
	func() {
		n := FreezeNow()
		defer Unfreeze()
		fatalIf(!IsFrozen(), t, "Expecting time to be frozen now")

		nn := Now()
		fatalIf(n != nn, t, "Now() should return a canned time")

		diff := Since(n)
		fatalIf(diff != 0, t, "Now() should be frozen, time no longer moves forward")

		time.Sleep(SLEEP)
		diff = Since(n)
		fatalIf(diff != 0, t, "time.Sleep()ing should no longer move a clock forward")

		nff := Now()
		until := Until(nff)
		fatalIf(until != 0, t, "Until() should be frozen because Now() is frozen")

		Sleep(SLEEP)
		diff = Since(n)
		fatalIf(diff != SLEEP, t, "timemachine.Sleep()ing should move clock forward a precise amount")

	}()

	fatalIf(IsFrozen(), t, "defer Unfreeze() should take effect by end of function closure")
	diff := Since(n)
	errorIf(diff <= SLEEP, t, "Now() should be increasing again, time.Sleep() and running code had an effect")
	n2 := Now()
	n2f := n2.Add(SLEEP)
	until := Until(n2f)
	fatalIf(until >= SLEEP || until <= 0, t,
		"With unfrozen time, timemachine.Sleep() should be behaving just like time.Sleep() (got: %v, should be larger than: %v)", until, SLEEP)
}

func TestTimeTravel(t *testing.T) {
	fatalIf(IsFrozen(), t, "Time should not be frozen")
	f := func() {
		defer func() {
			if recover() != nil {
				t.Log("Successfully catching deliberate panic when calling timemachine.Travel() when time is not frozen")
			}
		}()
		Travel(SLEEP)
	}
	f()

	_ = FreezeNow()
	defer Unfreeze()
	Travel(SLEEP)

}

//////////////////////////////////////////////////////////////
// helper functions

func errorIf(condition bool, t *testing.T, s string, args ...interface{}) {
	if condition {
		t.Errorf(s, args...)
	}
}
func fatalIf(condition bool, t *testing.T, s string, args ...interface{}) {
	if condition {
		t.Fatalf(s, args...)
	}
}

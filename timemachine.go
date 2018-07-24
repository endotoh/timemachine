package timemachine

import (
	"time"
)

var frozen bool
var frozenTime time.Time

///////////////////////////////////////////////////
// Swap-ins for time.* functions

// Now() is a swap-in replacement for time.Now() which can be overriden
// during testing by invoking Freeze(). In non-testing, production code, it is indistinguishable from time.Now(), also returning time.Time objects
func Now() time.Time {
	if frozen {
		//fmt.Println("Using frozen time: ", frozenTime)
		return frozenTime
	} else {
		return time.Now()
	}
}

func Sleep(d time.Duration) {
	if frozen {
		//fmt.Printf("Artificially moving time forward by %v\n", d)
		frozenTime = frozenTime.Add(d)
	} else {
		time.Sleep(d)
	}
}

func Since(t time.Time) time.Duration {
	return Now().Sub(t)
}

func Until(t time.Time) time.Duration {
	return t.Sub(Now())
}

///////////////////////////////////////////////////
// functions specific to timemachine

func FreezeNow() time.Time {
	frozen = true
	frozenTime = time.Now()
	return frozenTime
}

func Unfreeze() {
	frozen = false
}

func IsFrozen() bool {
	return frozen
}

func Travel(d time.Duration) time.Time {
	if !frozen {
		panic("You can only time travel after calling FreezeNow()")
	} else {
		frozenTime = frozenTime.Add(d)
		return frozenTime
	}
}

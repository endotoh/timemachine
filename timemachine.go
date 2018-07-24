// Package timemachine is a testing friendly drop-in replacement for Go's standard time package
//
// In production code, just use identically named static functions in time package
// like Now(), Since(), Until() and Sleep(). They all use and return the same
// time.Time and time.Duration types you know and love.
//
// In testing code, you can freeze and unfreeze time, sleep without incurring real wall clock time
// (i.e. making tests slower) and even travel forward in time. A typical test:
//
// func TestMy24HourExpiryBizLogic(t *testing.T) {
//   _ = timemachine.FreezeNow()
//   defer timemachine.Unfreeze()
//   x := mypkg.Init()
//   if x.IsExpired() {
//     t.Error("New objects should never be expired")
//   }
//   timemachine.Travel(24 * time.Hour + time.Second)
//   if ! x.IsExpired() {
//     t.Error("Just over a day old objects should be expired")
//   }
// }
//
// Inspired by HTTP mocking library Gock:
//     https://github.com/h2non/gock
//
// NOT concurrency safe in testing code, but that's probably OK for many use
// cases as test code is more often single threaded. More sophisticated alternatives
// which rely on objects instead can be found:
//     https://www.reddit.com/r/golang/comments/530cqp/how_to_mock_time_for_testing_purposes/
//
package timemachine

import (
	"time"
)

var frozen bool
var frozenTime time.Time

//////////////////////////////////////////////////////////////////////////
// Swap-ins for time.* functions

// Now behaves like time.Now() unless FreezeNow() has been called. In which
// case, it returns a cached time.Time object which only changes through
// Sleep() and Travel() functions
func Now() time.Time {
	if frozen {
		//fmt.Println("Using frozen time: ", frozenTime)
		return frozenTime
	} else {
		return time.Now()
	}
}

// Sleep behaves just like time.Sleep() unless FreezeNow has been called.
// In which case, it does not actually sleep it just moves the cached time forward.
func Sleep(d time.Duration) {
	if frozen {
		//fmt.Printf("Artificially moving time forward by %v\n", d)
		frozenTime = frozenTime.Add(d)
	} else {
		time.Sleep(d)
	}
}

// Since should be used instead of time.Since() if you are using this library, as it
// depends on Now() and FreezeNow() functions
func Since(t time.Time) time.Duration {
	return Now().Sub(t)
}

// Until should be used instead of time.Until() if you are using this library, as it
// depends on Now() and FreezeNow() functions
func Until(t time.Time) time.Duration {
	return t.Sub(Now())
}

//////////////////////////////////////////////////////////////////////////
// functions specific to timemachine

// FreezeNow should be used in tests to trigger this library's core behaviour,
// caching time.Now(). You should ONLY use this in test code.
func FreezeNow() time.Time {
	frozen = true
	frozenTime = time.Now()
	return frozenTime
}

// Unfreeze cleans things up, reverting to production mode. Use the FreezeNow(), defer Unfreeze()
// idiom.
func Unfreeze() {
	frozen = false
}

// IsFrozen tells you if FreezeNow() has been called without Unfreeze()
func IsFrozen() bool {
	return frozen
}

// Travel allows you to increment cached time by time.Duration. Only intended for test mode, not
// production mode. Panic's if called outside FreezeNow() and Unfreeze() block.
// You more explicitly communicate your intent using Travel() than Sleep().
func Travel(d time.Duration) time.Time {
	if !frozen {
		panic("You can only time travel after calling FreezeNow()")
	} else {
		frozenTime = frozenTime.Add(d)
		return frozenTime
	}
}

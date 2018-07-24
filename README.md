# timemachine
Mockable time (therefore time travel) in Golang, inspired by Gock

## INTRODUCTION

I didn't find something super simple and lightweight for mocking time
for frequent Golang testing use cases which do not require thread safety. This
does the job for me, it might for you too.

## USAGE

From godoc:

```
PACKAGE DOCUMENTATION

package timemachine

    Package timemachine is a testing friendly drop-in replacement for Go's
    standard time package

    In production code, just use identically named static functions in time
    package like Now(), Since(), Until() and Sleep(). They all use and
    return the same time.Time and time.Duration types you know and love.

    In testing code, you can freeze and unfreeze time, sleep without
    incurring real wall clock time (i.e. making tests slower) and even
    travel forward in time. A typical test:

    func TestMy24HourExpiryBizLogic(t *testing.T) {

	_ = timemachine.FreezeNow()
	defer timemachine.Unfreeze()
	x := mypkg.Init()
	if x.IsExpired() {
	  t.Error("New objects should never be expired")
	}
	timemachine.Travel(24 * time.Hour + time.Second)
	if ! x.IsExpired() {
	  t.Error("Just over a day old objects should be expired")
	}

    }

    Inspired by HTTP mocking library Gock:

	https://github.com/h2non/gock

    NOT concurrency safe in testing code, but that's probably OK for many
    use cases as test code is more often single threaded. More sophisticated
    alternatives which rely on objects instead can be found:

	https://www.reddit.com/r/golang/comments/530cqp/how_to_mock_time_for_testing_purposes/

FUNCTIONS

func FreezeNow() time.Time
    FreezeNow should be used in tests to trigger this library's core
    behaviour, caching time.Now(). You should ONLY use this in test code.

func IsFrozen() bool
    IsFrozen tells you if FreezeNow() has been called without Unfreeze()

func Now() time.Time
    Now behaves like time.Now() unless FreezeNow() has been called. In which
    case, it returns a cached time.Time object which only changes through
    Sleep() and Travel() functions

func Since(t time.Time) time.Duration
    Since should be used instead of time.Since() if you are using this
    library, as it depends on Now() and FreezeNow() functions

func Sleep(d time.Duration)
    Sleep behaves just like time.Sleep() unless FreezeNow has been called.
    In which case, it does not actually sleep it just moves the cached time
    forward.

func Travel(d time.Duration) time.Time
    Travel allows you to increment cached time by time.Duration. Only
    intended for test mode, not production mode. Panic's if called outside
    FreezeNow() and Unfreeze() block. You more explicitly communicate your
    intent using Travel() than Sleep().

func Unfreeze()
    Unfreeze cleans things up, reverting to production mode. Use the
    FreezeNow(), defer Unfreeze() idiom.

func Until(t time.Time) time.Duration
    Until should be used instead of time.Until() if you are using this
    library, as it depends on Now() and FreezeNow() functions

```

## SEE ALSO

Alternatives dicussed on Reddit:

https://www.reddit.com/r/golang/comments/530cqp/how_to_mock_time_for_testing_purposes/


## LICENSE

ASL 2.0 - https://www.apache.org/licenses/LICENSE-2.0

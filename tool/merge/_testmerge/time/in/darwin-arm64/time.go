// Selected declarations from the generated time.h binding.
package libc

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

// time()
type X__DarwinTimeT = c.Long

type TimeT = X__DarwinTimeT

// ISO/IEC 9899:201x 7.27.2.5 The timespec_get function
const TIME_UTC = 1

// Clock aliases.
type (
	// TickCount records clock ticks.
	//
	// The Darwin clock supplies the count.
	TickCount = c.Int // Measured in ticks.
)

// TimeSample stores a sampled value.
//
// The Darwin binding supplies its own sampling routine.
type TimeSample struct {
	Value c.Int
} // Darwin sample.

// Config accepts these options.
//
//   - Use the Darwin clock.
type Config struct {
	// Flags selects the sampling mode.
	//
	// The Darwin binding defines its own flags.
	Flags uint `json:"flags"` /* flags */
	_llcppg_anon_0 // The embedded struct.
}

// The embedded struct.
type _llcppg_anon_0 struct {
	Extra int
}

type Tm struct {
	TmSec    c.Int
	TmMin    c.Int
	TmHour   c.Int
	TmMday   c.Int
	TmMon    c.Int
	TmYear   c.Int
	TmWday   c.Int
	TmYday   c.Int
	TmIsdst  c.Int
	TmGmtoff c.Long
	TmZone   *c.Char
}

// Daylight is nonzero if the time zone observes daylight saving time.
//
// This variable is exported by the Darwin C library.
//
//go:linkname Daylight C.daylight
var Daylight c.Int

//llgo:link (*Tm).Asctime C.asctime
func (self *Tm) Asctime() *c.Char {
	return nil
}

// Mktime interprets the receiver in the local time zone.
//
// The receiver may be updated.
//
//llgo:link (*Tm).Mktime C.mktime
func (self *Tm) Mktime() TimeT {
	return 0
}

// Set time conversion information from the TZ environment variable.
//
// Darwin can derive the default from the system time zone configuration.
//
//go:linkname Tzset C.tzset
func Tzset()

type Timespec struct {
	TvSec  X__DarwinTimeT
	TvNsec c.Long
}

type ClockidT c.Uint

//go:linkname Tzname C.tzname
var Tzname **c.Char

//llgo:link ClockidT.ClockGettime C.clock_gettime
func (self ClockidT) ClockGettime(__tp *Timespec) c.Int {
	return 0
}

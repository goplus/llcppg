// Selected declarations from the generated time.h binding.
package libc

import (
	libc "github.com/goplus/lib/c"
	_ "unsafe"
)

// Seconds since the Epoch.
type X__TimeT = libc.Long

// Returned by `time'.
type TimeT = X__TimeT

// Time base values for timespec_get.
const TIME_UTC = 1

// Clock aliases.
type (
	// TickCount records clock ticks.
	//
	// The Linux clock supplies the count.
	TickCount = libc.Int // Measured in ticks.
)

type (
	// TimeSample stores a sampled value.
	//
	// The Linux binding supplies its own sampling routine.
	TimeSample struct {
		Value libc.Int
	} // Linux sample.
)

// Config accepts these options.
//
//   - Use the Linux clock.
type Config struct {
	// Flags selects the sampling mode.
	//
	// The Linux binding defines its own flags.
	Flags uint `json:"flags"` /* flags */
	_llcppg_anon_0 // The Linux embedded struct.
}

// The embedded struct.
type _llcppg_anon_0 struct {
	Extra int
}

// Used by other time functions.
type Tm struct {
	TmSec    libc.Int
	TmMin    libc.Int
	TmHour   libc.Int
	TmMday   libc.Int
	TmMon    libc.Int
	TmYear   libc.Int
	TmWday   libc.Int
	TmYday   libc.Int
	TmIsdst  libc.Int
	TmGmtoff libc.Long
	TmZone   *libc.Char
}

// Return the `time_t' representation of TP and normalize TP.
//
// The receiver may be updated.
//
//llgo:link (*Tm).Mktime C.mktime
func (self *Tm) Mktime() TimeT {
	return 0
}

// Return a string of the form "Day Mon dd hh:mm:ss yyyy\n"
// that is the representation of TP in this format.
//
//llgo:link (*Tm).Asctime C.asctime
func (self *Tm) Asctime() *libc.Char {
	return nil
}

// Set time conversion information from the TZ environment variable.
//
// If TZ is not defined, a locale-dependent default is used.
//
//go:linkname Tzset C.tzset
func Tzset()

var (
	// Daylight is nonzero if the time zone observes daylight saving time.
	//
	//go:linkname Daylight C.daylight
	Daylight libc.Int
)

//go:linkname Timezone C.timezone
var Timezone libc.Long // Seconds west of UTC.

type X__ClockidT = libc.Int

type X__SyscallSlongT = libc.Long

type Timespec struct {
	TvSec  X__TimeT
	TvNsec X__SyscallSlongT
}

type ClockidT = X__ClockidT

//go:linkname Tzname C.tzname
var Tzname [2]*libc.Char

//go:linkname ClockGettime C.clock_gettime
func ClockGettime(__clock_id ClockidT, __tp *Timespec) libc.Int

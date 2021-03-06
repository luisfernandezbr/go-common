package datetime

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const (
	// SignalTimeUnitNOW is now
	SignalTimeUnitNOW int32 = 0
	// SignalTimeUnitMONTH is 30 days
	SignalTimeUnitMONTH int32 = 30
	// SignalTimeUnitQUARTER is 90 days
	SignalTimeUnitQUARTER int32 = 90
	// SignalTimeUnitHALFYEAR is 180 days
	SignalTimeUnitHALFYEAR int32 = 180
	// SignalTimeUnitYEAR is 365 days
	SignalTimeUnitYEAR int32 = 365
	// SignalTimeUnitTHIRDQUARTER is 270 days
	SignalTimeUnitTHIRDQUARTER int32 = 270
	// SignalTimeUnitALLTIME is all time
	SignalTimeUnitALLTIME int32 = -1
	// SignalTimeUnitBIMONTH is two month
	SignalTimeUnitBIMONTH int32 = 60
	// DaysInMilliseconds is one day in milliseconds
	DaysInMilliseconds int64 = 86400000
)

// RFC3339 There is an edge case where there is no timezone set
// and that makes the format function return dates like
// 2019-09-03T20:48:57.073Z.
// We are using this RFC3339 custom format to
// always get an offset.
const RFC3339 = "2006-01-02T15:04:05.999999999-07:00"

// GetTimeUnitString will return the timeunit as a string
func GetTimeUnitString(timeUnit int32) string {
	switch timeUnit {
	case SignalTimeUnitNOW:
		{
			return "now"
		}
	case SignalTimeUnitMONTH:
		{
			return "month"
		}
	case SignalTimeUnitQUARTER:
		{
			return "quarter"
		}
	case SignalTimeUnitBIMONTH:
		{
			return "bimonth"
		}
	case SignalTimeUnitHALFYEAR:
		{
			return "halfyear"
		}
	case SignalTimeUnitTHIRDQUARTER:
		{
			return "thirdquarter"
		}
	case SignalTimeUnitYEAR:
		{
			return "year"
		}
	}
	return "alltime"
}

// ISODate returns a RFC 3339 formatted string for the current date time
func ISODate() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// ISODateFromTime returns a RFC 3339 formatted string from the supplied timestamp
func ISODateFromTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// ISODateToTime returns a RFC 3339 formatted string as a timestamp
func ISODateToTime(date string) (time.Time, error) {
	if strings.HasSuffix(date, "Z") {
		return time.Parse("2006-01-02T15:04:05Z", date)
	}
	return ISODateOffsetToTime(date)
}

// ISODateOffsetToTime returns a RFC 3339 formatted string as a timestamp
func ISODateOffsetToTime(date string) (time.Time, error) {
	if date == "" {
		return time.Time{}, nil
	}
	if strings.Contains(date, "Z") {
		// 2017-01-20T15:56:23.000000Z-08:00
		tv, err := time.Parse("2006-01-02T15:04:05.999999999Z-07:00", date)
		if err == nil {
			return tv, nil
		}
	}
	if strings.Contains(date, ".") {
		tv, err := time.Parse("2006-01-02T15:04:05.999999999-07:00", date)
		if err != nil {
			return time.Parse("2006-01-02T15:04:05.999999999-0700", date)
		}
		return tv, nil
	}

	match, _ := regexp.MatchString("([+-]\\d{2}:\\d{2})", date)
	if match {
		return time.Parse("2006-01-02T15:04:05-07:00", date)
	}

	return time.Parse("2006-01-02T15:04:05-0700", date)
}

// ISODateToEpoch returns an epoch date or 0 if invalid or empty
func ISODateToEpoch(date string) (int64, error) {
	if date == "" {
		return 0, nil
	}
	ts, err := ISODateToTime(date)
	if err != nil {
		return 0, err
	}
	return TimeToEpoch(ts), nil
}

// TimeToEpoch will convert a time to epoch (in UTC) with millisecond precision
func TimeToEpoch(tv time.Time) int64 {
	if tv.IsZero() {
		return 0
	}
	tv = tv.UTC()
	// we want to round down to microsecond precision from nano second before we return as milliseconds
	// so we can get the microseconds in the value of epoch
	return (tv.UnixNano() + 500000) / 1000000
}

// EpochNow will return the current time in epoch (in UTC) with millisecond precision
func EpochNow() int64 {
	return TimeToEpoch(time.Now())
}

// DateFromEpoch returns a time.Time from an epoch value in milliseconds
func DateFromEpoch(t int64) time.Time {
	return time.Unix(0, t*1000000)
}

// ShortDateFromEpoch will return a short date from a epoch value in milliseconds
func ShortDateFromEpoch(t int64) string {
	tv := DateFromEpoch(t)
	return tv.UTC().Format("2006-01-02")
}

// ShortDateFromTime will return a short date from a time
func ShortDateFromTime(tv time.Time) string {
	return tv.UTC().Format("2006-01-02")
}

// ShortDate returns a DATE (no time) formatted string from RFC 3339 formatted string
func ShortDate(date string) string {
	if strings.Contains(date, "T") {
		t, err := time.Parse("2006-01-02T15:04:05Z", date)
		if err != nil {
			return fmt.Sprintf("<error parsing date: %s. %v>", date, err)
		}
		return t.UTC().Format("2006-01-02")
	}
	return date
}

// DateRange will return the beginning and end of a date range for a given time unit.
// if the timeunit is -1 it will return epoch zero
func DateRange(ref time.Time, timeunit int64) (int64, int64) {
	end := EndofDay(TimeToEpoch(ref))
	begin := 1000 + (end - DaysInMilliseconds*timeunit)
	if timeunit == int64(-1) {
		begin = 0
	}
	return StartofDay(begin), end
}

// DateRangePrevious will previous time date for the previous range from timeunit
func DateRangePrevious(ref int64, timeunit int64) (int64, int64) {
	end := EndofDay(ref)
	begin := (end - DaysInMilliseconds*timeunit)
	priorstart, priorend := DateRange(DateFromEpoch(begin), timeunit)
	return priorstart, priorend // we go to the next day
}

// EndofDay returns the end of the day (midnight) for a given epoch time
func EndofDay(tv int64) int64 {
	t := DateFromEpoch(tv).UTC()
	return TimeToEpoch(time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 9999, time.UTC))
}

// StartofDay returns the start of the day (just after midnight) for a given epoch time
func StartofDay(tv int64) int64 {
	t := DateFromEpoch(tv).UTC()
	return TimeToEpoch(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC))
}

// ToTimeRange will return a start and end time range in epoch using tv as the reference day and adding days (use negative number to subtract)
func ToTimeRange(tv time.Time, days int) (int64, int64) {
	tv = tv.UTC()
	end := time.Date(tv.Year(), tv.Month(), tv.Day(), 23, 59, 59, 9999, time.UTC)
	startend := time.Date(tv.Year(), tv.Month(), tv.Day(), 0, 0, 0, 0, time.UTC)
	start := startend.AddDate(0, 0, days)
	return TimeToEpoch(start), TimeToEpoch(end)
}

// GetSignalDate returns a metric date in short form for a time unit from the ref date
func GetSignalDate(timeUnit int32, refDate time.Time) string {
	switch timeUnit {
	case SignalTimeUnitNOW:
		{
			return ShortDateFromTime(refDate)
		}
	case SignalTimeUnitMONTH:
		{
			return ShortDateFromTime(refDate.AddDate(0, 0, -30))
		}
	case SignalTimeUnitBIMONTH:
		{
			return ShortDateFromTime(refDate.AddDate(0, 0, -60))
		}
	case SignalTimeUnitQUARTER:
		{
			return ShortDateFromTime(refDate.AddDate(0, 0, -90))
		}
	case SignalTimeUnitHALFYEAR:
		{
			return ShortDateFromTime(refDate.AddDate(0, 0, -180))
		}
	case SignalTimeUnitTHIRDQUARTER:
		{
			return ShortDateFromTime(refDate.AddDate(0, 0, -270))
		}
	case SignalTimeUnitYEAR:
		{
			return ShortDateFromTime(refDate.AddDate(0, 0, -365))
		}
	}
	return ""
}

// GetSignalTime returns a metric date for a time unit from the ref date
// This will be changed from -30 to -29 because of the next example
// Lets say I want to take 1 day from "yesterday at the end" lets say "2017-02-26 23:59:59.9999"
// that will result in "2017-02-25 23:59:59.9999" and after truncate "2017-02-25 00:00:00" so it is actually taken 2 days
func GetSignalTime(timeUnit int32, refDate time.Time) time.Time {
	var t time.Time
	switch timeUnit {
	case SignalTimeUnitNOW:
		{
			return refDate.UTC().Truncate(time.Hour * 24)
		}
	case SignalTimeUnitMONTH:
		{
			t = refDate.UTC().AddDate(0, 0, -30)
		}
	case SignalTimeUnitBIMONTH:
		{
			t = refDate.UTC().AddDate(0, 0, -60)
		}
	case SignalTimeUnitQUARTER:
		{
			t = refDate.UTC().AddDate(0, 0, -90)
		}
	case SignalTimeUnitHALFYEAR:
		{
			t = refDate.UTC().AddDate(0, 0, -180)
		}
	case SignalTimeUnitTHIRDQUARTER:
		{
			t = refDate.UTC().AddDate(0, 0, -270)
		}
	case SignalTimeUnitYEAR:
		{
			t = refDate.UTC().AddDate(0, 0, -365)
		}
	}

	return t.Truncate(time.Hour * 24)
}

// ToMilliSec Convert time to milliseconds int64
func ToMilliSec(date time.Time) int64 {
	return date.UnixNano() / 1000000
}

// AddDaysToStrDate will add days to string date
func AddDaysToStrDate(date string, days int) (string, error) {
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", err
	}
	return ShortDateFromTime(d.AddDate(0, 0, days)), nil
}

// Date represents the object structure for date
type Date struct {
	// Epoch the date in epoch format
	Epoch int64 `json:"epoch" bson:"epoch" yaml:"epoch" faker:"-"`
	// Offset the timezone offset from GMT
	Offset int64 `json:"offset" bson:"offset" yaml:"offset" faker:"-"`
	// Rfc3339 the date in RFC3339 format
	Rfc3339 string `json:"rfc3339" bson:"rfc3339" yaml:"rfc3339" faker:"-"`
}

// NewDateNow returns a Date object as of now
func NewDateNow() Date {
	epoch := EpochNow()
	val := DateFromEpoch(epoch).Format(RFC3339)
	tv, _ := ISODateToTime(val)
	_, timezone := tv.Zone()
	return Date{
		Epoch:   epoch,
		Rfc3339: val,
		Offset:  int64(timezone) / 60,
	}
}

// NewDate returns a new Date object from a string date value
func NewDate(val string) (*Date, error) {
	tv, err := ISODateToTime(val)
	if err != nil {
		return nil, err
	}
	_, timezone := tv.Zone()
	return &Date{
		Epoch:   TimeToEpoch(tv),
		Rfc3339: tv.Round(time.Millisecond).Format(RFC3339),
		Offset:  int64(timezone) / 60,
	}, nil
}

// NewDateWithTime returns a new Date object from a time.Time value
func NewDateWithTime(tv time.Time) *Date {
	_, timezone := tv.Zone()
	return &Date{
		Epoch:   TimeToEpoch(tv),
		Rfc3339: tv.Round(time.Millisecond).Format(RFC3339),
		Offset:  int64(timezone) / 60,
	}
}

// NewDateFromEpoch returns a new Date object from a epoch time value
func NewDateFromEpoch(epoch int64) Date {
	val := DateFromEpoch(epoch).Format(RFC3339)
	tv, _ := ISODateToTime(val)
	_, timezone := tv.Zone()
	return Date{
		Epoch:   epoch,
		Rfc3339: val,
		Offset:  int64(timezone) / 60,
	}
}

// TimeFromDate returns a time.Time for a date
func TimeFromDate(date Date) time.Time {
	ts := DateFromEpoch(date.Epoch)
	if ts.IsZero() {
		return ts
	}
	// apply timezone
	loc := time.FixedZone("", int(date.Offset*60))
	return ts.In(loc)
}

// EpochMinuteApart returns true if both epochs are less than or equal
// to one minute apart
func EpochMinuteApart(epoch1, epoch2 int64) bool {
	big := epoch1
	small := epoch2
	// return tv1.Truncate(time.Minute).Equal(tv2.Truncate(time.Minute))
	// get the names right
	if small > big {
		big = epoch2
		small = epoch1
	}
	return big-small <= 1000*60
}

// ConvertToModel will fill dateModel based on passed time
func ConvertToModel(ts time.Time, dateModel interface{}) {
	if ts.IsZero() {
		return
	}

	date := NewDateWithTime(ts)

	t := reflect.ValueOf(dateModel).Elem()
	t.FieldByName("Rfc3339").Set(reflect.ValueOf(date.Rfc3339))
	t.FieldByName("Epoch").Set(reflect.ValueOf(date.Epoch))
	t.FieldByName("Offset").Set(reflect.ValueOf(date.Offset))
}

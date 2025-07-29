package stdlib

import (
	"time"

	"github.com/malivvan/vv/vvm"
)

var timesModule = map[string]vvm.Object{
	"format_ansic":        &vvm.String{Value: time.ANSIC},
	"format_unix_date":    &vvm.String{Value: time.UnixDate},
	"format_ruby_date":    &vvm.String{Value: time.RubyDate},
	"format_rfc822":       &vvm.String{Value: time.RFC822},
	"format_rfc822z":      &vvm.String{Value: time.RFC822Z},
	"format_rfc850":       &vvm.String{Value: time.RFC850},
	"format_rfc1123":      &vvm.String{Value: time.RFC1123},
	"format_rfc1123z":     &vvm.String{Value: time.RFC1123Z},
	"format_rfc3339":      &vvm.String{Value: time.RFC3339},
	"format_rfc3339_nano": &vvm.String{Value: time.RFC3339Nano},
	"format_kitchen":      &vvm.String{Value: time.Kitchen},
	"format_stamp":        &vvm.String{Value: time.Stamp},
	"format_stamp_milli":  &vvm.String{Value: time.StampMilli},
	"format_stamp_micro":  &vvm.String{Value: time.StampMicro},
	"format_stamp_nano":   &vvm.String{Value: time.StampNano},
	"nanosecond":          &vvm.Int{Value: int64(time.Nanosecond)},
	"microsecond":         &vvm.Int{Value: int64(time.Microsecond)},
	"millisecond":         &vvm.Int{Value: int64(time.Millisecond)},
	"second":              &vvm.Int{Value: int64(time.Second)},
	"minute":              &vvm.Int{Value: int64(time.Minute)},
	"hour":                &vvm.Int{Value: int64(time.Hour)},
	"january":             &vvm.Int{Value: int64(time.January)},
	"february":            &vvm.Int{Value: int64(time.February)},
	"march":               &vvm.Int{Value: int64(time.March)},
	"april":               &vvm.Int{Value: int64(time.April)},
	"may":                 &vvm.Int{Value: int64(time.May)},
	"june":                &vvm.Int{Value: int64(time.June)},
	"july":                &vvm.Int{Value: int64(time.July)},
	"august":              &vvm.Int{Value: int64(time.August)},
	"september":           &vvm.Int{Value: int64(time.September)},
	"october":             &vvm.Int{Value: int64(time.October)},
	"november":            &vvm.Int{Value: int64(time.November)},
	"december":            &vvm.Int{Value: int64(time.December)},
	"sleep": &vvm.BuiltinFunction{
		Name:      "sleep",
		Value:     timesSleep,
		NeedVMObj: true,
	}, // sleep(int)
	"parse_duration": &vvm.UserFunction{
		Name:  "parse_duration",
		Value: timesParseDuration,
	}, // parse_duration(str) => int
	"since": &vvm.UserFunction{
		Name:  "since",
		Value: timesSince,
	}, // since(time) => int
	"until": &vvm.UserFunction{
		Name:  "until",
		Value: timesUntil,
	}, // until(time) => int
	"duration_hours": &vvm.UserFunction{
		Name:  "duration_hours",
		Value: timesDurationHours,
	}, // duration_hours(int) => float
	"duration_minutes": &vvm.UserFunction{
		Name:  "duration_minutes",
		Value: timesDurationMinutes,
	}, // duration_minutes(int) => float
	"duration_nanoseconds": &vvm.UserFunction{
		Name:  "duration_nanoseconds",
		Value: timesDurationNanoseconds,
	}, // duration_nanoseconds(int) => int
	"duration_seconds": &vvm.UserFunction{
		Name:  "duration_seconds",
		Value: timesDurationSeconds,
	}, // duration_seconds(int) => float
	"duration_string": &vvm.UserFunction{
		Name:  "duration_string",
		Value: timesDurationString,
	}, // duration_string(int) => string
	"month_string": &vvm.UserFunction{
		Name:  "month_string",
		Value: timesMonthString,
	}, // month_string(int) => string
	"date": &vvm.UserFunction{
		Name:  "date",
		Value: timesDate,
	}, // date(year, month, day, hour, min, sec, nsec) => time
	"now": &vvm.UserFunction{
		Name:  "now",
		Value: timesNow,
	}, // now() => time
	"parse": &vvm.UserFunction{
		Name:  "parse",
		Value: timesParse,
	}, // parse(format, str) => time
	"unix": &vvm.UserFunction{
		Name:  "unix",
		Value: timesUnix,
	}, // unix(sec, nsec) => time
	"add": &vvm.UserFunction{
		Name:  "add",
		Value: timesAdd,
	}, // add(time, int) => time
	"add_date": &vvm.UserFunction{
		Name:  "add_date",
		Value: timesAddDate,
	}, // add_date(time, years, months, days) => time
	"sub": &vvm.UserFunction{
		Name:  "sub",
		Value: timesSub,
	}, // sub(t time, u time) => int
	"after": &vvm.UserFunction{
		Name:  "after",
		Value: timesAfter,
	}, // after(t time, u time) => bool
	"before": &vvm.UserFunction{
		Name:  "before",
		Value: timesBefore,
	}, // before(t time, u time) => bool
	"time_year": &vvm.UserFunction{
		Name:  "time_year",
		Value: timesTimeYear,
	}, // time_year(time) => int
	"time_month": &vvm.UserFunction{
		Name:  "time_month",
		Value: timesTimeMonth,
	}, // time_month(time) => int
	"time_day": &vvm.UserFunction{
		Name:  "time_day",
		Value: timesTimeDay,
	}, // time_day(time) => int
	"time_weekday": &vvm.UserFunction{
		Name:  "time_weekday",
		Value: timesTimeWeekday,
	}, // time_weekday(time) => int
	"time_hour": &vvm.UserFunction{
		Name:  "time_hour",
		Value: timesTimeHour,
	}, // time_hour(time) => int
	"time_minute": &vvm.UserFunction{
		Name:  "time_minute",
		Value: timesTimeMinute,
	}, // time_minute(time) => int
	"time_second": &vvm.UserFunction{
		Name:  "time_second",
		Value: timesTimeSecond,
	}, // time_second(time) => int
	"time_nanosecond": &vvm.UserFunction{
		Name:  "time_nanosecond",
		Value: timesTimeNanosecond,
	}, // time_nanosecond(time) => int
	"time_unix": &vvm.UserFunction{
		Name:  "time_unix",
		Value: timesTimeUnix,
	}, // time_unix(time) => int
	"time_unix_nano": &vvm.UserFunction{
		Name:  "time_unix_nano",
		Value: timesTimeUnixNano,
	}, // time_unix_nano(time) => int
	"time_format": &vvm.UserFunction{
		Name:  "time_format",
		Value: timesTimeFormat,
	}, // time_format(time, format) => string
	"time_location": &vvm.UserFunction{
		Name:  "time_location",
		Value: timesTimeLocation,
	}, // time_location(time) => string
	"time_string": &vvm.UserFunction{
		Name:  "time_string",
		Value: timesTimeString,
	}, // time_string(time) => string
	"is_zero": &vvm.UserFunction{
		Name:  "is_zero",
		Value: timesIsZero,
	}, // is_zero(time) => bool
	"to_local": &vvm.UserFunction{
		Name:  "to_local",
		Value: timesToLocal,
	}, // to_local(time) => time
	"to_utc": &vvm.UserFunction{
		Name:  "to_utc",
		Value: timesToUTC,
	}, // to_utc(time) => time
}

func timesSleep(args ...vvm.Object) (ret vvm.Object, err error) {
	vm := args[0].(*vvm.VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	i1, ok := vvm.ToInt64(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}
	ret = vvm.UndefinedValue
	if time.Duration(i1) <= time.Second {
		time.Sleep(time.Duration(i1))
		return
	}

	done := make(chan struct{})
	go func() {
		time.Sleep(time.Duration(i1))
		select {
		case <-vm.AbortChan:
		case done <- struct{}{}:
		}
	}()

	select {
	case <-vm.AbortChan:
		return nil, vvm.ErrVMAborted
	case <-done:
	}
	return
}

func timesParseDuration(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	s1, ok := vvm.ToString(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	dur, err := time.ParseDuration(s1)
	if err != nil {
		ret = wrapError(err)
		return
	}

	ret = &vvm.Int{Value: int64(dur)}

	return
}

func timesSince(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: int64(time.Since(t1))}

	return
}

func timesUntil(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: int64(time.Until(t1))}

	return
}

func timesDurationHours(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	i1, ok := vvm.ToInt64(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Float{Value: time.Duration(i1).Hours()}

	return
}

func timesDurationMinutes(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	i1, ok := vvm.ToInt64(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Float{Value: time.Duration(i1).Minutes()}

	return
}

func timesDurationNanoseconds(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	i1, ok := vvm.ToInt64(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: time.Duration(i1).Nanoseconds()}

	return
}

func timesDurationSeconds(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	i1, ok := vvm.ToInt64(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Float{Value: time.Duration(i1).Seconds()}

	return
}

func timesDurationString(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	i1, ok := vvm.ToInt64(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.String{Value: time.Duration(i1).String()}

	return
}

func timesMonthString(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	i1, ok := vvm.ToInt64(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.String{Value: time.Month(i1).String()}

	return
}

func timesDate(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 7 {
		err = vvm.ErrWrongNumArguments
		return
	}

	i1, ok := vvm.ToInt(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}
	i2, ok := vvm.ToInt(args[1])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}
	i3, ok := vvm.ToInt(args[2])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "third",
			Expected: "int(compatible)",
			Found:    args[2].TypeName(),
		}
		return
	}
	i4, ok := vvm.ToInt(args[3])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "fourth",
			Expected: "int(compatible)",
			Found:    args[3].TypeName(),
		}
		return
	}
	i5, ok := vvm.ToInt(args[4])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "fifth",
			Expected: "int(compatible)",
			Found:    args[4].TypeName(),
		}
		return
	}
	i6, ok := vvm.ToInt(args[5])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "sixth",
			Expected: "int(compatible)",
			Found:    args[5].TypeName(),
		}
		return
	}
	i7, ok := vvm.ToInt(args[6])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "seventh",
			Expected: "int(compatible)",
			Found:    args[6].TypeName(),
		}
		return
	}

	ret = &vvm.Time{
		Value: time.Date(i1,
			time.Month(i2), i3, i4, i5, i6, i7, time.Now().Location()),
	}

	return
}

func timesNow(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 0 {
		err = vvm.ErrWrongNumArguments
		return
	}

	ret = &vvm.Time{Value: time.Now()}

	return
}

func timesParse(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 2 {
		err = vvm.ErrWrongNumArguments
		return
	}

	s1, ok := vvm.ToString(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	s2, ok := vvm.ToString(args[1])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	parsed, err := time.Parse(s1, s2)
	if err != nil {
		ret = wrapError(err)
		return
	}

	ret = &vvm.Time{Value: parsed}

	return
}

func timesUnix(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 2 {
		err = vvm.ErrWrongNumArguments
		return
	}

	i1, ok := vvm.ToInt64(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	i2, ok := vvm.ToInt64(args[1])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	ret = &vvm.Time{Value: time.Unix(i1, i2)}

	return
}

func timesAdd(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 2 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	i2, ok := vvm.ToInt64(args[1])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	ret = &vvm.Time{Value: t1.Add(time.Duration(i2))}

	return
}

func timesSub(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 2 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	t2, ok := vvm.ToTime(args[1])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "time(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: int64(t1.Sub(t2))}

	return
}

func timesAddDate(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 4 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	i2, ok := vvm.ToInt(args[1])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	i3, ok := vvm.ToInt(args[2])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "third",
			Expected: "int(compatible)",
			Found:    args[2].TypeName(),
		}
		return
	}

	i4, ok := vvm.ToInt(args[3])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "fourth",
			Expected: "int(compatible)",
			Found:    args[3].TypeName(),
		}
		return
	}

	ret = &vvm.Time{Value: t1.AddDate(i2, i3, i4)}

	return
}

func timesAfter(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 2 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	t2, ok := vvm.ToTime(args[1])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "time(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	if t1.After(t2) {
		ret = vvm.TrueValue
	} else {
		ret = vvm.FalseValue
	}

	return
}

func timesBefore(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 2 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	t2, ok := vvm.ToTime(args[1])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	if t1.Before(t2) {
		ret = vvm.TrueValue
	} else {
		ret = vvm.FalseValue
	}

	return
}

func timesTimeYear(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: int64(t1.Year())}

	return
}

func timesTimeMonth(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: int64(t1.Month())}

	return
}

func timesTimeDay(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: int64(t1.Day())}

	return
}

func timesTimeWeekday(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: int64(t1.Weekday())}

	return
}

func timesTimeHour(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: int64(t1.Hour())}

	return
}

func timesTimeMinute(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: int64(t1.Minute())}

	return
}

func timesTimeSecond(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: int64(t1.Second())}

	return
}

func timesTimeNanosecond(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: int64(t1.Nanosecond())}

	return
}

func timesTimeUnix(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: t1.Unix()}

	return
}

func timesTimeUnixNano(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Int{Value: t1.UnixNano()}

	return
}

func timesTimeFormat(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 2 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	s2, ok := vvm.ToString(args[1])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	s := t1.Format(s2)
	if len(s) > vvm.MaxStringLen {

		return nil, vvm.ErrStringLimit
	}

	ret = &vvm.String{Value: s}

	return
}

func timesIsZero(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	if t1.IsZero() {
		ret = vvm.TrueValue
	} else {
		ret = vvm.FalseValue
	}

	return
}

func timesToLocal(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Time{Value: t1.Local()}

	return
}

func timesToUTC(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.Time{Value: t1.UTC()}

	return
}

func timesTimeLocation(args ...vvm.Object) (
	ret vvm.Object,
	err error,
) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.String{Value: t1.Location().String()}

	return
}

func timesTimeString(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		err = vvm.ErrWrongNumArguments
		return
	}

	t1, ok := vvm.ToTime(args[0])
	if !ok {
		err = vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &vvm.String{Value: t1.String()}

	return
}

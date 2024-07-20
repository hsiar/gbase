package gbase

import "github.com/golang-module/carbon"

type Time struct {
	Cb carbon.Carbon
	//ts int64
}

func (this *Time) Timestamp() int64 {
	return this.Cb.Timestamp()
}

func (this *Time) WithTime(time int64) *Time {
	this.Cb = this.Cb.CreateFromTimestamp(time)
	return this
}

func (this *Time) SubMonth() *Time {
	this.Cb = this.Cb.SubMonth()
	return this
}

func (this *Time) DayStageCn() string {
	h := this.Cb.Hour()
	switch {
	case h < 12:
		return "早上"
	case h < 18:
		return "下午"
	default:
		return "晚上"
	}
}

func (this *Time) DayFirstTime() int64 {
	return this.Cb.StartOfDay().Timestamp()
}

func (this *Time) DayEndTime() int64 {
	return this.Cb.EndOfDay().Timestamp()
}

func (this *Time) MonthFirstTime() int64 {
	return this.Cb.StartOfMonth().Timestamp()
}

func (this *Time) MonthEndTime() int64 {
	return this.Cb.EndOfMonth().Timestamp()
}

func (this *Time) YearFirstTime() int64 {
	return this.Cb.StartOfYear().Timestamp()
}

func (this *Time) YearEndTime() int64 {
	return this.Cb.EndOfYear().Timestamp()
}

func NewTime(ts ...int64) (t *Time) {
	t = &Time{}
	t.Cb = carbon.Now()

	if len(ts) > 0 {

		t.Cb = t.Cb.CreateFromTimestamp(ts[0])
	}
	return
}

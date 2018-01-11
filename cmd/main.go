package main

import (
	"github.com/bxcodec/faker"
	"github.com/johnmcdnl/compare"
)

type SomeStruct struct {
	Latitude         float32
	Long             float32
	CreditCardType   string
	CreditCardNumber string
	Email            string
	IPV4             string
	IPV6             string
	Password         string
	PhoneNumber      string
	MacAddress       string
	Url              string
	UserName         string
	ToolFreeNumber   string
	E164PhoneNumber  string
	TitleMale        string
	TitleFemale      string
	FirstNameMale    string
	FirstNameFemale  string
	LastName         string
	Name             string
	UnixTime         int64
	Date             string
	Time             string
	MonthName        string
	Year             string
	DayOfWeek        string
	DayOfMonth       string
	Timestamp        string
	Century          string
	TimeZone         string
	TimePeriod       string
	Word             string
	Sentence         string
	Sentences        string
}

func main() {
	var t1 SomeStruct
	var t2 SomeStruct
	faker.FakeData(&t1)
	t2 = t1

	compare.Struct(t1, t2, nil, nil )
}

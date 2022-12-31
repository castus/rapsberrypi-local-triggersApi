package lightController_test

import (
	"testing"
	"time"

	light "raspberrypi.local/triggersApi/lightController"
	sunAPI "raspberrypi.local/triggersApi/sunPositionAPI"
)

var sunPositionAPI = sunAPI.SunPosition{
	Sunrise: getParsedDate("Sun Dec 4 07:25:13 CET 2022"),
	Sunset:  getParsedDate("Sun Dec 4 15:27:25 CET 2022"),
}

var location = getLocation()
var sut = light.NewLightController(sunPositionAPI)

func TestTriggerAtNight(t *testing.T) {
	if getTriggerKey(getParsedDate("Sat Dec 3 23:59:00 CET 2022")) != "trigger-1" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 02:59:00 CET 2022")) != "trigger-1" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 01:00:00 CET 2022")) != "trigger-1" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 00:00:00 CET 2022")) != "trigger-1" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 23:01:00 CET 2022")) != "trigger-1" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 23:59:00 CET 2022")) != "trigger-1" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Mon Dec 5 00:01:00 CET 2022")) != "trigger-1" {
		t.Fail()
	}
}

func TestTriggerAtNoon(t *testing.T) {
	if getTriggerKey(getParsedDate("Sun Dec 4 11:00:00 CET 2022")) != "trigger-5" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 12:59:00 CET 2022")) != "trigger-5" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 13:00:00 CET 2022")) != "trigger-5" {
		t.Fail()
	}
}

func TestTriggerRightBeforeAndAfterNoon(t *testing.T) {
	if getTriggerKey(getParsedDate("Sun Dec 4 08:30:00 CET 2022")) != "trigger-4" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 10:59:59 CET 2022")) != "trigger-4" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 13:00:01 CET 2022")) != "trigger-4" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 14:27:01 CET 2022")) != "trigger-4" {
		t.Fail()
	}
}

func TestTriggerSunriseToSunrisePlusOne(t *testing.T) {
	if getTriggerKey(getParsedDate("Sun Dec 4 07:26:00 CET 2022")) != "trigger-3" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 08:25:00 CET 2022")) != "trigger-3" {
		t.Fail()
	}
}

func TestTriggerSunsetMinusOneToSunset(t *testing.T) {
	if getTriggerKey(getParsedDate("Sun Dec 4 14:28:00 CET 2022")) != "trigger-3" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 15:27:00 CET 2022")) != "trigger-3" {
		t.Fail()
	}
}

func TestTriggerEarlyMorning(t *testing.T) {
	if getTriggerKey(getParsedDate("Sun Dec 4 03:01:00 CET 2022")) != "trigger-2" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 07:25:00 CET 2022")) != "trigger-2" {
		t.Fail()
	}
}

func TestTriggerEvening(t *testing.T) {
	if getTriggerKey(getParsedDate("Sun Dec 4 15:28:00 CET 2022")) != "trigger-6" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 22:00:00 CET 2022")) != "trigger-6" {
		t.Fail()
	}
}

func TestTriggerLateEvening(t *testing.T) {
	if getTriggerKey(getParsedDate("Sun Dec 4 22:00:01 CET 2022")) != "trigger-7" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("Sun Dec 4 23:00:00 CET 2022")) != "trigger-7" {
		t.Fail()
	}
}

func getParsedDate(date string) time.Time {
	d, _ := time.Parse(time.UnixDate, date)
	return d
}

func getTriggerKey(date time.Time) string {
	return sut.GetTriggerKey(location, date).Current
}

func getLocation() *time.Location {
	l, _ := time.LoadLocation("Europe/Warsaw")
	return l
}

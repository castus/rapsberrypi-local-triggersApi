package lightController_test

import (
	"testing"
	"time"

	light "raspberrypi.local/triggersApi/lightController"
	sunAPI "raspberrypi.local/triggersApi/sunPositionAPI"
)

// Important: Dates are in UTC
// So noon in life is 2022-12-04T11:00:00 in UTC
var sunPositionAPI = sunAPI.SunPosition{
	Sunrise: getParsedDate("2022-12-04T06:25:13"),
	Sunset:  getParsedDate("2022-12-04T14:27:25"),
}

var location = getLocation()
var sut = light.NewLightController(sunPositionAPI)

func TestTriggerAtNoon(t *testing.T) {
	if getTriggerKey(getParsedDate("2022-12-04T10:25:00")) != "trigger-5" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T10:00:00")) != "trigger-5" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T11:00:00")) != "trigger-5" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T11:00:00")) != "trigger-5" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T12:00:00")) != "trigger-5" {
		t.Fail()
	}
}

func TestTriggerRightBeforeAndAfterNoon(t *testing.T) {
	if getTriggerKey(getParsedDate("2022-12-04T07:30:00")) != "trigger-4" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T09:59:59")) != "trigger-4" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T12:00:01")) != "trigger-4" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T13:27:01")) != "trigger-4" {
		t.Fail()
	}
}

func TestTriggerSunriseToSunrisePlusOne(t *testing.T) {
	if getTriggerKey(getParsedDate("2022-12-04T06:26:00")) != "trigger-3" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T07:25:00")) != "trigger-3" {
		t.Fail()
	}
}

func TestTriggerSunsetMinusOneToSunset(t *testing.T) {
	if getTriggerKey(getParsedDate("2022-12-04T13:28:00")) != "trigger-3" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T14:27:00")) != "trigger-3" {
		t.Fail()
	}
}

func TestTriggerEarlyMorning(t *testing.T) {
	if getTriggerKey(getParsedDate("2022-12-04T02:01:00")) != "trigger-2" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T06:25:00")) != "trigger-2" {
		t.Fail()
	}
}

func TestTriggerEvening(t *testing.T) {
	if getTriggerKey(getParsedDate("2022-12-04T14:28:00")) != "trigger-6" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T21:00:00")) != "trigger-6" {
		t.Fail()
	}
}

func TestTriggerLateEvening(t *testing.T) {
	if getTriggerKey(getParsedDate("2022-12-04T21:00:01")) != "trigger-7" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T22:00:00")) != "trigger-7" {
		t.Fail()
	}
}

func TestTriggerNightSameDay(t *testing.T) {
	if getTriggerKey(getParsedDate("2022-12-04T23:00:01")) != "trigger-1" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T00:00:00")) != "trigger-1" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T01:00:01")) != "trigger-1" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-04T22:00:01")) != "trigger-1" {
		t.Fail()
	}
}

func TestTriggerNightTheNextDay(t *testing.T) {
	if getTriggerKey(getParsedDate("2022-12-05T23:00:01")) != "trigger-1" {
		t.Fail()
	}
	if getTriggerKey(getParsedDate("2022-12-05T00:00:01")) != "trigger-1" {
		t.Fail()
	}
}

func getParsedDate(date string) time.Time {
	d, _ := time.Parse("2006-01-02T15:04:05", date)
	return d
}

func getTriggerKey(date time.Time) string {
	return sut.GetTriggerKey(location, date).Current
}

func getLocation() *time.Location {
	l, _ := time.LoadLocation("Europe/Warsaw")
	return l
}

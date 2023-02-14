package lightController

import (
	"fmt"
	"log"
	"time"

	"github.com/uniplaces/carbon"

	api "raspberrypi.local/triggersApi/sunPositionAPI"
)

/*
Alexa Color definitions:
+--------------------------+-----------------+
|          Color           | Temperature (K) |
+--------------------------+-----------------+
| warm, warm white         |            2200 |
| incandescent, soft white |            2700 |
| white                    |            4000 |
| daylight, daylight white |            5500 |
| cool, cool white         |            7000 |
+--------------------------+-----------------+

Triggers definitions:
+-----------+-----------+---------+-------+------------+
| Monkey id | From (HH) | To (HH) | Color | Brightness |
+-----------+-----------+---------+-------+------------+
| trigger-1 | 23:00     | 03:00   |  2700 |         40 |
| trigger-2 | 03:00     | SR      |  2200 |          0 |
| trigger-3 | SR        | SR+1    |  4000 |         80 |
| trigger-3 | SS-1      | SS      |  4000 |         80 |
| trigger-4 | SR+1      | NOON-1  |  5500 |        100 |
| trigger-4 | NOON+1    | SS-1    |  5500 |        100 |
| trigger-5 | NOON-1    | NOON+1  |  7000 |        100 |
| trigger-6 | SS        | 22:00   |  2700 |         80 |
| trigger-7 | 22:00     | 23:00   |  2700 |         60 |
| trigger-500 | red light for debug purposes           |
+-----------+-----------+---------+-------+------------+
Legend:
* SR - Sunrise
* SS - Sunset

Timeline:
00:00     03:00         SR        SR+1        NOON-1      NOON+1      SS-1         SS         22:00       23:00       23:59
  +---------+-----------------------------------------------------------------------------------------------------------+
  |trigger-1| trigger-2 | trigger-3 | trigger-4 | trigger-5 | trigger-4 | trigger-3 | trigger-6 | trigger-7 | trigger-1 |
  +---------+-----------------------------------------------------------------------------------------------------------+
*/

const (
	trigger1   = "trigger-1"
	trigger2   = "trigger-2"
	trigger3   = "trigger-3"
	trigger4   = "trigger-4"
	trigger5   = "trigger-5"
	trigger6   = "trigger-6"
	trigger7   = "trigger-7"
	trigger500 = "trigger-500"
)

type APIResponse struct {
	Current     string   `json:"current"`
	Checkpoints []string `json:"checkpoints"`
}

type LightController struct {
	API api.SunPosition
}

func NewLightController(api api.SunPosition) *LightController {
	return &LightController{
		API: api,
	}
}

func (l *LightController) GetTriggerKey(loc *time.Location, now time.Time) APIResponse {
	sunrise := carbon.NewCarbon(l.API.Sunrise.In(loc))
	sunrisePlusOne := carbon.NewCarbon(l.API.Sunrise.Add(time.Hour * 1).In(loc))
	sunset := carbon.NewCarbon(l.API.Sunset.In(loc))
	sunsetMinusOne := carbon.NewCarbon(l.API.Sunset.Add(time.Hour * -1).In(loc))
	nowCarbon := carbon.NewCarbon(now.In(loc))
	noonPlusOne := carbon.NewCarbon(time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, loc).Add(time.Hour * 1))
	noonMinusOne := carbon.NewCarbon(time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, loc).Add(time.Hour * -1))
	evening := carbon.NewCarbon(time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, loc))
	lateEvening := carbon.NewCarbon(time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, loc))
	night := carbon.NewCarbon(time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, loc))

	var checkpoints []string
	checkpoints = append(checkpoints,
		sunrise.Format(time.UnixDate),
		sunrisePlusOne.Format(time.UnixDate),
		noonMinusOne.Format(time.UnixDate),
		noonPlusOne.Format(time.UnixDate),
		sunsetMinusOne.Format(time.UnixDate),
		sunset.Format(time.UnixDate),
		evening.Format(time.UnixDate),
		lateEvening.Format(time.UnixDate),
		night.Format(time.UnixDate),
	)

	var trigger string

	if !nowCarbon.IsSameDay(sunset) {
		trigger = trigger1
	} else if nowCarbon.Between(noonMinusOne, noonPlusOne, true) {
		trigger = trigger5
	} else if nowCarbon.Between(sunrisePlusOne, noonMinusOne, true) || nowCarbon.Between(noonPlusOne, sunsetMinusOne, true) {
		trigger = trigger4
	} else if nowCarbon.Between(sunrise, sunrisePlusOne, true) || nowCarbon.Between(sunsetMinusOne, sunset, true) {
		trigger = trigger3
	} else if nowCarbon.Between(night, sunrise, true) {
		trigger = trigger2
	} else if nowCarbon.Between(sunset, evening, true) {
		trigger = trigger6
	} else if nowCarbon.Between(evening, lateEvening, true) {
		trigger = trigger7
	} else if nowCarbon.Lte(night) || nowCarbon.Gte(lateEvening) {
		trigger = trigger1
	} else {
		trigger = trigger500
	}

	log.Printf("type=debug msg=\"For the given time %s, trigger is %s\"\n", now.String(), trigger)
	var checkpointString = ""
	for index, element := range checkpoints {
		checkpointString = fmt.Sprintf("%s check-point-%d=\"%s\"", checkpointString, index+1, element)
	}
	log.Printf("type=debug tag=checkpoints %s\n", checkpointString)

	return APIResponse{
		trigger, checkpoints,
	}
}

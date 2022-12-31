package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"

	"raspberrypi.local/triggersApi/lightController"
	"raspberrypi.local/triggersApi/sunPositionAPI"
)

var sun = sunPositionAPI.NewSunPositionAPI()

func main() {
	automaticallyRefreshDataWhenDayStarts()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		l := lightController.NewLightController(sun.Get())
		loc, _ := time.LoadLocation("Europe/Warsaw")
		now := time.Now().In(loc)

		b, err := json.Marshal(l.GetTriggerKey(loc, now))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		} else {
			w.WriteHeader(http.StatusOK)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = fmt.Fprintf(w, string(b))
		if err != nil {
			panic(err)
		}
	})

	port := ":8080"
	fmt.Println("Triggers Api server is running on port" + port)

	log.Fatal(http.ListenAndServe(port, nil))
}

func automaticallyRefreshDataWhenDayStarts() {
	c := cron.New()
	_, err := c.AddFunc("0 2 * * *", func() { // At 2 AM
		fmt.Println("Cron function executes, refreshing data when day starts")
		s := sunPositionAPI.NewSunPositionAPI()
		s.Get()
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}

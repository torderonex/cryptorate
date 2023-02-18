package timemanager

import (
	"strconv"
	"strings"
	"time"
)

var (
	location = time.Now().Location()
)

func WaitUntil(tm string, foo func(), checkOk func() bool) {
	for true {
		time.Sleep(time.Duration(getWaitingTime(tm)+1) * time.Second)
		if !checkOk() {
			break
		}
		foo()
	}

}

func getWaitingTime(tm string) int {
	tmarr := strings.Split(tm, ":")
	hours, _ := strconv.Atoi(tmarr[0])
	minutes, _ := strconv.Atoi(tmarr[1])
	temp := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hours, minutes, 0, 0, location)
	if time.Now().Hour() > temp.Hour() || (time.Now().Hour() == temp.Hour() && time.Now().Minute() >= temp.Minute()) {
		temp = temp.AddDate(0, 0, 1)

	}
	return int(temp.Sub(time.Now()).Seconds())
}
 
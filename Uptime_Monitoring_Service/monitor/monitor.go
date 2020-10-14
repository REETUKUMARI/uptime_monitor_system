package monitor

import (
	"fmt"
	"net/http"
	"time"
)

func Checklink() {
	rows, err := MonitorRepo.GetRows()
	if err != nil {
		fmt.Println("connection failed")
	}

	defer rows.Close()

	var (
		id              string
		urllink         string
		crawltimeout    time.Duration
		frequency       int
		failurethresold int
		status          string
		failurecount    int
	)

	c := make(chan string)
	for rows.Next() {
		rows.Scan(&id, &urllink, &crawltimeout, &frequency, &failurethresold, &status, &failurecount)
		fmt.Println(urllink, crawltimeout, frequency)
		go checkLink(urllink, c)
	}

	for l := range c {
		go func(link string) {
			time.Sleep(1 * time.Second)
			if link != "" {
				checkLink(link, c)
			}
		}(l)
	}
}

func checkLink(urllink string, c chan string) {
	p := MonitorRepo.GetUrlData(urllink)
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	res, err := client.Get(urllink)
	if err != nil || res.StatusCode != http.StatusOK {
		MonitorRepo.IncreaseFailureCount(p.ID)
		if p.FailureCount >= p.FailureThreshold {
			MonitorRepo.FailureCountToZero(p.ID)
			c <- ""
		}
		MonitorRepo.UpdateStatus(p.ID, "inactive")
	} else {
		if p.Status != "active" {
			MonitorRepo.UpdateStatus(p.ID, "active")
			MonitorRepo.FailureCountToZero(p.ID)
		}
	}
	c <- urllink
}
func Checkurl(url string, crawltimeout time.Duration) string {
	client := http.Client{
		Timeout: crawltimeout * time.Second,
	}
	res, err := client.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		return "inactive"
	} else {
		return "active"
	}
}

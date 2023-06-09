package main

import (
	"strconv"
	"strings"
	"time"

	"tinygo.org/x/drivers/net"
)

type nordnet struct {
	server string
	client HttpClient
}

func newNordnet() *nordnet {
	return &nordnet{
		server: "https://api.prod.nntech.io",
		client: HttpClient{
			timeout:     time.Second,
			connections: map[string]net.Conn{},
		},
	}
}

func (b *nordnet) getReturns(periods []string, instrumentId int) ([]float64, error) {
	returns := make([]float64, len(periods))
	url := b.server + "/market-data/price-time-series/v2/returns/" + strconv.Itoa(instrumentId)
	println(url)
	req := newGET(url, nil)
	res, err := b.client.sendHttp(req, false)
	if err != nil {
		return returns, err
	} else {
		trace(string(res.bytes))
	}

	data := string(res.bytes)
	for i, period := range periods {
		head := `"` + period + `","development":`
		start := strings.Index(data, head) + len(head)
		end := strings.Index(data[start:], `,`) + start
		returns[i], _ = strconv.ParseFloat(data[start:end], 32)
	}
	return returns, nil
}

func (b *nordnet) getLast(instrumentId int) (float64, error) {
	url := b.server + "/market-data/price-time-series/v2/returns/" + strconv.Itoa(instrumentId)
	println(url)
	req := newGET(url, nil)
	res, err := b.client.sendHttp(req, false)
	if err != nil {
		return 0, err
	} else {
		trace(string(res.bytes))
	}

	data := string(res.bytes)
	day1 := `"DAY_1"`
	day1_start := strings.Index(data, day1)
	absdev := `absoluteDevelopment":`
	absolutedevelopmentStart := strings.Index(data[day1_start:], absdev) + len(absdev) + day1_start

	absolutedevelopmentEnd := strings.Index(data[absolutedevelopmentStart:], `,`) + absolutedevelopmentStart
	todayChange, _ := strconv.ParseFloat(data[absolutedevelopmentStart:absolutedevelopmentEnd], 32)

	closeText := `close":`

	closeStart := strings.Index(data[day1_start:], closeText) + len(closeText) + day1_start
	closeEnd := strings.Index(data[closeStart:], `,`) + closeStart
	close, _ := strconv.ParseFloat(data[closeStart:closeEnd], 32)

	return close + todayChange, nil

}

func extractDevelopment(period string, data string) (float64, error) {
	head := `"` + period + `","development":`
	start := strings.Index(data, head) + len(head)
	end := strings.Index(data[start:], `,`) + start
	return strconv.ParseFloat(data[start:end], 32)
}

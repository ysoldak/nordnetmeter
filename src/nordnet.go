package main

import (
	"strconv"
	"strings"
	"time"

	"tinygo.org/x/drivers/net"
)

type nordnet struct {
	server string
	// token  string
	client HttpClient
}

func newNordnet() *nordnet {
	return &nordnet{
		server: "https://api.prod.nntech.io",
		// token:  blynkToken,
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

func extractDevelopment(period string, data string) (float64, error) {
	head := `"` + period + `","development":`
	start := strings.Index(data, head) + len(head)
	end := strings.Index(data[start:], `,`) + start
	return strconv.ParseFloat(data[start:end], 32)
}

/*
func (b *blynk) updateInt(name string, value int) (err error) {
	if b.token == "" {
		return
	}
	url := b.server + "/external/api/update?token=" + b.token + "&" + name + "=" + strconv.Itoa(value)
	req := newGET(url, nil)
	res, err := b.client.sendHttp(req, false)
	if err != nil {
		return err
	} else {
		trace(string(res.bytes))
	}
	return nil
}

func (b *blynk) updateFloat(name string, value float64) (err error) {
	if b.token == "" {
		return
	}
	url := b.server + "/external/api/update?token=" + b.token + "&" + name + "=" + strconv.FormatFloat(value, 'f', 2, 64)
	req := newGET(url, nil)
	res, err := b.client.sendHttp(req, false)
	if err != nil {
		return err
	} else {
		trace(string(res.bytes))
	}
	return nil
}

func (b *blynk) sendEvent(name string) (err error) {
	if b.token == "" {
		return
	}
	url := b.server + "/external/api/logEvent?token=" + b.token + "&code=" + name
	req := newGET(url, nil)
	res, err := b.client.sendHttp(req, false)
	if err != nil {
		return err
	} else {
		trace(string(res.bytes))
	}
	return nil
}
*/

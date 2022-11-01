package stream

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type Reply struct {
	Msg string
}

func HeartBeats(token string, streamname string, edgename string, namespace string, streamcode string) Reply {
	reply := Reply{}
	client := resty.New()
	resp, err := client.R().
		SetResult(&reply).
		SetAuthToken(token).
		SetQueryParams(map[string]string{
			"streamname": streamname + "." + namespace,
			"edgename":   edgename,
			"streamcode": streamcode,
		}).
		Post("http://192.168.20.151:30085/api/v1/stream/heartbeat")

	if err != nil {
		fmt.Println("resp:", resp)
		fmt.Println("error on stream heart beats ", err.Error())
	}
	return reply
}

func Sendheartbeats(token string, streamname string, edgename string, namespace string, streamcode string, hbsignal chan string) {
	i := 1
	myTimer := time.NewTimer(2000 * time.Millisecond)
	for {
		select {
		case v := <-hbsignal:
			if v == "stop" {
				fmt.Println("stop hearbeats to " + streamcode)
				myTimer.Stop()
				return
			}
		case <-myTimer.C:
			reply := HeartBeats(token, streamname, edgename, namespace, streamcode)
			fmt.Println("msg"+strconv.Itoa(i)+"", reply.Msg)
			myTimer.Reset(2000 * time.Millisecond)
			i++
		}
	}
}

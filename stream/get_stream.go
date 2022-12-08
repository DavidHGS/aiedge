package stream

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Stream struct {
	Hls        string
	Httpflv    string
	Rtmp       string
	Streamcode string
}

func GetStream(token string, streamname string, edgename string, namespace string) Stream {
	stream := Stream{}
	client := resty.New()
	resp, err := client.R().
		SetResult(&stream).
		SetAuthToken(token).
		SetQueryParams(map[string]string{
			"streamname": streamname + "." + namespace,
			"edgename":   edgename,
		}).
		Get("http://192.168.20.151:30085/api/v1/stream")

	if err != nil {
		fmt.Println("resp:", resp)
		fmt.Println("error on get stream address ", err.Error())
	}
	return stream
}

package post

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-resty/resty/v2"
)

type SigninBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type JwtToken struct { //获取token
	Msg string
}
type Rectangle struct {
	Bottom float64
	Left   float64
	Right  float64
	Top    float64
}
type DetectionObject struct {
	Label int8
	Rect  Rectangle
	Score float64
}
type DetectionResult struct {
	Rects []*DetectionObject
}

func Signin(username, password string) JwtToken {
	jwt := JwtToken{}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(SigninBody{
			Username: username,
			Password: password,
		}).
		SetResult(&jwt).
		Post("http://1.117.224.103:6100/api/v1/auth/signin")

	if err != nil {
		fmt.Println("resp:", resp)
		fmt.Println("error on sign in ", err.Error())
	}
	return jwt
}

func PostDetection(filePath, token string) DetectionResult {
	file, _ := os.Open(filePath)
	fileBytes, _ := ioutil.ReadAll(file)

	detectionResult := DetectionResult{}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "image/jpeg").
		SetBody(fileBytes).
		SetAuthToken(token).
		SetResult(&detectionResult).
		// SetAuthToken("Bearer " + token).
		// Post("http://172.26.70.122:30001/v1/object-detection")//zerotier Ip
		Post("http://192.168.1.200:30001/v1/object-detection")
	if err != nil {
		fmt.Println("resp:", resp)
		fmt.Println("error on PostDetection ", err.Error())
	}
	return detectionResult
}

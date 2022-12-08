package post

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"strconv"

	"github.com/go-resty/resty/v2"
	"gocv.io/x/gocv"
)

type Trajectory struct {
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
	Top    float64 `json:"top"`
	Id     int64   `json:"id"`
	Score  float64 `json:"score"`
}

type TrackingResult struct {
	TrajectoryMissonID int64        `json:"trajectory_misson_id"`
	Trajectorys        []Trajectory `json:"trajectorys"`
}

type CreateMissonResult struct {
	TrajectoryMissonID int64 `json:"trajectory_misson_id"`
}

func CreateMission(post_url, auth string) int64 {
	// post_url = "http://192.168.1.84:30002/v1/tracking/multiple-object-tracking/create"
	client := resty.New()
	client.SetCloseConnection(true)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", auth).
		SetQueryParam("video_frame_rate", "15").
		SetQueryParam("track_buffer", "30").
		Get(post_url)

	if err != nil {
		fmt.Println(resp)
		fmt.Println("can't create mot mission ", err.Error())
		return -1
	}
	create_misson_result := CreateMissonResult{}
	create_misson_result.TrajectoryMissonID = -1

	json.Unmarshal([]byte(resp.String()), &create_misson_result)
	return create_misson_result.TrajectoryMissonID
}

func RemoveMission(post_url, auth string, missionID int64) {
	// post_url = "http://192.168.1.84:30002/v1/tracking/multiple-object-tracking/remove"
	client := resty.New()
	client.SetCloseConnection(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", auth).
		SetQueryParam("trajectory_misson_id", strconv.FormatInt(missionID, 10)).
		Get(post_url)

	if err != nil {
		fmt.Println(resp)
		fmt.Println("can't remove mission ", err.Error())
		return
	}
}
func main() {
	auth := "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inh1c2hlbmciLCJzdGF0ZW1lbnRzIjpbeyJlZmZlY3QiOiJhbGxvdyIsImFjdGlvbiI6WyIqIl0sInJlc291cmNlIjpbIioiXX1dLCJleHAiOjE2Njg2ODI3NDAsImlhdCI6MTY2ODA3Nzk0MCwiaXNzIjoiYWllZGdlLWF1dGgifQ.WZFf4x3SYD1mBms17MTWAepJpjWXkuceoJ8Hd8EoTK6zhg8_tt7G7vg6QmDqrgcqd6ZALli_HRwfHWlihaa3mTrOmoXEZ30Yp1NU6ez_BHGfCOiS4K-M2j6UwPZ5MdIMGIfg0miEq8pjpJGrPbNleSSAT_JktmbXFhkMpyPjI4jAVRa0OWIUl-vi5YHzY5G9seLzhW2rUW6yD0mi2o0YSK5ufCnZdOvg17JIevKXbGUetfSfjsxkVvQNW3ggyzX5AqYttAFcA9J6Q5tH2KDKJlhOqGIP1W6cqkJddtOpRxOgCAmvjzxjq89S3Fga33vkZgKV_HpURQdpHKTvlqALcQ"

	// Create a Resty Client
	client := resty.New()
	client.SetCloseConnection(true)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", auth).
		SetQueryParam("video_frame_rate", "15").
		SetQueryParam("track_buffer", "30").
		Get("http://192.168.1.84:30002/v1/tracking/multiple-object-tracking/create")

	if err != nil {
		fmt.Println("can't create mot mission ", err.Error())
		return
	}

	create_misson_result := CreateMissonResult{}
	json.Unmarshal([]byte(resp.String()), &create_misson_result)

	for pic_i := 0; pic_i <= 5272; pic_i++ {

		image_file_path := fmt.Sprintf("/home/liujiaxi/Cam1/rgb_%05d_1.jpg", pic_i)

		mat := gocv.IMRead(image_file_path, gocv.IMReadColor)

		if mat.Empty() {
			fmt.Println("can't run open image ", "image is empty")
			return
		}

		buf, err := gocv.IMEncode(".jpg", mat)
		if err != nil {
			fmt.Println("can't run encode image ", err.Error())
			return
		}

		// 检测
		resp, err = client.R().
			SetHeader("Content-Type", "image/jpeg").
			SetHeader("Authorization", auth).
			SetBody(buf.GetBytes()).
			Post("http://192.168.1.84:30001/v1/object-detection")

		if err != nil {
			fmt.Println("can't run detection ", err.Error())
			return
		}

		// label 0 人类 1 ... 根据 label 做筛选
		// 人脸检测 格式要跟 object-detection 统一
		detection_result := resp.Body()

		// 跟踪
		// 跟踪的输入是检测的输出
		// 跟踪的输出其实就是给每个检测框打了一个 ID
		resp, err = client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", auth).
			SetBody(detection_result).
			SetQueryParam("trajectory_misson_id", strconv.FormatInt(create_misson_result.TrajectoryMissonID, 10)).
			Post("http://192.168.1.84:30002/v1/tracking/multiple-object-tracking")

		if err != nil {
			fmt.Println("can't run tracking ", err.Error())
			return
		}

		fmt.Println(resp.String())

		tracking_result := TrackingResult{}
		json.Unmarshal([]byte(resp.String()), &tracking_result)

		for i := 0; i < len(tracking_result.Trajectorys); i++ {
			trajectory := tracking_result.Trajectorys[i]

			rect := image.Rectangle{Min: image.Point{int(trajectory.Left), int(trajectory.Bottom)},
				Max: image.Point{int(trajectory.Right), int(trajectory.Top)}}
			gocv.Rectangle(&mat, rect, color.RGBA{255, 255, 0, 0}, 1)

			text := fmt.Sprintf("%d : %f", trajectory.Id, trajectory.Score)
			location := image.Point{int(trajectory.Left), int(trajectory.Bottom)}

			gocv.PutText(&mat, text, location, gocv.FontHersheyComplexSmall, 0.5, color.RGBA{122, 0, 122, 0}, 1)
		}

		result_path := fmt.Sprintf("/home/liujiaxi/white_box/go_test/result/%d.jpg", pic_i)

		gocv.IMWrite(result_path, mat)
	}

	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", auth).
		SetQueryParam("trajectory_misson_id", strconv.FormatInt(create_misson_result.TrajectoryMissonID, 10)).
		Get("http://192.168.1.84:30002/v1/tracking/multiple-object-tracking/remove")

	if err != nil {
		fmt.Println("can't remove mission ", err.Error())
		return
	}
}

package post

type BlurResult struct {
	Blur_degree []float64 `json:"blur_degree"`
}

type HeadPoseResult struct {
	Pitch float64 `json:"pitch"`
	Roll  float64 `json:"roll"`
	Yaw   float64 `json:"yaw"`
}

// func test() {
// 	//模糊度
// 	resp, err = client.R().
// 		SetHeader("Content-Type", "image/jpeg").
// 		SetHeader("Authorization", auth).
// 		SetBody(buf_cutted.GetBytes()).
// 		Post("http://192.168.20.150:30004/v1/blur-detection")
// 	if err != nil {
// 		fmt.Println(resp)
// 		fmt.Println("can't run blur-detection ", err.Error())
// 		return
// 	}
// 	fmt.Println(resp)
// 	blur_result := post.BlurResult{}
// 	json.Unmarshal([]byte(resp.String()), &blur_result)

// 	//偏转
// 	resp, err = client.R().
// 		SetHeader("Content-Type", "image/jpeg").
// 		SetHeader("Authorization", auth).
// 		SetBody(buf_cutted.GetBytes()).
// 		Post("http://192.168.20.150:30005/v1/face/head-pose")
// 	if err != nil {
// 		fmt.Println(resp)
// 		fmt.Println("can't run blur-detection ", err.Error())
// 		return
// 	}
// 	fmt.Println(resp)
// 	head_post_result := []post.HeadPoseResult{}
// 	json.Unmarshal([]byte(resp.String()), &head_post_result)
// }

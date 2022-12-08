package algorithm

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

type ImgType struct {
	Imgfile  []byte
	Filename string
}
type BlurResultSort struct {
	Filename    []string
	Blur_degree []float64 `json:"blur_degree"`
}

type HeadPoseResultSort struct {
	Id        int
	Filename  string
	Pitch     float64 `json:"pitch"`
	Roll      float64 `json:"roll"`
	Yaw       float64 `json:"yaw"`
	Sort      string
	PoseScore float64
	Socre     float64
}

const w1, w2, w3 = 0.3, 0.4, 0.3 //w1+w2+w3=100%
const blur_w, pose_w = 0.5, 0.5  //w4+w5=100%
var weight = []float64{w1, w2, w3, blur_w, pose_w}

// @去重算法
// 综合评分：每项指标的排名，折算成得百分制得分，两项得分进行占比权重的加权平均得出总分。

// 总图片数为N, pitch：(越接近)  -14  roll:(越接近) -15  yaw：（越接近） 81三项排名分别为n1，n2，n3，每项权重分别为w1，w2，w3。w1+w2+w3=100%

// 姿态得分 = （N - n1 + 1） * （100 / N） * w1 +（N - n2 + 1） * （100 / N） * w2 + （N - n3 + 1） * （100 / N） * w3

// 总图片数为N,模糊得分（值越小  越清晰），偏转得分两项排名分别为n4，n5，每项权重分别为blur_w，blur_w。blur_w+pose_w=100%
// 总得分 = （N - n4 + 1）* （100 / N） * blur_w +（N - n5 + 1） * （100 / N） * pose_w

func Deduplication(imgMap map[string]([]ImgType), auth string) {
	var imgStr, size string
	client := resty.New()

	for _, imgArr := range imgMap {
		imgStr = ""
		size = ""
		blur_result_sort := BlurResultSort{}
		head_post_result_sort := []HeadPoseResultSort{}

		for k, value := range imgArr {
			strTmp := strconv.Itoa(len(value.Imgfile)) + ","
			size += strTmp
			imgStr += string(value.Imgfile)
			blur_result_sort.Filename = append(blur_result_sort.Filename, value.Filename)
			head_post_result_sort = append(head_post_result_sort, HeadPoseResultSort{
				Filename: value.Filename,
				Id:       k,
			})
			fmt.Println(value.Filename)
		}
		size = strings.TrimRight(size, ",")
		//模糊检测
		resp, err := client.R().
			SetHeader("Content-Type", "image/jpeg").
			SetHeader("Authorization", auth).
			SetQueryParam("pic_offsets", size).
			SetBody([]byte(imgStr)).
			Post("http://192.168.20.150:30004/v1/blur-detection")
		if err != nil {
			fmt.Println(resp)
			fmt.Println("can't run blur-detection ", err.Error())
			return
		}
		fmt.Println(resp)

		json.Unmarshal([]byte(resp.String()), &blur_result_sort)
		// s := NewSortIndexs([]float64{5, 8, 10, 2, 9, 6})
		//模糊排序
		s := NewSortIndexs(blur_result_sort.Blur_degree)
		sort.Sort(s)
		indexArr := s.indexs
		fmt.Println(indexArr)
		N := len(blur_result_sort.Blur_degree)
		blurScoreArr := make([]float64, N)
		for k, v := range indexArr {
			blurScoreArr[k] = float64(N-v) * 100 / float64(N) * float64(blur_w)
		}
		// blur_result.Blur_degree

		//头部姿态检测
		resp, err = client.R().
			SetHeader("Content-Type", "image/jpeg").
			SetHeader("Authorization", auth).
			SetQueryParam("pic_offsets", size).
			SetBody([]byte(imgStr)).
			Post("http://192.168.20.150:30005/v1/face/head-pose")
		if err != nil {
			fmt.Println(resp)
			fmt.Println("can't run blur-detection ", err.Error())
			return
		}
		fmt.Println(resp)
		json.Unmarshal([]byte(resp.String()), &head_post_result_sort)
		//v.Sort = pitch_roll_yaw
		//pitch排序
		sort.Stable(HeadPoseResultPitchDecrement(head_post_result_sort))
		for k, v := range head_post_result_sort {
			tmp := "" + strconv.Itoa(k)
			v.Sort += tmp
		}
		//roll排序
		sort.Stable(HeadPoseResultRollDecrement(head_post_result_sort))
		for k, v := range head_post_result_sort {
			tmp := "_" + strconv.Itoa(k)
			v.Sort += tmp
		}
		//yaw排序
		sort.Stable(HeadPoseResultYawDecrement(head_post_result_sort))
		for k, v := range head_post_result_sort {
			tmp := "_" + strconv.Itoa(k)
			v.Sort += tmp
		}
		//计算 headPoseScore
		for _, v := range head_post_result_sort {
			sortArr := strings.Split(v.Sort, "_")
			for i, j := range sortArr {
				tmp, _ := strconv.ParseFloat(j, 64)
				v.PoseScore += tmp * weight[i]
			}
		}
		//按照PostScore 排序
		//总成绩存在head_post_result_sort.Score中
		sort.Stable(HeadPoseResultPoseScoreDecrement(head_post_result_sort))
		for k, v := range head_post_result_sort {
			//总得分 = （N - n4 + 1）* （100 / N） * blur_w +（N - n5 + 1） * （100 / N） * pose_w
			//blurScoreArr[k] = float64(N-v) * 100 / float64(N) * float64(blur_w)
			v.Socre = float64(N-k)*100/float64(N)*float64(pose_w) + blurScoreArr[v.Id]
		}
		//按照Score 排序 最终结果
		sort.Stable(HeadPoseResultScoreDecrement(head_post_result_sort))
		//全部排序完成  发送到接口
		fmt.Println("test")
	}
}

//———————————————————————————————————————————排序———————————————————————————————————————————————————
type SortIndexs struct {
	sort.Float64Slice // 可以替换成其他实现了 sort.Interface
	indexs            []int
}

func (p *SortIndexs) Swap(i, j int) {
	p.Float64Slice.Swap(i, j)
	p.indexs[i], p.indexs[j] = p.indexs[j], p.indexs[i]
}

func NewSortIndexs(arr []float64) *SortIndexs {
	s := &SortIndexs{Float64Slice: sort.Float64Slice(arr), indexs: make([]int, len(arr))}
	for i := range s.indexs {
		s.indexs[i] = i // 原有排序 indexs
	}
	return s
}

//patch 排序函数
type HeadPoseResultPitchDecrement []HeadPoseResultSort

func (h HeadPoseResultPitchDecrement) Len() int { return len(h) }

func (h HeadPoseResultPitchDecrement) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h HeadPoseResultPitchDecrement) Less(i, j int) bool {
	return math.Abs(h[i].Pitch-(-14)) > math.Abs(h[j].Pitch-(-14))
}

// roll 排序函数
type HeadPoseResultRollDecrement []HeadPoseResultSort

func (h HeadPoseResultRollDecrement) Len() int { return len(h) }

func (h HeadPoseResultRollDecrement) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h HeadPoseResultRollDecrement) Less(i, j int) bool {
	return math.Abs(h[i].Roll-(-15)) > math.Abs(h[j].Roll-(-15))
}

// yaw 排序函数
type HeadPoseResultYawDecrement []HeadPoseResultSort

func (h HeadPoseResultYawDecrement) Len() int { return len(h) }

func (h HeadPoseResultYawDecrement) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h HeadPoseResultYawDecrement) Less(i, j int) bool {
	return math.Abs(h[i].Yaw-(81)) > math.Abs(h[j].Yaw-(81))
}

//headPoseScore 排序
type HeadPoseResultPoseScoreDecrement []HeadPoseResultSort

func (h HeadPoseResultPoseScoreDecrement) Len() int { return len(h) }

func (h HeadPoseResultPoseScoreDecrement) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h HeadPoseResultPoseScoreDecrement) Less(i, j int) bool {
	return h[i].PoseScore < h[j].PoseScore
}

//Score 排序
type HeadPoseResultScoreDecrement []HeadPoseResultSort

func (h HeadPoseResultScoreDecrement) Len() int { return len(h) }

func (h HeadPoseResultScoreDecrement) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h HeadPoseResultScoreDecrement) Less(i, j int) bool {
	return h[i].Socre < h[j].Socre
}

// id 排序函数
type HeadPoseResultIdDecrement []HeadPoseResultSort

func (h HeadPoseResultIdDecrement) Len() int { return len(h) }

func (h HeadPoseResultIdDecrement) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h HeadPoseResultIdDecrement) Less(i, j int) bool {
	return h[i].Id > h[j].Id
}

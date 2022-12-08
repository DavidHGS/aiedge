package ffmpeg

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"os/exec"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// /*
// 1.将视频抽帧成图片
// 2.将抽帧图片发送请求到接口
// 3.用返回的boundingbox给图片中画框
// 4.将画框图片拼接成视频
// */

// // func main() {
// // 	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
// // 	fmt.Println("dir", dir)                 //workDir
// // 	var ffmpegPath string = dir + "/ffmpeg" //ffmpeg地址
// // 	urlPath := os.Args[1]                   //下载视频的连接
// // 	// saveFilename := dir + "/vidio.mp4"
// // 	currTime := time.Now().UnixNano() //当前时间
// // 	//var freq float64 =60

// // 	//根据视频地址下载视频到本地

// // 	//todo
// // 	var flag bool = true
// // 	//创建images文件夹，存储图片
// // 	if flag {
// // 		path := "./images/"
// // 		if err := os.Mkdir(path, os.ModePerm); err != nil {
// // 			fmt.Println("mkdir failed")
// // 		}

// // 		//ffmpeg使用命令: ffmpeg -i http://video.pearvideo.com/head/20180301/cont-1288289-11630613.mp4 -r 1 -t 4 -f image2 image-%05d.jpeg
// // 		/*
// // 		   -t 代表持续时间，单位为秒
// // 		   -f 指定保存图片使用的格式，可忽略。
// // 		   -r 指定抽取的帧率，即从视频中每秒钟抽取图片的数量。1代表每秒抽取一帧。
// // 		   -ss 指定起始时间
// // 		   -vframes 指定抽取的帧数
// // 		*/
// // 		//获取视频长度
// // 		videoLen, _ := GenerateLength(ffmpegPath, urlPath, "1111111111")
// // 		testFfmpegParams(urlPath, path, ffmpegPath, 60, videoLen)
// // 	}
// // 	cost := float64((time.Now().UnixNano() - currTime) / 1000000)
// // 	fmt.Println("process timecost", cost)
// // }

// func GenerateLength(ffmpegPath string, url string, reqId string) (int, error) {
// 	ct := time.Now().UnixNano()
// 	var len int
// 	for i := 0; i < 2; i++ {
// 		//视频处理，延长超时时间
// 		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
// 		cmd := exec.CommandContext(ctx, ffmpegPath, "-i", url)
// 		defer cancel()

// 		var stdOut bytes.Buffer
// 		var stdErr bytes.Buffer
// 		cmd.Stdout = &stdOut
// 		cmd.Stderr = &stdErr
// 		cmd.Run()
// 		str := "2006-01-02" + strings.TrimPrefix(str, "Duration:")
// 		if videoTime, err := time.Parse("2006-01-02 15:04:05", str); err != nil {
// 			if ctx.Err() != nil {
// 				fmt.Println("GenerateLength Err:", ctx.Err())
// 			}
// 			len = 0
// 		} else {
// 			len = videoTime.Hour()*3600 + videoTime.Minute()*60 + videoTime.Second()
// 			break
// 		}
// 	}
// 	cost := float64((time.Now().UnixNano() - ct) / 1000000)
// 	fmt.Println("videolengthcost", cost)
// 	fmt.Println("---------->>>videolength:", len)
// 	return len, nil
// }

// //通过-ss参数 获取视频中的图片帧
// func testFfmpegParams(url string, path string, ffmpegPath string, freq int, videoLen int) string {
// 	ct := time.Now().UnixNano()
// 	var outPutError string
// 	for i := 0; i < videoLen; i += freq {
// 		sec := strconv.Itoa(i)
// 		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5000)*time.Millisecond)
// 		cmd := exec.CommandContext(ctx, ffmpegPath,
// 			"-fmtlevel", "error",
// 			"-y",
// 			"-ss", sec,
// 			"-t", "1",
// 			"-i", url,
// 			"-vframes", "1",
// 			path+"/"+sec+".jpg")
// 		defer cancel()
// 		var stdErr bytes.Buffer
// 		cmd.Stderr = &stdErr
// 		err := cmd.Run()
// 		if err != nil {
// 			outPutError += fmt.Sprintf("lastframecmderr:%v;", err)
// 		}
// 		if stdErr.Len() != 0 {
// 			outPutError += fmt.Sprintf("lastframestderr:%v;", stdErr.String())
// 		}
// 		if ctx.Err() != nil {
// 			outPutError += fmt.Sprintf("lastframectxerr:%v;", ctx.Err())
// 		}
// 	}
// 	cost := float64((time.Now().UnixNano() - ct) / 1000000)
// 	fmt.Println("jiepingcost", cost)
// 	return outPutError
// }

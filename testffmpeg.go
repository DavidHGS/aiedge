package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"unsafe"

// 	"github.com/giorgisio/goav/avcodec"
// 	"github.com/giorgisio/goav/avdevice"
// 	"github.com/giorgisio/goav/avfilter"
// 	"github.com/giorgisio/goav/avformat"
// 	"github.com/giorgisio/goav/avutil"
// 	"github.com/giorgisio/goav/swresample"
// 	"github.com/giorgisio/goav/swscale"
// )

// // avcodec 对应于 ffmpeg 库：libavcodec ——————————[提供更多编解码器的实现]
// // avformat 对应于 ffmpeg 库：libavformat ——————————[实现流媒体协议、容器格式和基本的I/O访问]。
// // avutil 对应于 ffmpeg 库：libavutil ——————————[包括哈希器、解压器和其他实用功能] 。
// // avfilter 对应于 ffmpeg 库：libavfilter ——————————[提供了一种通过过滤器链来改变已解码的音频和视频的方法]。
// // avdevice 对应于 ffmpeg 库：libavdevice ——————————[提供了一个访问捕获和播放设备的抽象概念] 。
// // swresample 对应 ffmpeg 库：libswresample ——————————[实现音频混合和重采样例程]。
// // swscale 对应于 ffmpeg 库：libswscale ——————————[实现颜色转换和缩放例程]。

// func main() {
// 	filename := "./src/sample.mp4"
// 	//加载ffmpeg网络库
// 	avformat.AvRegisterAll()
// 	//加载ffmpeg的编解码库
// 	avcodec.AvcodecRegisterAll()
// 	log.Printf("AvFilter Version:\t%v", avfilter.AvfilterVersion())
// 	log.Printf("AvDevice Version:\t%v", avdevice.AvdeviceVersion())
// 	log.Printf("SWScale Version:\t%v", swscale.SwscaleVersion())
// 	log.Printf("AvUtil Version:\t%v", avutil.AvutilVersion())
// 	log.Printf("AvCodec Version:\t%v", avcodec.AvcodecVersion())
// 	log.Printf("Resample Version:\t%v", swresample.SwresampleLicense())

// 	//打开视频流
// 	formatCtx := openInput(filename) //返回一个带有视频信息的Context
// 	if formatCtx == nil {
// 		return
// 	}
// 	//检索流信息
// 	success := findStreamInfo(formatCtx)
// 	if !success {
// 		return
// 	}
// 	fmt.Println("检索流信息：", success)
// 	//Print detailed information about the input or output format, such as duration, bitrate, streams, container, programs, metadata, side data, codec and time base.
// 	formatCtx.AvDumpFormat(0, filename, 0) //av_dump_format
// 	//读取一帧视频帧
// 	avCodecCtx := findFirstVideoStreamCodecContext(formatCtx)
// 	if avCodecCtx == nil {
// 		fmt.Println("无视频帧")
// 		return
// 	}
// 	//查找并打开解码器
// 	codecCtx := findAndOpenCodec(avCodecCtx)
// 	if codecCtx == nil {
// 		fmt.Println("没有发现解码器")
// 		return
// 	}

// 	//Allocate video frame
// 	pFrame := avutil.AvFrameAlloc() //Allocate an Frame and set its fields to default values.
// 	// Allocate an AVFrame structure
// 	pFrameRGB := avutil.AvFrameAlloc()
// 	// Determine required buffer size and allocate buffer
// 	numBytes := uintptr(avcodec.AvpictureGetSize((avcodec.PixelFormat)(avcodec.AV_PIX_FMT_RGB24), codecCtx.Width(), codecCtx.Height())) // Calculate the size in bytes that a picture of the given width and height would occupy if stored in the given picture format.
// 	buffer := avutil.AvMalloc(numBytes)
// 	// Assign appropriate parts of buffer to image planes in pFrameRGB
// 	// Note that pFrameRGB is an AVFrame, but AVFrame is a superset
// 	// of AVPicture
// 	avp := (*avcodec.Picture)(unsafe.Pointer(pFrameRGB))
// 	// AvpictureFill - Setup the picture fields based on the specified image parameters and the provided image data buffer. 设置picture的一些字段
// 	avp.AvpictureFill((*uint8)(buffer), (avcodec.PixelFormat)(avcodec.AV_PIX_FMT_RGB24), codecCtx.Width(), codecCtx.Height())

// 	//initialize SWS context for software scalling
// 	swsCtx := swscale.SwsGetcontext( //Allocate and return an Context.
// 		codecCtx.Width(),
// 		codecCtx.Height(),
// 		swscale.PixelFormat(codecCtx.PixFmt()),
// 		codecCtx.Width(),
// 		codecCtx.Height(),
// 		swscale.PixelFormat(avcodec.AV_PIX_FMT_RGB24),
// 		int(avcodec.SWS_BILINEAR),
// 		nil,
// 		nil,
// 		nil,
// 	)
// 	for i := 0; i < int(formatCtx.NbStreams()); i++ {
// 		switch formatCtx.Streams()[i].CodecParameters().AvCodecGetType() {
// 		case avformat.AVMEDIA_TYPE_VIDEO:

// 			//循环读取视频帧并解码成rgb 默认为yuv数据
// 			packet := avcodec.AvPacketAlloc()
// 			frameNumber := 1                         //帧数
// 			for formatCtx.AvReadFrame(packet) >= 0 { //Return the next frame of a stream.
// 				//是否为视频流的某一个包
// 				if packet.StreamIndex() == i {
// 					//Decode video frame
// 					response := codecCtx.AvcodecSendPacket(packet) //发送请求
// 					if response < 0 {
// 						fmt.Println("sending packet error")
// 					}
// 					if response >= 0 {
// 						response = codecCtx.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(pFrame))) //接收报文
// 						if response == avutil.AvErrorEAGAIN || response == avutil.AvErrorEOF {
// 							fmt.Println("AvcodecReceiveFrame Error:  avutil.AvErrorEAGAIN or  avutil.AvErrorEOF ")
// 							break
// 						} else if response < 0 {
// 							fmt.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
// 							return
// 						}
// 						// 从原生数据转换成RGB， 转换 5 个视频帧
// 						// Convert the image from its native format to RGB
// 						if frameNumber < 5 {
// 							swscale.SwsScale2(swsCtx, avutil.Data(pFrame), avutil.Linesize(pFrame), 0, codecCtx.Height(), avutil.Data(pFrameRGB), avutil.Linesize(pFrameRGB))
// 							//保存到本地磁盘
// 							fmt.Printf("Writing frame %d\n", frameNumber)
// 							SaveFrame(pFrameRGB, codecCtx.Width(), codecCtx.Height(), frameNumber)
// 						} else {
// 							return
// 						}
// 						frameNumber++
// 					}
// 					// 释放资源
// 					// Free the packet that was allocated by av_read_frame
// 					packet.AvFreePacket() //Free a packet.
// 				}
// 			}
// 			//Free the RGB image
// 			avutil.AvFree(buffer)
// 			avutil.AvFrameFree(pFrameRGB)
// 			// Free the YUV frame
// 			avutil.AvFrameFree(pFrame)
// 			//Close the codecs 关闭解码器
// 			codecCtx.AvcodecClose()
// 			(*avcodec.Context)(unsafe.Pointer(avCodecCtx)).AvcodecClose()
// 			// Close the video file 关闭视频文件
// 			formatCtx.AvformatCloseInput()
// 		default:
// 			fmt.Println("Didn't find a cideo stream")
// 		}
// 	}
// }

// /**
// 打开视频流
// */
// func openInput(filename string) *avformat.Context {
// 	formatCtx := avformat.AvformatAllocContext() ///Allocate an Context.
// 	//打开视频流
// 	// avformat.AvFindInputFormat("")
// 	if avformat.AvformatOpenInput(&formatCtx, filename, nil, nil) != 0 { //Open an input stream and read the header. 数据传到formatCtx
// 		fmt.Printf("unable to open file %s\n", filename)
// 		return nil
// 	}
// 	return formatCtx
// }

// /**
// 检索流信息
// */
// func findStreamInfo(ctx *avformat.Context) bool {
// 	if ctx.AvformatFindStreamInfo(nil) < 0 { ////Read packets of a media file to get stream information. 返回值为负数则失败 avformat_find_stream_info————C
// 		log.Println("Error:Couldn't find stream information.")
// 		ctx.AvformatCloseInput()
// 		return false
// 	}
// 	return true
// }

// /**
// 获取第一帧视频位置
// */
// func findFirstVideoStreamIndex(ctx *avformat.Context) int {
// 	videoStreamIndex := -1
// 	for index, stream := range ctx.Streams() {
// 		switch stream.CodecParameters().AvCodecGetType() {
// 		case avformat.AVMEDIA_TYPE_VIDEO:
// 			return index
// 		}
// 	}
// 	return videoStreamIndex
// }

// /**
// 读取一帧视频帧
// */
// func findFirstVideoStreamCodecContext(ctx *avformat.Context) *avformat.CodecContext {
// 	//得到一个指针的视频编解码器上下文
// 	for _, stream := range ctx.Streams() {
// 		switch stream.CodecParameters().AvCodecGetType() {
// 		case avformat.AVMEDIA_TYPE_VIDEO:
// 			return stream.Codec()
// 		}
// 	}
// 	return nil
// }

// /***********************************************************************/
// /**
// 查找并打开编解码器
// */
// func findAndOpenCodec(codecCtx *avformat.CodecContext) *avcodec.Context {
// 	codec := avcodec.AvcodecFindDecoder(avcodec.CodecId(codecCtx.GetCodecId())) //Find a registered decoder with a matching codec ID.通过codecCtx找到一个解码器
// 	if codec == nil {
// 		fmt.Println("Unsupported codec")
// 		return nil
// 	}
// 	// Copy context
// 	codecContext := codec.AvcodecAllocContext3()                                            //Allocate an Context and set its fields to default values. 解码器分配一个codecContext
// 	if codecContext.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(codecCtx))) != 0 { //Copy the settings of the source Context into the destination Context. 将avformat.CodecContext转化为
// 		fmt.Println("Couldn't copy codec context")
// 		return nil
// 	}
// 	return codecContext
// }

// //保存一帧图片到磁盘
// func SaveFrame(frame *avutil.Frame, width, height, frameNumber int) {
// 	//openfile
// 	fileName := fmt.Sprintf("images/frame%d.ppm", frameNumber)
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		log.Println("file Create faile")
// 	}
// 	defer file.Close()
// 	//Writer header
// 	header := fmt.Sprintf("P6\n%d %d\n255\n", width, height)
// 	file.Write([]byte(header))

// 	//Write puixel data
// 	for k := 0; k < height; k++ {
// 		data0 := avutil.Data(frame)[0]
// 		buf := make([]byte, width*3) //RGB so width *3
// 		startPos := uintptr(unsafe.Pointer(&data0)) + uintptr(k)*uintptr(avutil.Linesize(frame)[0])
// 		for i := 0; i < width*3; i++ {
// 			element := *(*uint8)(unsafe.Pointer(startPos + uintptr(i)))
// 			buf[i] = element
// 		}
// 		file.Write(buf)
// 	}
// }
/************************************************************************************/

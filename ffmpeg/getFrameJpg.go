package ffmpeg

// // tutorial01.c
// // Code based on a tutorial at http://dranger.com/ffmpeg/tutorial01.html

// // A small sample program that shows how to use libavformat and libavcodec to
// // read video from a file.
// //
// // Use
// //
// // gcc -o tutorial01 tutorial01.c -lavformat -lavcodec -lswscale -lz
// //
// // to build (assuming libavformat and libavcodec are correctly installed
// // your system).
// //
// // Run using
// //
// // tutorial01 myvideofile.mpg
// //
// // to write the first five frames from "myvideofile.mpg" to disk in PPM
// // format.
// import (
// 	"fmt"
// 	"fmt"
// 	"os"
// 	"unsafe"

// 	"github.com/giorgisio/goav/avcodec"
// 	"github.com/giorgisio/goav/avformat"
// 	"github.com/giorgisio/goav/avutil"
// 	"github.com/giorgisio/goav/swscale"
// )

// // 视频中提取照片：
// // 1、fmpeg对将像素数据写入到JPG图片中也封装到了avformat_xxx系列接口中，它的使用流程和封装视频数据到mp4文件一模一样，只不过一个JPG文件中只包含了一帧视频数据而已;
// // 2、ffmpeg对JPG文件的封装支持模式匹配，即如果想要将多张图片写入到多张jpg中只需要文件名包含百分号即可，例如 name%3d.jpg，那么在每一次调用av_write_frame()函数写入视频数据到jpg图片时都会生成一张jpg图片。这样做的好处是不需要每一张要写入的jpg文件都创建一个AVFormatContext与之对应。其它流程和写入一张jpg一样。

// // 流程为：
// // 1、先从MP4中提取指定时刻AVPacket解码成AVFrame
// // 2、然后将步骤1得到的AVFrame进行从素格式YUV420P到JPG需要的YUVJ420P像素格式的转换
// // 3、再重新编码，然后再封装到jpg中

// // 照片合成视频：
// // 因为JPG的编码方式为AV_CODEC_ID_MJPEG，MP4如果采用h264编码，那么两者的编码方式是不一致的，所以就需要先解码再编码，具体流程为：
// // 1、先将JPG解码成AVFrame
// // 2、将JPG解码后的源像素格式YUVJ420P转换成x264编码需要的YUV420P像素格式
// // 3、再重新编码，然后再封装到mp4中

// // SaveFrame writes a single frame to disk as a PPM file
// func SaveFrame(frame *avutil.Frame, width, height, frameNumber int) {
// 	// Open file
// 	fileName := fmt.Sprintf("./images/frame%d.ppm", frameNumber)
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		fmt.Println("Error Reading")
// 	}
// 	defer file.Close()

// 	// Write header
// 	header := fmt.Sprintf("P6\n%d %d\n255\n", width, height)
// 	file.Write([]byte(header))

// 	// Write pixel data
// 	for y := 0; y < height; y++ {
// 		data0 := avutil.Data(frame)[0]
// 		buf := make([]byte, width*3)
// 		startPos := uintptr(unsafe.Pointer(data0)) + uintptr(y)*uintptr(avutil.Linesize(frame)[0])
// 		for i := 0; i < width*3; i++ {
// 			element := *(*uint8)(unsafe.Pointer(startPos + uintptr(i)))
// 			buf[i] = element
// 		}
// 		file.Write(buf)
// 	}
// }

// func main() {
// 	filename := "./src/sample.mp4"
// 	dstPath := "./images/output%d.jpg"
// 	// Open video file
// 	pFormatContext := avformat.AvformatAllocContext()
// 	if avformat.AvformatOpenInput(&pFormatContext, filename, nil, nil) != 0 {
// 		fmt.Printf("Unable to open file %s\n", filename)
// 		os.Exit(1)
// 	}

// 	// Retrieve stream information
// 	if pFormatContext.AvformatFindStreamInfo(nil) < 0 {
// 		fmt.Println("Couldn't find stream information")
// 		os.Exit(1)
// 	}

// 	// Dump information about file onto standard error
// 	pFormatContext.AvDumpFormat(0, filename, 0)
// 	video_index := -1
// 	var DeCodecCtx *avcodec.Context
// 	var DeCodecCtxOrig *avformat.CodecContext
// 	//遍历出视频索引 Find the first video stream
// 	for i := 0; i < int(pFormatContext.NbStreams()); i++ {
// 		stream := pFormatContext.Streams()[i]
// 		if stream.CodecParameters().AvCodecGetType() == avformat.AVMEDIA_TYPE_VIDEO { //说明是video
// 			video_index = i
// 			//初始化解码器用于解码 Get a pointer to the codec context for the video stream
// 			DeCodecCtxOrig := stream.Codec()
// 			// Find the decoder for the video stream
// 			DeCodec := avcodec.AvcodecFindDecoder(avcodec.CodecId(DeCodecCtxOrig.GetCodecId()))
// 			if DeCodec == nil {
// 				fmt.Println("Unsupported codec!")
// 				os.Exit(1)
// 			}
// 			// Copy context
// 			DeCodecCtx := DeCodec.AvcodecAllocContext3()
// 			// 设置解码参数，这里直接从源视频流中拷贝
// 			if DeCodecCtx.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(DeCodecCtxOrig))) != 0 {
// 				fmt.Println("Couldn't copy codec context")
// 				os.Exit(1)
// 			}

// 			// 初始化解码器上下文 Open codec
// 			if DeCodecCtx.AvcodecOpen2(DeCodec, nil) < 0 {
// 				fmt.Println("Could not open codec")
// 				os.Exit(1)
// 			}
// 			break
// 		}
// 	}
// 	// 初始化编码器;因为最终是要写入到JPEG，所以使用的编码器ID为AV_CODEC_ID_MJPEG
// 	EnCodec := avcodec.AvcodecFindEncoder(avcodec.CodecId(avcodec.AV_CODEC_ID_MJPEG))
// 	EnCodecCtx := EnCodec.AvcodecAllocContext3()
// 	//设置编码参数
// 	inStream := pFormatContext.Streams()[video_index]
// 	EnCodecCtxOrig := inStream.Codec()
// 	// 对于MJPEG编码器来说，它支持的是YUVJ420P/YUVJ422P/YUVJ444P格式的像素
// 	EnCodecCtx.SetEncodeParams(EnCodecCtxOrig.GetWidth(), EnCodecCtxOrig.GetHeight(), avcodec.AV_PIX_FMT_YUV420P9)
// 	EnCodecCtx.SetTimebase(EnCodecCtxOrig.GetTimeBase().Num(), EnCodecCtxOrig.GetTimeBase().Den())

// 	// 初始化编码器上下文
// 	if EnCodecCtx.AvcodecOpen2(EnCodec, nil) < 0 {
// 		fmt.Println("EnCodecCtx.AvcodecOpen2 fail")
// 		os.Exit(1)
// 	}
// 	// 创建用于输出JPG的封装器
// 	var outFormat *avformat.Context
// 	if avformat.AvformatAllocOutputContext2(&outFormat, nil, "", dstPath) < 0 {
// 		fmt.Println("avformat.AvformatAllocOutputContext2 fail")
// 		os.Exit(1)
// 	}

// 	/** 添加流
// 	 *  对于图片封装器来说，也可以把它想象成只有一帧视频的视频封装器。所以它实际上也需要一路视频流，而事实上图片的流是视频流类型
// 	 */
// 	stream := outFormat.AvformatNewStream(nil)
// 	// 设置流参数；直接从编码器拷贝参数即可
// 	if EnCodecCtx.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(stream.Codec()))) < 0 {
// 		fmt.Println("EnCodecCtx.AvcodecCopyContext fail")
// 		os.Exit(1)
// 	}
// 	/** 初始化上下文
// 	 *  对于写入JPG来说，它是不需要建立输出上下文IO缓冲区的的，所以avio_open2()没有调用到，但是最终一样可以调用av_write_frame()写入数据
// 	 */

// 	/** 为输出文件写入头信息
// 	 *  不管是封装音视频文件还是图片文件，都需要调用此方法进行相关的初始化，否则av_write_frame()函数会崩溃
// 	 */
// 	if outFormat.AvformatWriteHeader(nil) < 0 {
// 		fmt.Println("outFormat.AvformatWriteHeader fail")
// 		os.Exit(1)
// 	}

// 	// Allocate video frame
// 	pFrame := avutil.AvFrameAlloc()

// 	// Allocate an AVFrame structure
// 	pFrameRGB := avutil.AvFrameAlloc()
// 	if pFrameRGB == nil {
// 		fmt.Println("Unable to allocate RGB Frame")
// 		os.Exit(1)
// 	}

// 	// Determine required buffer size and allocate buffer
// 	numBytes := uintptr(avcodec.AvpictureGetSize(avcodec.AV_PIX_FMT_RGB24, DeCodecCtx.Width(),
// 		DeCodecCtx.Height()))
// 	buffer := avutil.AvMalloc(numBytes)

// 	// Assign appropriate parts of buffer to image planes in pFrameRGB
// 	// Note that pFrameRGB is an AVFrame, but AVFrame is a superset
// 	// of AVPicture
// 	avp := (*avcodec.Picture)(unsafe.Pointer(pFrameRGB))
// 	avp.AvpictureFill((*uint8)(buffer), avcodec.AV_PIX_FMT_RGB24, DeCodecCtx.Width(), DeCodecCtx.Height())

// 	// initialize SWS context for software scaling
// 	swsCtx := swscale.SwsGetcontext(
// 		DeCodecCtx.Width(),
// 		DeCodecCtx.Height(),
// 		(swscale.PixelFormat)(DeCodecCtx.PixFmt()),
// 		DeCodecCtx.Width(),
// 		DeCodecCtx.Height(),
// 		avcodec.AV_PIX_FMT_RGB24,
// 		avcodec.SWS_BILINEAR,
// 		nil,
// 		nil,
// 		nil,
// 	)

// 	// Read frames and save first five frames to disk
// 	frameNumber := 1
// 	packet := avcodec.AvPacketAlloc()
// 	for pFormatContext.AvReadFrame(packet) >= 0 {
// 		// Is this a packet from the video stream?
// 		if packet.StreamIndex() == i {
// 			// Decode video frame
// 			response := DeCodecCtx.AvcodecSendPacket(packet)
// 			if response < 0 {
// 				fmt.Printf("Error while sending a packet to the decoder: %s\n", avutil.ErrorFromCode(response))
// 			}
// 			for response >= 0 {
// 				response = DeCodecCtx.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(pFrame)))
// 				if response == avutil.AvErrorEAGAIN || response == avutil.AvErrorEOF {
// 					break
// 				} else if response < 0 {
// 					if response == -11 {
// 						fmt.Println("response error = - 11") //ReadEagain错误
// 						break
// 					}
// 					fmt.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
// 					return
// 				}

// 				if frameNumber <= 30 {
// 					// Convert the image from its native format to RGB
// 					swscale.SwsScale2(swsCtx, avutil.Data(pFrame),
// 						avutil.Linesize(pFrame), 0, DeCodecCtx.Height(),
// 						avutil.Data(pFrameRGB), avutil.Linesize(pFrameRGB))

// 					// Save the frame to disk
// 					fmt.Printf("Writing frame %d\n", frameNumber)
// 					SaveFrame(pFrameRGB, DeCodecCtx.Width(), DeCodecCtx.Height(), frameNumber)
// 				} else {
// 					return
// 				}
// 				frameNumber++
// 			}
// 		}

// 		// Free the packet that was allocated by av_read_frame
// 		packet.AvFreePacket()
// 	}

// 	// Free the RGB image
// 	avutil.AvFree(buffer)
// 	avutil.AvFrameFree(pFrameRGB)

// 	// Free the YUV frame
// 	avutil.AvFrameFree(pFrame)

// 	// Close the codecs
// 	DeCodecCtx.AvcodecClose()
// 	(*avcodec.Context)(unsafe.Pointer(DeCodecCtxOrig)).AvcodecClose()

// 	// Close the video file
// 	pFormatContext.AvformatCloseInput()

// 	// Stop after saving frames of first video straem

// }

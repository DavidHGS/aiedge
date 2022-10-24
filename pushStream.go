package main

// import (
// 	"bytes"
// 	"fmt"
// 	"log"
// 	"os/exec"
// )

// 	func main() {
// 	command := "ffmpeg -re -stream_loop -1 -i video02.mp4 -vcodec copy -acodec copy -f flv -y rtmp://aiedge.ndsl-lab.cn:8035/live/stream"
// 	cmd := exec.Command("/bin/bash", "-c", command)
// 	var out bytes.Buffer
// 	cmd.Stdout = &out

// 	if err := cmd.Run(); err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Printf("command output: %q", out.String())
// }

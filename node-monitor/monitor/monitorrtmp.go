package monitor

import (
	"context"
	"errors"
	"node-monitor/common"
	"os/exec"
	"path"
	"strings"
	"time"

	"node-monitor/libs/logs"
	//rtmp "github.com/retailnext/rtmpclient"
	"github.com/xxjwxc/gowp/workpool"
)

type urlStream struct {
	url    string
	stream string
}

func Monitorrtmp() {

	// var urlAll []urlStream = make([]urlStream)
	urlAll := []urlStream{}
	for _, application := range common.Cfg.Rtmpmonitor {
		url := "rtmp://" + application.Domain + "/" + application.App
		streams := strings.Split(application.Stream, ",")
		for _, s := range streams {
			urlAll = append(urlAll, urlStream{url, s})
		}
	}
	// logs.Info(urlAll)
	//使用工作池,批量监控
	wp := workpool.New(5) // 设置最大线程数
	// wp.SetTimeout(1200 * time.Second) // 设置超时时间
	for i := 0; i < len(urlAll); i++ {
		url := urlAll[i].url
		s := urlAll[i].stream
		wp.Do(func() error {

			monitorRtmpStream(url, s)

			return nil
		})
	}
	wp.Wait()
	// err := wp.Wait()
	// if err != nil {
	// 	logs.Error("Workpool error,err", err)
	// }

}

func monitorRtmpStream(url string, s string) (err error) {
	logs.Info("Monitoring rtmp stream, url, ", url, s)
	// go func() {
	// var err error
	// var err1 error

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	err = rtmpplayerWithtimeout(ctx, url, s)

	if err != nil {
		logs.Error("Failed to play,err,", err, ",url,", url, s, ",try again")
		//失败了，重新连接一次，减少误判
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
		defer cancel()
		err = rtmpplayerWithtimeout(ctx, url, s)

		if err != nil {
			logs.Error("Failed to play,err,", err, ",url,", url, s)
			msg := "视频流检查,url," + url + "/" + s + ",不能正常播放,err," + err.Error()
			Send_Telegram_message(msg)
		}
	}
	return
}

//url rtmp://192.168.20.111/vid3
//streamName streamname
func rtmpplayerWithtimeout(ctx context.Context, url string, streamName string) (err error) {

	hctx, hcancel := context.WithTimeout(ctx, time.Second*30) //设置超时时间300s
	defer hcancel()
	resp := make(chan string, 1)
	respError := make(chan error, 1)

	cmd := exec.Command(path.Join(common.Workdir, "/ffmpeg/ffprobe"), "-v", "quiet", "-print_format", "json", "-show_streams", url+"/"+streamName)
	// cmd := exec.Command("ls", "-v quiet -print_format json -show_streams", "rtmp://43.245.200.167/live/dgbc0101spc")
	// if runtime.GOOS = "windows" {
	//     cmd := exec.Command("tasklist")
	// }

	go func() {

		output, err := cmd.CombinedOutput()

		if err != nil {

			logs.Error("cmd running failed,err,", err, ",output,", string(output))
			respError <- err
			return

		}
		resp <- "finish"
	}()

	// 超时机制
	select {
	//	case <-ctx.Done():
	//		fmt.Println("ctx timeout")
	//		fmt.Println(ctx.Err())
	case <-hctx.Done():
		err = errors.New(url + "/" + streamName + ", No data after 30 seconds")
		logs.Error("Error Monitoring,exceed 30s,timeout", err)
		cmd.Process.Kill()

	case v := <-resp:
		if v == "finish" {
			logs.Info("Finish monitoring rtmp", url+"/"+streamName)
		}
	case e := <-respError:
		err = e
	}

	return
}

//rtmp播放库,耗内存,废弃
//url rtmp://192.168.20.111/vid3
//streamName streamname
// func rtmpplayerRetailnext(url string, streamName string) (err error) {
// 	conn, err := rtmp.Dial(url, 100)
// 	if err != nil {
// 		logs.Error(url+"/"+streamName, "Failed to Dial url,", url, ",err,", err)
// 		return
// 	}

// 	f, err := os.Create(path.Join(common.Workdir, "video_dump"))
// 	if err != nil {
// 		logs.Error(url+"/"+streamName, "Failed to create video dump file,err,", err)
// 		return
// 	}
// 	defer f.Close()

// 	defer conn.Close()

// 	logs.Debug(url+"/"+streamName, "Trying to connect url,", url, ",conn,", conn)

// 	err = conn.Connect()
// 	if err != nil {
// 		logs.Error(url+"/"+streamName, "Failed to connect url,", url, ",conn,", conn)
// 		return
// 	}

// 	streamIDs := make([]uint32, 0)

// 	for done := false; !done; {
// 		select {
// 		case msg, ok := <-conn.Events():
// 			if !ok {
// 				done = true
// 			}
// 			switch ev := msg.Data.(type) {
// 			case *rtmp.StatusEvent:
// 				logs.Debug(url+"/"+streamName, "evt status:", ev.Status)
// 			case *rtmp.ClosedEvent:
// 				logs.Debug(url+"/"+streamName, "evt closed")
// 			case *rtmp.VideoEvent:
// 				logs.Debug(url+"/"+streamName, "evt video:", ev.Message.Timestamp, ev.Message.Buf.Len())
// 				logs.Debug(url+"/"+streamName, "Play video successfully,conn,", conn, streamName, ",get message,", ev.Message)
// 				return
// 			case *rtmp.AudioEvent:
// 				logs.Debug(url+"/"+streamName, "evt audio")
// 			case *rtmp.CommandEvent:
// 				logs.Debug(url+"/"+streamName, "evt command")
// 			case *rtmp.StreamBegin:
// 				logs.Debug(url+"/"+streamName, "case *rtmp.StreamBegin")
// 			case *rtmp.StreamEOF:
// 				logs.Debug(url+"/"+streamName, "case *rtmp.StreamEOF", ev.StreamID, streamIDs)
// 				return
// 			case *rtmp.StreamDry:
// 				logs.Debug(url+"/"+streamName, "case *rtmp.StreamDry")
// 			case *rtmp.StreamIsRecorded:
// 				logs.Debug(url+"/"+streamName, "case *rtmp.StreamIsRecorded")
// 			case *rtmp.StreamCreatedEvent:
// 				streamIDs = append(streamIDs, ev.Stream.ID())
// 				logs.Debug(url+"/"+streamName, "evt stream created", ev.Stream.ID())
// 				err = ev.Stream.Play(streamName, nil, nil, nil)
// 				if err != nil {
// 					logs.Error(url+"/"+streamName, "Play error: %s", err.Error())
// 					return
// 				}
// 			}
// 		case <-time.After(30 * time.Second):
// 			err = errors.New(url + "/" + streamName + ", No data after 30 seconds")
// 			logs.Error(url+"/"+streamName, "failed to play,err,", err)
// 			return
// 		}
// 	}
// 	return
// }

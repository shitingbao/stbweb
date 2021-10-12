//首先下一个 ffmpeg 的可执行文件
package ffmpeg

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	defaultOutPath = "outFiles" // 输出文件地址
	// defaultInPath  = "./"       //
)

type Ffmpeg interface {
	MovToMp4() error
}

type (
	Option func(*options)
)

type options struct {
	InPath, OutPath string
}

type ffmpeg struct {
	order, inPath, outPath string
}

func NewFfmpeg(opts ...Option) Ffmpeg {
	f := &options{
		// InPath:  defaultInPath,
		OutPath: defaultOutPath,
	}
	for _, o := range opts {
		o(f)
	}
	return &ffmpeg{
		order:   "ffmpeg",
		inPath:  f.InPath,
		outPath: f.OutPath,
	}
}

// MovToMp4
func (f *ffmpeg) MovToMp4() error {
	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}
	fds, err := os.ReadDir(f.inPath)
	if err != nil {
		// log.Println("read file path have err:", err)
		return err
	}
	wholePath := path.Join(rootPath, f.outPath+strconv.Itoa(int(time.Now().Unix())))
	if err := os.MkdirAll(wholePath, os.ModePerm); err != nil {
		return err
	}
	for _, v := range fds {
		l := strings.Split(v.Name(), ".") // 防止名字中有多个英文 “ . ”
		if len(l) > 1 && l[len(l)-1] == "mov" {
			fileName := strings.Join(l[:len(l)-1], ".")
			fileName += ".mp4"
			wholeOutPath := path.Join(wholePath, fileName)

			ecPath := path.Join(rootPath, f.order)
			//不要写整条命令！！！
			//不要写整条命令！！！
			//不要写整条命令！！！
			cmd := exec.Command(ecPath, "-i", path.Join(rootPath, v.Name()), "-qscale", "0", wholeOutPath) //ffmpeg -i input.mov -qscale 0 output.mp4
			if err := cmd.Run(); err != nil {
				return err
			} else {
				log.Println("success change ", v.Name())
			}
		}
	}
	return nil
}

func WithFfmpegRpath(inpath string) Option {
	return func(o *options) {
		o.InPath = inpath
	}
}

func WithFfmpegOutpath(outpath string) Option {
	return func(o *options) {
		o.OutPath = outpath
	}
}

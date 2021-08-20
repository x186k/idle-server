package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var gstformat = `/usr/bin/gst-launch-1.0 filesrc location="%s" !
decodebin name=decode !
videoconvert !
x264enc option-string=slice-max-size=1200 speed-preset=medium tune=zerolatency key-int-max=1 !
video/x-h264,profile=constrained-baseline !
queue max-size-time=100000000 !
rtph264pay config-interval=-1 name=payloader !
multifilesink location="%s"`

const infile = "/foo/idle-media"
const outdir = "/foo/out.nosync"

var gstcmd = fmt.Sprintf(gstformat, infile, outdir+"/rtp%d.rtp")

func checkFatal(err error) {
	if err != nil {
		_, fileName, fileLine, _ := runtime.Caller(1)
		log.Fatalf("FATAL %s:%d %v", filepath.Base(fileName), fileLine, err)
	}
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		checkFatal(err)
	}
	return true
}

func main() {
	var err error

	if !Exists(infile) {
		checkFatal(fmt.Errorf("no in file"))
		// http server
	}

	_ = os.RemoveAll(outdir)
	err = os.Mkdir(outdir, 0777)
	checkFatal(err)

	args := strings.Fields(gstcmd)

	// for k, v := range args {
	// 	println(k, v)
	// }
	log.Printf("Running command2")

	cmd := exec.Command(args[0], args[1:]...)
	stdoutStderr, err := cmd.CombinedOutput()
	checkFatal(err)

	fmt.Printf("output: %s\n", stdoutStderr)

	log.Printf("Command err code: %v", err)

	n, err := ioutil.ReadDir(outdir)
	checkFatal(err)

	log.Printf("n files after gstreamer: %v", len(n))
}

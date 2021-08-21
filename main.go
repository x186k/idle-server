package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/pflag"
)

var gstformat = `/usr/bin/gst-launch-1.0 filesrc location="%s" !
decodebin name=decode !
videoconvert !
x264enc option-string=slice-max-size=1200 speed-preset=medium tune=zerolatency key-int-max=1 !
video/x-h264,profile=constrained-baseline !
queue max-size-time=100000000 !
rtph264pay config-interval=-1 name=payloader !
multifilesink location="%s"`

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

	input := pflag.String("input", "", "input file")

	pflag.Parse()

	if *input != "" {
		if !Exists(*input) {
			checkFatal(fmt.Errorf("input file does not exist"))
		}

		zipbuf, err := runGstreamer(*input)
		checkFatal(err)

		//ext := path.Ext(*input)
		outfile := "idle-clip.zip"

		err = ioutil.WriteFile(outfile, zipbuf, 0666)
		checkFatal(err)

		log.Println("output written to ", outfile)

		return

	}

	log.Println("No --input flag, waiting for http requests")

	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			httpError(fmt.Errorf("http POST ONLY"), w)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			httpError(err, w)
			return
		}

		infile, err := ioutil.TempFile("", "idle-server-input")
		if err != nil {
			httpError(err, w)
			return
		}
		defer os.Remove(infile.Name())
		log.Println("temp file is ", infile.Name())

		n, err := infile.Write(body)
		if err != nil {
			httpError(err, w)
			return
		}
		if n != len(body) {
			httpError(fmt.Errorf("bad file write len"), w)
			return
		}

		zipbuf, err := runGstreamer(infile.Name())
		if err != nil {
			httpError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename='idle-clip.zip'")

		n, err = w.Write(zipbuf)
		if err != nil {
			httpError(err, w)
			return
		}

		if n != len(zipbuf) {
			httpError(fmt.Errorf("bad resp write len"), w)
			return
		}

	})

	err = http.ListenAndServe(":8080", nil)
	checkFatal(err)

}

func httpError(err error, rw http.ResponseWriter) {
	log.Println(err.Error())
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}

func runGstreamer(infile string) ([]byte, error) {
	outdir, err := ioutil.TempDir("", "idle-server-dir")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(outdir)

	var gstcmd = fmt.Sprintf(gstformat, infile, outdir+"/rtp%d.rtp")
	log.Printf("Running command2")
	args := strings.Fields(gstcmd)
	cmd := exec.Command(args[0], args[1:]...)
	stdoutStderr, err := cmd.CombinedOutput()
	log.Printf("output: %s\n", stdoutStderr)
	checkFatal(err)
	if err != nil {
		return nil, fmt.Errorf("cmd.CombinedOutput() %w", err)
	}

	log.Printf("Command err code: %v", err)

	rtpfiles, err := ioutil.ReadDir(outdir)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadDir(outdir) %w", err)
	}

	buf := new(bytes.Buffer)

	w := zip.NewWriter(buf)

	for _, file := range rtpfiles {
		f, err := w.Create(filepath.Base(file.Name()))
		if err != nil {
			return nil, fmt.Errorf("w.Create() %w", err)
		}
		pktbody, err := ioutil.ReadFile(path.Join(outdir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("ioutil.ReadFile(...) %w", err)
		}

		_, err = f.Write(pktbody)
		if err != nil {
			return nil, fmt.Errorf("f.Write(pktbody) %w", err)
		}
	}

	w.Close() //important

	log.Println(len(rtpfiles), "packets zipped up")

	return buf.Bytes(), nil
}

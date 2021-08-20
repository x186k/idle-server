package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {

	log.Printf("Running command and waiting for it to finish...")

	cmd := exec.Command("/usr/bin/gst-inspect-1.0")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)

	log.Printf("Command finished with error: %v", err)
}

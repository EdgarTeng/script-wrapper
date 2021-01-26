package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const ShellToUse = "bash"

func executeContent(content []byte) error {
	cmd := exec.Command(ShellToUse)
	cmd.Stdin = bytes.NewReader(content)
	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	_, err := io.Copy(os.Stdout, &buf)
	return err
}

func readFile() []byte {
	content, err := ioutil.ReadFile("plain.sh")
	if err != nil {
		log.Fatal(err)
	}

	return content
}

func main() {
	executeContent(readFile())
}

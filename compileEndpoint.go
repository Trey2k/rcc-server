package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func compileEndpoint(rw http.ResponseWriter, req *http.Request) {

	req.ParseMultipartForm(1000 << 20)

	file, handler, err := req.FormFile("tarball")
	if err != nil {
		panic(err)
	}

	dir := fmt.Sprintf("%s/%s", Config.WorkingDir, handler.Filename)

	defer file.Close()

	err = Untar(Config.WorkingDir, file)
	if err != nil {
		panic(err)
	}

	output, outFile, err := compile(dir)
	if err != nil {
		panic(err)
	}

	err = sendbackFile(outFile, string(output), rw)
	if err != nil {
		panic(err)
	}

}

func sendbackFile(filename, output string, rw http.ResponseWriter) error {

	body := &bytes.Buffer{}
	err := Tar(filename, body)
	if err != nil {
		rw.Header().Set("output", fmt.Sprintf("%s\n%s", output, "Output file/folder specified not found"))
		rw.WriteHeader(http.StatusInternalServerError)
		_, err = rw.Write(body.Bytes())
		return err
	}

	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.Header().Set("output", output)
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(body.Bytes())

	return err

}

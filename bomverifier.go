package main

import (
	"archive/zip"
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Bom struct {
	Libs []struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"libs"`
}

/*

This module is designed to run as a docker container at a Jenkins pipeline agent as part of the BOM verification after maven or cradle packaging Springboot applications.
The idea is that it verifies if the libraries contained do match the bill of material set by the company.


1. reads BOM Bill Of Material for Springboot libraries from github or bitbucket repo
2. Scans all Springboot jar files for libraries not matching the BOM
3. Report libraries that do not conform BOM
4. Exists with a non 0 exist code if -terminate flag is specified on comandline, that will break the jenkins build
 */
func main() {

	exitCode := 0

	param1 := flag.String("url", "https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.txt", "URL BOM bill of material yaml file.")

	param2 := flag.Bool("terminate", false, "Terminate jenkins pipeline on violation.")

	flag.Parse()

	uri := *param1
	terminate := *param2

	//uri = "https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.txt"

	fmt.Printf("URI bom txt file : %s termination is : %t\n", uri, terminate)

	fmt.Printf("DownloadToFile From: %s.\n", uri)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if d, err := HTTPDownload(uri); err == nil {
		fmt.Printf("downloaded %s.\n", uri)
		if WriteFile("bom.txt", d) == nil {
			fmt.Printf("saved %s as bom.txt\n", uri)
		}
	} else {
		panic(err)
	}

	var libs []string

	readFile("./bom.txt", &libs)

	fmt.Printf("size array = %d\n", len(libs))

	// Go through the Jenkins workspace directory, scanning for Springboot jar files
	files, err := FilePathWalkDir(".")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file, "jar") {
			fmt.Printf("Found SpringBoot application for verification => %s\n\n", file)
			read, err := zip.OpenReader(file)
			if err != nil {
				msg := "Failed to open: %s"
				log.Fatalf(msg, err)
			}
			defer read.Close()

			for _, infile := range read.File {

				if err := scanLibraries(infile.Name, file, libs, &exitCode, terminate); err != nil {
					log.Fatalf("Failed to read %s from zip: %s", infile.Name, err)
				}
			}

		}
	}

	os.Exit(exitCode)
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

/*
Match all contained libraries to matching the BOM listed jar files, no match...report conflict on stdout
 */
func scanLibraries(file string, filename string, libs []string, exitCode *int, terminate bool) error {

	var match = false

	if strings.Contains(file, ".jar") {

		for i := 0; i < len(libs); i++ {
			foundLib := strings.Split(file, "lib/")[1]
			foundLib = strings.Split(foundLib, ".jar")[0]
			matchLib := strings.Split(libs[i], ".jar")[0]
			//fmt.Printf("FOUND LIB %s\n", foundLib)
			//fmt.Printf("MATCHING LIB %s\n", matchLib)
			if strings.Compare(matchLib, foundLib) == 0 {
				match = true
			}

		}

		if match == false {
			fmt.Println("offending library : ", strings.Split(file, "lib/")[1])
			if terminate == true {
				*exitCode = 1
			}
		}

	}

	return nil
}

func HTTPDownload(uri string) ([]byte, error) {
	fmt.Printf("HTTPDownload From: %s.\n", uri)
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ReadFile: Size of download: %d\n", len(d))
	return d, err
}

func WriteFile(dst string, d []byte) error {
	fmt.Printf("WriteFile: Size of download: %d\n", len(d))
	err := ioutil.WriteFile(dst, d, 0444)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func readFile(fn string, libs *[]string) (err error) {

	file, err := os.Open(fn)
	defer file.Close()

	if err != nil {
		return err
	}

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)
	var line string
	for {
		line, err = reader.ReadString('\n')

		// Process the line here.
		//fmt.Println(line)
		*libs = append(*libs, line)
		if err != nil {
			break
		}
	}

	if err != io.EOF {
		fmt.Printf(" > Failed!: %v\n", err)
	}

	return
}
func limitLength(s string, length int) string {
	if len(s) < length {
		return s
	}

	return s[:length]
}

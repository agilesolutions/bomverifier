package main
 
import (
	"encoding/json"
	"archive/zip"
	"fmt"
	"log"
	"os"
	"io"
	"path/filepath"
	"strings"
	"net/http"
	"gopkg.in/yaml.v2"
)


 
 
func main() {

    if len(os.Args) != 2 {
        fmt.Println("Usage:", os.Args[0], "uri")
        return
    }

    uri := os.Args[1]
    
    uri = "https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.yaml"

    fmt.Println("URI bom yaml file : ", uri )
    fmt.Println()
    
    fmt.Printf("DownloadToFile From: %s.\n", uri)
    if d, err := HTTPDownload(uri); err == nil {
        fmt.Printf("downloaded %s.\n", uri)
        if WriteFile("bom.yaml", d) == nil {
            fmt.Printf("saved %s as %s\n", uri, dst)
        }
    }

	type bom struct {
    	Libs []struct {
            Name       string `yaml:"name"`
            Version    int    `yaml:"version"`
        } `yaml:"libs"`
	}
	var bom Bom
	
	err = yaml.Unmarshal("bom.yaml", &bom)
	if err != nil {
	    panic(err)
	}

	fmt.Print(service.libs[0].Name)
    

 	// filepath.Walk
 	files, err := FilePathWalkDir(".")
 	if err != nil {
  	panic(err)
 	}
 	for _, file := range files{
  		if (strings.HasSuffix(file, "jar")) {
	  		//fmt.Println(file)
			read, err := zip.OpenReader(file )
			if err != nil {
				msg := "Failed to open: %s"
				log.Fatalf(msg, err)
			}
			defer read.Close()

			for _, infile := range read.File {
				
			
				if err := listFiles(infile, file, expression, configuration.Copyto); err != nil {
				log.Fatalf("Failed to read %s from zip: %s", infile.Name, err)
				}
			}

  		}
 	}


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

// http://www.golangprograms.com/go-program-to-extracting-or-unzip-a-zip-format-file.html
func listFiles(file *zip.File, filename string, expression string, location string) error {
	fileread, err := file.Open()
	if err != nil {
		msg := "Failed to open zip %s for reading: %s"
		return fmt.Errorf(msg, file.Name, err)
	}
	defer fileread.Close()
 
 	if (strings.Contains(file.Name, expression)) {
 		// display zipfilename and contained file
		fmt.Fprintf(os.Stdout, "File extracted from %s : copied to -> %s/%s:", filename, location, file.Name) 	
	    fmt.Println()
	    
	    
	    zippedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer zippedFile.Close()

   		desktop := location +"/" + file.Name
   		
		
		
		outputFile, err := os.OpenFile(
		
				desktop,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				log.Fatal("*** error opening file ",err)
			}
	    
		defer outputFile.Close()


 
		_, err = io.Copy(outputFile, zippedFile)
		if err != nil {
			log.Fatal(err)
		}
    }

 
	if err != nil {
		msg := "Failed to read zip %s for reading: %s"
		return fmt.Errorf(msg, file.Name, err)
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


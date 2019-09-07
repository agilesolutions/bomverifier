package main
 
import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"net/http"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)


	type Bom struct {
    	Libs []struct {
            Name       string `yaml:"name"`
            Version    int    `yaml:"version"`
        } `yaml:"libs"`
	}
 
 
func main() {

    if len(os.Args) != 2 {
        fmt.Println("Usage:", os.Args[0], "uri")
        return
    }

    uri := os.Args[1]
    
    uri = "https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.yaml"
    
    var dst = "bom.yaml"

    fmt.Println("URI bom yaml file : ", uri )
    fmt.Println()
    
    fmt.Printf("DownloadToFile From: %s.\n", uri)
    if d, err := HTTPDownload(uri); err == nil {
        fmt.Printf("downloaded %s.\n", uri)
        if WriteFile("bom.yaml", d) == nil {
            fmt.Printf("saved %s as %s\n", uri, dst)
        }
    }

   filename, _ := filepath.Abs(dst)
    yamlFile, err := ioutil.ReadFile(filename)

    if err != nil {
        panic(err)
    }

	var bom Bom
	
	err = yaml.Unmarshal(yamlFile, &bom)
	if err != nil {
	    panic(err)
	}

	fmt.Print(bom.Libs[0].Name)
    

 	// filepath.Walk
 	//var inFile
 	files, err := FilePathWalkDir(".")
 	if err != nil {
  	panic(err)
 	}
 	for _, file := range files{
  		if (strings.HasSuffix(file, "jar")) {
	  		fmt.Println(file)
			read, err := zip.OpenReader(file )
			if err != nil {
				msg := "Failed to open: %s"
				log.Fatalf(msg, err)
			}
			defer read.Close()

			for _, infile := range read.File {
				
			
				if err := listFiles(infile, file, ".jar", bom); err != nil {
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
/*
* LIST FILES ***********************************************************
*/

func listFiles(file *zip.File, filename string, expression string, bom Bom) error {
	fileread, err := file.Open()
	if err != nil {
		msg := "Failed to open zip %s for reading: %s"
		return fmt.Errorf(msg, file.Name, err)
	}
	defer fileread.Close()
 
 	if (strings.Contains(file.Name, expression)) {
 		// display zipfilename and contained file
		//fmt.Fprintf(os.Stdout, "%s ", strings.Split(file.Name,"lib/")[1])
		fmt.Fprintf(os.Stdout, " frm origin %s ", bom.Libs[0].Name)
	    fmt.Println()
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


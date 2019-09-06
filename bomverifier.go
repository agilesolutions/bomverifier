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
)
 
 
func main() {

    if len(os.Args) != 2 {
        fmt.Println("Usage:", os.Args[0], "token")
        return
    }

    expression := os.Args[1]

    fmt.Println("searching for zip files with names containing : ", expression )
    fmt.Println()

	// https://stackoverflow.com/questions/16465705/how-to-handle-configuration-in-go
	type Configuration struct {
	    Copyfrom    string
	    Copyto		string
	}
	
	dir, _ := os.Getwd()
	fmt.Println(dir)

	file, _ := os.Open(dir +  "\\config.json")
	
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
  		fmt.Println("error reading config file : ", err)
	}
	fmt.Fprintf(os.Stdout, "reading PDF files from  -> %s:", configuration.Copyfrom) 	
	fmt.Println()
	fmt.Fprintf(os.Stdout, "writing PDF files to  -> %s:", configuration.Copyto) 	
	fmt.Println()
	fmt.Println()
	

    
    

 	// filepath.Walk
 	files, err := FilePathWalkDir(configuration.Copyfrom)
 	if err != nil {
  	panic(err)
 	}
 	for _, file := range files{
  		if (strings.HasSuffix(file, "zip")) {
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



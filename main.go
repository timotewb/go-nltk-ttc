package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ncruces/zenity"
)

func main(){
	inDir, err := zenity.SelectFile(
		zenity.Filename(""),
		zenity.Directory(),
		zenity.DisallowEmpty(),
		zenity.Title("Select input directory"),
	)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}

	// define vars
	currentTime := time.Now()
	outFileName := "output-"+currentTime.Format("2006-01-02_150405")+".txt"
	pattern := `output-\d{4}-\d{2}-\d{2}_\d{6}\.txt`
	csvFile, err := os.Create(filepath.Join(inDir,outFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	removeStrings := []string{"'","\"","“","”"}

	// get list of .txt file names
	files, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		matched, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			fmt.Println(err)
		}
		if filepath.Ext(file.Name()) == ".txt" && !matched{
			inFile, err := os.Open(filepath.Join(inDir,file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			scanner := bufio.NewScanner(inFile)
			for scanner.Scan(){
				sentences := strings.Split(scanner.Text(), ".")
				for _, sentence := range sentences {
					for _, removeString := range removeStrings {
						sentence = strings.ReplaceAll(sentence, removeString, "")
					}
					sentence = strings.TrimSpace(sentence)
					if sentence != "" && (regexp.MustCompile(`[A-Za-z]`).MatchString(sentence) || regexp.MustCompile(`\d`).MatchString(sentence)){
						err := writer.Write([]string{sentence})
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		}
	}

	zenity.Info("File created: "+outFileName,
		zenity.Title("Complete"),
		zenity.InfoIcon,
	)

}
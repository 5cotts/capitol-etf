package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/5cotts/capitol-etf/backend/models"
)

type MissingFileError struct {
	FileUrl string
}

func (e *MissingFileError) Error() string {
	return fmt.Sprintf("Missing file %s.", e.FileUrl)
}

func getFileNamesFromPath(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}

func parseFdrXml(filePath string) (*models.FinancialDisclosure, error) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer xmlFile.Close()

	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	fdr := models.FinancialDisclosure{}
	xml.Unmarshal(bytes, &fdr)
	return &fdr, nil
}

func downloadDisclosure(year int, DocId string, filePath string, dryRun bool) error {
	// TOOD(5cotts): It looks like some disclosures have different kinds of URLs.
	// See the below. Figure out how to account for different URLs.
	//
	// Type 1
	// https://disclosures-clerk.house.gov/public_disc/ptr-pdfs/{YEAR}/{DOCUMENT_ID}.pdf"
	//
	// Type 2
	// https://disclosures-clerk.house.gov/public_disc/financial-pdfs/{YEAR}/{DOCUMENT_ID}.pdf
	//
	url_type_1 := fmt.Sprintf(
		"https://disclosures-clerk.house.gov/public_disc/ptr-pdfs/%v/%s.pdf",
		year,
		DocId,
	)
	url_type_2 := fmt.Sprintf(
		"https://disclosures-clerk.house.gov/public_disc/ptr-pdfs/%v/%s.pdf",
		year,
		DocId,
	)

	urls := []string{url_type_1, url_type_2}
	for _, url := range urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		// Some documents do not exist even though they are in the XML.
		if res.StatusCode == 404 {
			return &MissingFileError{url}
		}

		if !dryRun {
			// Create a blank file
			file, err := os.Create(filePath)
			if err != nil {
				return err
			}

			defer file.Close()

			// Fill the blank file
			_, err = io.Copy(file, res.Body)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	var (
		totalXmlFiles = 0
		xmlFilePath   = "./fdr/xml/"
		totalPdfs     = 0
		missingPdfs   = 0
		dryRun        = true
	)

	xmlFileNames, err := getFileNamesFromPath(xmlFilePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, xmlFileName := range xmlFileNames {
		totalXmlFiles++

		xmlFilePath := fmt.Sprintf("%s%s", xmlFilePath, xmlFileName)
		fdr, err := parseFdrXml(xmlFilePath)
		if err != nil {
			log.Fatal(err)
		}

		numPdfs := len(fdr.Members)
		totalPdfs += numPdfs
		log.Println(
			"Will attempt to download",
			numPdfs,
			"financial disclosure reports from",
			xmlFilePath,
			"with dryRun =",
			dryRun,
		)

		for idx, member := range fdr.Members {
			err := downloadDisclosure(
				member.Year,
				member.DocId,
				fmt.Sprintf("./fdr/pdf/%s.pdf", member.DocId),
				dryRun,
			)

			switch err.(type) {
			case *MissingFileError:
				missingPdfs++
			}

			if idx%100 == 0 {
				log.Println(
					"On index:",
					idx,
					"Total PDFs:",
					missingPdfs,
					"Missing PDFs:",
					missingPdfs,
				)
				time.Sleep(5 * time.Second)
			}
		}
		log.Println("Finished downloading PDFs from", xmlFilePath)
	}
	log.Println("END", "Total PDFs:", totalPdfs, "Missing PDFs", missingPdfs)
}

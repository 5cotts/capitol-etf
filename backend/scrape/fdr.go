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
	url := fmt.Sprintf(
		"https://disclosures-clerk.house.gov/public_disc/ptr-pdfs/%v/%s.pdf",
		year,
		DocId,
	)

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

	return nil
}

func main() {
	var (
		missingFiles = 0
		dryRun       = true
	)

	file := "./fdr/xml/2021FD.xml"
	fdr, err := parseFdrXml(file)
	if err != nil {
		log.Fatal(err)
	}

	totalFiles := len(fdr.Members)
	log.Println(
		"Will attempt to download",
		totalFiles,
		"financial disclosure reports from",
		file,
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
			missingFiles++
		}

		if idx%100 == 0 {
			log.Println(
				"On index:",
				idx,
				"Total Files:",
				totalFiles,
				"Missing Files:",
				missingFiles,
			)
			time.Sleep(5 * time.Second)
		}
	}
	log.Println("END", "Total Files:", totalFiles, "Missing Files", missingFiles)
}

package csv

import (
	"mime/multipart"

	"github.com/gocarina/gocsv"
)

func ParseCSVFile(fileheader *multipart.FileHeader, output any) error {
	file, err := fileheader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gocsv.Unmarshal(file, output); err != nil { // Load clients from file
		return err
	}
	return nil
}

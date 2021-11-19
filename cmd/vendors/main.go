package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var directoryPtr = flag.String("directory", "", "the absolute path to the directory containing the accounting exports")

func main() {
	flag.Parse()

	if directoryPtr == nil {
		log.Fatalf("please provide the -directory flag")
	}
	vendors := make(map[string]bool)
	if err := filepath.Walk(*directoryPtr,
		func(path string, _ os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileVendors, err := checkFileForVendors(path)
			if err != nil {
				return err
			}
			for _, fileVendor := range fileVendors {
				vendors[fileVendor] = true
			}
			return nil
		}); err != nil {
		log.Fatal(err)
	}
	var unorderedVendors []string
	for vendor := range vendors {
		unorderedVendors = append(unorderedVendors, vendor)
	}
	sort.Strings(unorderedVendors)
	for _, vendor := range unorderedVendors {
		fmt.Println(vendor)
	}
}

func checkFileForVendors(path string) (vendors []string, err error) {
	if !strings.HasSuffix(path, ".csv") {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var vendorField *VendorField
	csvReader := csv.NewReader(f)
	for row := 0; true; row++ {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if row == 0 {
			for column, fieldName := range rec {
				switch {
				case fieldName == string(Description):
					vendorField = &VendorField{
						Index: column,
						Type:  Description,
					}
				case fieldName == string(VendorName):
					vendorField = &VendorField{
						Index: column,
						Type:  VendorName,
					}
				}
			}
			if vendorField == nil {
				return nil, nil
			}
		} else {
			vendors = append(vendors, vendorField.ParseVendor(rec))
		}
	}
	return vendors, nil
}

type VendorFieldType string

const (
	Description VendorFieldType = "Description"
	VendorName  VendorFieldType = "Vendor Name"
)

type VendorField struct {
	Index int
	Type  VendorFieldType
}

func (f *VendorField) ParseVendor(rec []string) string {
	value := rec[f.Index]
	switch {
	case f.Type == Description:
		return strings.Split(strings.ToLower(value), ",")[0]
	case f.Type == VendorName:
		return strings.ToLower(value)
	}
	return ""
}

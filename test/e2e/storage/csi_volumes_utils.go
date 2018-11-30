package storage

import (
	"os"
	"fmt"
	"encoding/json"
	"k8s.io/kubernetes/test/e2e/storage/testsuites/testdriver"
	"k8s.io/kubernetes/test/e2e/storage/utils"
)
type DriverInfo struct {
	Name       string // Name of the driver
	FeatureTag string // FeatureTag for the driver

	MaxFileSize          int64       // Max file size to be tested for this driver
	SupportedFsType      []string    // list of string for supported fs type
	SupportedMountOption []string    // list of string for supported mount option
	RequiredMountOption  []string    // list of string for required mount option (Optional)
	IsPersistent         bool        // Flag to represent whether it provides persistency
	IsFsGroupSupported   bool        // Flag to represent whether it supports fsGroup
	IsBlockSupported     bool        // Flag to represent whether it supports Block Volume

	Config testdriver.TestConfig // Test configuration for the current test.
}

type manifestJsonStruct struct {
	DriverInfo   DriverInfo
	PatchOptions utils.PatchCSIOptions
	Manifests    []string
	ScManifest   string
	ClaimSize    string
}

func getManifestJson(fileName string) (*manifestJsonStruct, error){
	var manifestObject manifestJsonStruct
	//Read Json File
	manifestFile, readErr := os.Open(fileName)
	defer manifestFile.Close()
	if readErr != nil{
		return nil, fmt.Errorf("Unable to read file: %v", readErr)
	}

	jsonParser := json.NewDecoder(manifestFile)
	jsonParser.Decode(&manifestObject)
	//output, err := json.MarshalIndent(manifestObject, "", "\t")
	//if err != nil {
	//	return &manifestObject, nil
	//}
	//ioutil.WriteFile("driver_manifest_read.json", output, 777)

	return &manifestObject, nil
}
package storage

import (
	"os"
	"fmt"
	"encoding/json"
	"k8s.io/kubernetes/test/e2e/storage/testsuites/testdriver"
	"k8s.io/kubernetes/test/e2e/storage/utils"
	)

type manifestJsonStruct struct {
	DriverInfo   testdriver.DriverInfo
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
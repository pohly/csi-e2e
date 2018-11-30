package storage

import (
	"encoding/json"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/storage/testsuites/testdriver"
	"k8s.io/kubernetes/test/e2e/storage/utils"
	"os"
	"path"
	"strings"
)

type DriverInfo struct {
	Name       string // Name of the driver
	FeatureTag string // FeatureTag for the driver

	MaxFileSize          int64    // Max file size to be tested for this driver
	SupportedFsType      []string // list of string for supported fs type
	SupportedMountOption []string // list of string for supported mount option
	RequiredMountOption  []string // list of string for required mount option (Optional)
	IsPersistent         bool     // Flag to represent whether it provides persistency
	IsFsGroupSupported   bool     // Flag to represent whether it supports fsGroup
	IsBlockSupported     bool     // Flag to represent whether it supports Block Volume

	Config testdriver.TestConfig // Test configuration for the current test.
}

type manifestJsonStruct struct {
	DriverInfo   DriverInfo
	PatchOptions utils.PatchCSIOptions
	Manifests    []string
	ScManifest   string
	ClaimSize    string
}

func getManifestJson(fileName string) (*manifestJsonStruct, error) {
	var manifestObject manifestJsonStruct
	//Read Json File
	manifestFile, readErr := os.Open(fileName)
	defer manifestFile.Close()
	if readErr != nil {
		return nil, fmt.Errorf("Unable to read file: %v", readErr)
	}

	jsonParser := json.NewDecoder(manifestFile)
	jsonParser.Decode(&manifestObject)

	return &manifestObject, nil
}

func PatchCSIDeployment(f *framework.Framework, o utils.PatchCSIOptions, object interface{}) error {
	rename := o.OldDriverName != "" && o.NewDriverName != "" &&
		o.OldDriverName != o.NewDriverName

	patchVolumes := func(volumes []v1.Volume) {
		if !rename {
			return
		}
		for i := range volumes {
			volume := &volumes[i]
			if volume.HostPath != nil {
				// Update paths like /var/lib/kubelet/plugins/<provisioner>.
				p := &volume.HostPath.Path
				dir, file := path.Split(*p)
				if file == o.OldDriverName {
					*p = path.Join(dir, o.NewDriverName)
				}
			}
		}
	}

	patchContainers := func(containers []v1.Container) {
		for i := range containers {
			container := &containers[i]
			if rename {
				for e := range container.Args {
					// Inject test-specific provider name into paths like this one:
					// --kubelet-registration-path=/var/lib/kubelet/plugins/csi-hostpath/csi.sock
					container.Args[e] = strings.Replace(container.Args[e], "/"+o.OldDriverName+"/", "/"+o.NewDriverName+"/", 1)
				}
			}
			// Overwrite driver name resp. provider name
			// by appending a parameter with the right
			// value.
			switch container.Name {
			//case o.DriverContainerName:
			//	container.Args = append(container.Args, "--drivername="+o.NewDriverName)
			case o.ProvisionerContainerName:
				// Driver name is expected to be the same
				// as the provisioner here.
				container.Args = append(container.Args, "--provisioner="+o.NewDriverName)
			}
		}
	}

	patchPodSpec := func(spec *v1.PodSpec) {
		patchContainers(spec.Containers)
		patchVolumes(spec.Volumes)
		if o.NodeName != "" {
			spec.NodeName = o.NodeName
		}
	}

	switch object := object.(type) {
	case *appsv1.ReplicaSet:
		patchPodSpec(&object.Spec.Template.Spec)
	case *appsv1.DaemonSet:
		patchPodSpec(&object.Spec.Template.Spec)
	case *appsv1.StatefulSet:
		patchPodSpec(&object.Spec.Template.Spec)
	case *appsv1.Deployment:
		patchPodSpec(&object.Spec.Template.Spec)
	case *storagev1.StorageClass:
		if o.NewDriverName != "" {
			// Driver name is expected to be the same
			// as the provisioner name here.
			object.Provisioner = o.NewDriverName
		}
	}

	return nil
}

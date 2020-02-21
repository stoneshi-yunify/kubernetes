/*
Copyright (C) 2018 Yunify, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this work except in compliance with the License.
You may obtain a copy of the License in the LICENSE file, or at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package qingcloud

import (
	v12 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/storage/testpatterns"
	"k8s.io/kubernetes/test/e2e/storage/testsuites"
)

type driver struct {
	driverInfo testsuites.DriverInfo
}

// initDriver returns driver that implements TestDriver interface
func initDriver(name string) testsuites.TestDriver {
	supportedTypes := sets.NewString(
		"", // Default fsType
	)
	return &driver{
		driverInfo: testsuites.DriverInfo{
			Name:                 name,
			MaxFileSize:          testpatterns.FileSizeLarge,
			SupportedFsType:      supportedTypes,
			SupportedMountOption: sets.NewString(),
			//SupportedSizeRange:   volume.SizeRange{Min: "5Gi"},
			Capabilities: map[testsuites.Capability]bool{
				testsuites.CapPersistence:         true,
				testsuites.CapBlock:               true,
				testsuites.CapFsGroup:             true,
				testsuites.CapExec:                true,
				testsuites.CapSnapshotDataSource:  true,
				testsuites.CapPVCDataSource:       true,
				testsuites.CapMultiPODs:           true,
				testsuites.CapControllerExpansion: true,
				testsuites.CapNodeExpansion:       true,
				testsuites.CapVolumeLimits:        true,
				testsuites.CapSingleNodeVolume:    true,
				//testsuites.CapTopology:            true,
			},
		},
	}
}

var _ testsuites.TestDriver = &driver{}
var _ testsuites.DynamicPVTestDriver = &driver{}

func (n *driver) GetDriverInfo() *testsuites.DriverInfo {
	return &n.driverInfo
}

func (n *driver) SkipUnsupportedTest(pattern testpatterns.TestPattern) {
}

func (n *driver) PrepareTest(f *framework.Framework) (*testsuites.PerTestConfig, func()) {
	config := &testsuites.PerTestConfig{
		Driver:    n,
		Prefix:    "qingcloud",
		Framework: f,
	}
	return config, func() {}
}

func (n *driver) GetDynamicProvisionStorageClass(config *testsuites.PerTestConfig, fsType string) *v12.StorageClass {
	parameters := map[string]string{
		"maxSize":  "5000",
		"minSize":  "100",
		"replica":  "2",
		"stepSize": "50",
		"type":     "2",
		"tags":     "tag-56odd6sg",
	}
	if fsType != "" {
		parameters["fsType"] = fsType
	}
	return testsuites.GetStorageClass(
		"disk.csi.qingcloud.com",
		parameters,
		nil,
		config.Framework.Namespace.Name,
		"qingcloud")
}

func (n *driver) GetClaimSize() string{
	return "5Gi"
}


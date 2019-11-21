package gounity

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dell/gounity/payloads"
	"strings"
	"testing"
	"time"
)

func TestCreateVolume(t *testing.T) {
	now := time.Now()
	volName := "test-" + now.Format("20060102150405")

	var vol *payloads.Volume
	var err error
	hostIOLimit, err := testConf.volumeApi.FindHostIOLimitByName(testConf.hostIOLimitName)
	fmt.Println("hostIOLimit:", prettyPrintJson(hostIOLimit), "Error:", err)

	if hostIOLimit != nil {
		//Ex JSon Request Body: {size: 5368709120,isThinEnabled: true,pool: {id: pool_1},fastVPParameters: {tieringPolicy: 0}}}
		vol, err = testConf.volumeApi.CreateLun(volName, testConf.poolId, "Description", 5368709120, 0, hostIOLimit.IoLimitPolicyContent.Id, true, false)
	} else {
		vol, err = testConf.volumeApi.CreateLun(volName, testConf.poolId, "Description", 5368709120, 0, "", true, false)
	}
	fmt.Println("Create volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Create volume failed: %v", err)
	}

	vol, err = testConf.volumeApi.FindVolumeByName(volName)
	fmt.Println("Find volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find volume failed: %v", err)
	}

	err = testConf.volumeApi.DeleteVolume(vol.VolumeContent.ResourceId)
	fmt.Println("Delete volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}

	//Ex JSon Request Body: {name:test-20190923105137,description:Description,lunParameters:{size:5368709120,isThinEnabled:false,pool:{id:pool_1},isDataReductionEnabled:false,fastVPParameters:{tieringPolicy:0}}}"
	fmt.Println("Test to verify thinEnabled false")
	vol, err = testConf.volumeApi.CreateLun(volName, testConf.poolId, "Description", 5368709120, 0, "", false, false)

	fmt.Println("Create volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Create volume failed: %v", err)
	}

	vol, err = testConf.volumeApi.FindVolumeByName(volName)
	fmt.Println("Find volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find volume failed: %v", err)
	}

	vols, _, err := testConf.volumeApi.ListVolumes(0, 10)
	fmt.Println("List volumes: ", len(vols))
	if len(vols) <= 10 {
		fmt.Println("List volume success")
	} else {
		t.Fatalf("Find volume failed: %v", err)
	}
	err = testConf.volumeApi.DeleteVolume(vol.VolumeContent.ResourceId)
	fmt.Println("Delete volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}
}
func TestCreateVolumeWithThinAndDataReduction(t *testing.T) {
	now := time.Now()
	volName := "test-" + now.Format("20060102150405")

	//var vol *payloads.Volume
	var err error
	err = testVerifyparameter(volName, true, true)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	} else {
		fmt.Println("Test success")
	}

	err = testVerifyparameter(volName, true, false)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	} else {
		fmt.Println("Test success")
	}

	err = testVerifyparameter(volName, false, true)
	if strings.ContainsAny(err.Error(), "Enable data reduction for the specified LUN also requires thin to be enabled.") {
		fmt.Println("Test success")
	} else {
		t.Fatalf("Test failed: %v", err)
	}

	err = testVerifyparameter(volName, false, false)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	} else {
		fmt.Println("Test success")
	}
}

func testVerifyparameter(volName string, thin, dataReduction bool) error {
	//Ex JSon Request Body: {name:test-20190923105137,description:Description,lunParameters:{size:5368709120,isThinEnabled:false,pool:{id:pool_1},isDataReductionEnabled:false,fastVPParameters:{tieringPolicy:0}}}"
	volName = fmt.Sprintf("%s%v%v", volName, thin, dataReduction)
	fmt.Printf("******Test to verify thin:%v and dataReduction:%v VolName: %s\n", thin, dataReduction, volName)
	vol, err := testConf.volumeApi.CreateLun(volName, testConf.poolId, "Description", 5368709120, 0, "", thin, dataReduction)

	str, _ := json.Marshal(vol)
	fmt.Println("*****Create volume:", string(str), err)
	if err != nil {
		return errors.New(fmt.Sprintf("Create volume failed for volume: %s error :%v", volName, err))
	}
	vol, err = testConf.volumeApi.FindVolumeByName(volName)
	fmt.Println("*****Find volume:", volName, prettyPrintJson(vol), err)
	if vol.VolumeContent.IsThinEnabled == thin && vol.VolumeContent.IsDataReductionEnabled == dataReduction {
		fmt.Println("Success: ", volName)
	} else {
		fmt.Println(vol.VolumeContent.IsThinEnabled, thin, vol.VolumeContent.IsDataReductionEnabled, dataReduction)
		return errors.New(fmt.Sprintf("Parametes are not matched.: %s", volName))
	}
	if err != nil {
		return errors.New(fmt.Sprintf("Find volume failed for volume: %s error :%v", volName, err))
	}

	err = testConf.volumeApi.DeleteVolume(vol.VolumeContent.ResourceId)
	fmt.Println("Delete volume:", prettyPrintJson(vol), err)
	return nil
}

func TestExpandVolume(t *testing.T) {
	now := time.Now()
	volName := "test-" + now.Format("20060102150405")

	var vol *payloads.Volume
	var err error

	vol, err = testConf.volumeApi.CreateLun(volName, testConf.poolId, "Description", 1368709120, 0, "", true, false)
	fmt.Println("Create volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Create volume failed: %v", err)
	}

	vol, err = testConf.volumeApi.FindVolumeByName(volName)
	fmt.Println("Find volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find volume failed: %v", err)
	}
	err = testConf.volumeApi.ExpandVolume(vol.VolumeContent.ResourceId, 2368709120)
	if err != nil {
		t.Fatalf("Expand volume failed: %v", err)
	}

	vol, err = testConf.volumeApi.FindVolumeByName(volName)
	fmt.Println("Find volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find volume failed: %v", err)
	}

	if vol.VolumeContent.SizeTotal != 2368709120 {
		t.Fatalf("Find volume failed: %v", err)
	}

	err = testConf.volumeApi.DeleteVolume(vol.VolumeContent.ResourceId)
	fmt.Println("Delete volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}
}

func TestExportVolume(t *testing.T) {
	now := time.Now()
	volName := "test-" + now.Format("20060102150405")

	var vol *payloads.Volume
	var host *payloads.Host
	var err error

	vol, err = testConf.volumeApi.CreateLun(volName, testConf.poolId, "Description", 5368709120, 0, "", true, false)
	fmt.Println("Create volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Create volume failed: %v", err)
	}

	vol, err = testConf.volumeApi.FindVolumeByName(volName)
	fmt.Println("Find volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find volume failed: %v", err)
	}

	host, err = testConf.hostApi.CreateHost(testConf.nodeHostName)
	fmt.Println("Create Host:", prettyPrintJson(host), err)
	if err != nil {
		//t.Fatalf("Create Host failed: %v", err)
	}

	host, err = testConf.hostApi.FindHostByName(testConf.nodeHostName)
	fmt.Println("Find Host:", prettyPrintJson(host), err)
	if err != nil {
		t.Fatalf("Find Host failed: %v", err)
	}

	err = testConf.volumeApi.ExportVolume(vol.VolumeContent.ResourceId, host.HostContent.ID)
	fmt.Println("ExportVolume:", prettyPrintJson(host), err)
	if err != nil {
		t.Fatalf("ExportVolume failed: %v", err)
	}
	err = testConf.volumeApi.UnexportVolume(vol.VolumeContent.ResourceId)
	fmt.Println("UnExportVolume:", prettyPrintJson(host), err)
	if err != nil {
		t.Fatalf("UnExportVolume failed: %v", err)
	}
	err = testConf.volumeApi.DeleteVolume(vol.VolumeContent.ResourceId)
	fmt.Println("Delete volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}
	err = testConf.hostApi.DeleteHost(testConf.nodeHostName)
	fmt.Println("Delete Host:", err)
	if err != nil {
		t.Fatalf("Delete Host failed: %v", err)
	}
}

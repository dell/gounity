package gounity

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

type testConfig struct {
	unityEndPoint   string
	username        string
	password        string
	poolId          string
	nodeHostName    string
	nodeHostIp      string
	wwns            []string
	iqn             string
	hostIOLimitName string
	nasServer       string
	volumeApi       *volume
	hostApi         *host
	poolApi         *storagepool
	snapApi         *snapshot
	ipinterfaceApi  *ipinterface
	fileApi         *filesystem
	metricsApi      *metrics
}

var testConf *testConfig

func TestMain(m *testing.M) {
	fmt.Println("------------In TestMain--------------")
	os.Setenv("GOUNITY_DEBUG", "true")

	// for this tutorial, we will hard code it to config.txt
	testProp, err := readTestProperties("test.properties")
	if err != nil {
		panic("The system cannot find the file specified")
	}

	insecure, _ := strconv.ParseBool(testProp["X_CSI_UNITY_INSECURE"])
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}

	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()

	testConf = &testConfig{}
	testConf.unityEndPoint = testProp["GOUNITY_ENDPOINT"]
	testConf.username = testProp["X_CSI_UNITY_USER"]
	testConf.password = testProp["X_CSI_UNITY_PASSWORD"]
	testConf.poolId = testProp["STORAGE_POOL"]
	testConf.nodeHostName = testProp["NODE_HOSTNAME"]
	testConf.hostIOLimitName = testProp["HOST_IO_LIMIT_NAME"]
	testConf.nodeHostIp = testProp["NODE_HOSTIP"]
	testConf.nasServer = testProp["UNITY_NAS_SERVER"]
	testConf.iqn = testProp["NODE_IQN"]
	wwnStr := testProp["NODE_WWNS"]

	os.Setenv("GOUNITY_ENDPOINT", testConf.unityEndPoint)
	os.Setenv("X_CSI_UNITY_USER", testConf.username)
	os.Setenv("X_CSI_UNITY_PASSWORD", testConf.password)

	testConf.username = testProp["X_CSI_UNITY_USER"]
	testConf.password = testProp["X_CSI_UNITY_PASSWORD"]

	testClient := getTestClient(ctx, testConf.unityEndPoint, testConf.username, testConf.password, testConf.unityEndPoint, insecure)
	testConf.wwns = strings.Split(wwnStr, ",")

	testConf.hostApi = NewHost(testClient)
	testConf.poolApi = NewStoragePool(testClient)
	testConf.snapApi = NewSnapshot(testClient)
	testConf.volumeApi = NewVolume(testClient)
	testConf.ipinterfaceApi = NewIpInterface(testClient)
	testConf.fileApi = NewFilesystem(testClient)
	testConf.metricsApi = NewMetrics(testClient)

	code := m.Run()
	fmt.Println("------------End of TestMain--------------")
	os.Exit(code)
}

func getTestClient(ctx context.Context, url, username, password, endpoint string, insecure bool) *Client {
	fmt.Println("Test:", url, username, password)

	c, err := NewClientWithArgs(ctx, endpoint, insecure)
	if err != nil {
		fmt.Println(err)
	}

	err = c.Authenticate(ctx, &ConfigConnect{
		Username: username,
		Password: password,
		Endpoint: url,
	})
	return c
}

func readTestProperties(filename string) (map[string]string, error) {
	// init with some bogus data
	configPropertiesMap := map[string]string{}
	if len(filename) == 0 {
		return nil, errors.New("Error reading properties file " + filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		// check if the line has = sign
		// and process the line. Ignore the rest.
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				// assign the config map
				configPropertiesMap[key] = value
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return configPropertiesMap, nil
}

func prettyPrintJson(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}

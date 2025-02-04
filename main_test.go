/*
 Copyright © 2019-2025 Dell Inc. or its subsidiaries. All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
      http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package gounity

import (
	"bufio"
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

	"github.com/dell/gounity/mocks"
)

type testConfig struct {
	unityEndPoint   string
	username        string
	password        string
	poolID          string
	nodeHostName    string
	nodeHostIP      string
	wwns            []string
	iqn             string
	hostIOLimitName string
	nasServer       string
	tenant          string
	hostList        []string
	client          UnityClient
}

var testConf *testConfig

func TestMain(m *testing.M) {
	fmt.Println("------------In TestMain--------------")
	os.Setenv("GOUNITY_DEBUG", "true")

	// for this tutorial, we will hard code it to config.txt
	testProp, err := readTestProperties("test.properties_template")
	if err != nil {
		panic("The system cannot find the file specified")
	}

	insecure, _ := strconv.ParseBool(testProp["X_CSI_UNITY_INSECURE"])
	/* #nosec G402 */
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}

	if err != nil {
		fmt.Println(err)
	}

	testConf = &testConfig{}
	testConf.unityEndPoint = "https://mock-endpoint"
	testConf.username = "user"
	testConf.password = "password"
	testConf.poolID = "pool_3"
	testConf.nodeHostName = "Unit-test-host-20231023052923"
	testConf.hostIOLimitName = "Autotyre"
	testConf.nodeHostIP = "10.20.30.40"
	testConf.nasServer = "nas_1"
	testConf.iqn = "iqn.1996-04.de.suse:01:f8298e544dc"
	wwnStr := ""
	hostListStr := "Unit-test-host-20231023052923"
	testConf.tenant = "tenant_1"

	os.Setenv("GOUNITY_ENDPOINT", testConf.unityEndPoint)
	os.Setenv("X_CSI_UNITY_USER", testConf.username)
	os.Setenv("X_CSI_UNITY_PASSWORD", testConf.password)

	testConf.client = getTestClient()
	testConf.wwns = strings.Split(wwnStr, ",")
	testConf.hostList = strings.Split(hostListStr, ",")

	code := m.Run()
	fmt.Println("------------End of TestMain--------------")
	os.Exit(code)
}

func getTestClient() *UnityClientImpl {
	return &UnityClientImpl{
		api:           &mocks.Client{},
		configConnect: &ConfigConnect{},
	}
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

func prettyPrintJSON(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}

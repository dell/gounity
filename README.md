<!--
 Copyright © 2019 Dell Inc. or its subsidiaries. All Rights Reserved.
 
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
      http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
-->
# GoUnity
A portable Go library which represents API bindings that allow you to manage Unity XT storage platforms.

## Integration Tests Execution
Follow the steps to run integration tests:
1. Create a properties file `test.properties` using the template `test.properties_template`.
2. Populate the properties file with values from your Unity XT system as shown below: 

![test prop1](https://user-images.githubusercontent.com/92028646/161742532-bafc1927-4cbe-4b10-ab7a-d671d883d493.JPG) 

3. To run the integration tests, run `make go-unittest`. Once all the tests in each module are run successfully, you will see `Output` as `PASS` for each of the module else `Output` is `FAIL`.
4. To get the integration test coverage for each module, run `make go-coverage`.
5. To generate and analyze coverage statistics, run `go tool cover -html=gounity_coverprofile.out`.

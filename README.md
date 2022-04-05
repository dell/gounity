# gounity
A portable Go library for unity related operations such as create-volume, list-volume, etc.__

## Integration Tests execution
Follow the steps to run integration tests:
1. Create a properties file `test.properties` using the template `test.properties_template`.
2. Populate the properties file with values for your DellEMC Unity system as shown below: 

  ![test prop](https://user-images.githubusercontent.com/92028646/161708153-f3e1fb38-2e2d-4ca0-9796-fc61a181cb9f.JPG)

3. To run the integration tests, run `make go-unittest`. After all the tests in each module run successfully, you will see `Output` as `PASS` for the module else `Output` is `FAIL`.
4. To get the integration test coverage for each module, run `make go-coverage`.
5. To generate and analyze coverage statistics, run `go tool cover -html=gounity_coverprofile.out`.

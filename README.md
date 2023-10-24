# proxy-svc

## Description

Proxy service is the module that allows rarimo front-ends send requests to the third party APIs. 

## Install

  ```bash
  git clone githut.com/rarimo/proxy-svc
  cd main
  go build main.go
  export KV_VIPER_FILE=./config.yaml
  ./main run api
  ```

## Running from Source

* Set up environment value with config file path `KV_VIPER_FILE=./config.yaml`
* Provide valid config file
* Launch the service with `run api` command

## License
[MIT](./LICENSE)

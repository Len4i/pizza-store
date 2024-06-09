# Pizza store order api
> [!Warning]
> This is the demo project.
><details>
>  <summary><i>Details</i></summary>
>
>1. Write a server  
>Should have the following endpoints:  
>`/order` ( POST )  
>`/health` ( GET ) (this is also a hint )  
>`/order` receives a POST with the body:
>   ```json
>   {
>       "pizza-type": "<margherita|pugliese|marinara>"
>       "size": "<personal|family>"
>       "amount": <int>
>   }
>   ```
>2. Write a Secure Dockerfile
>3. Write a CI pipeline
>
></details>

![release badge](https://github.com/Len4i/pizza-store/actions/workflows/release.yaml/badge.svg?event=push)

## Helm chart installation
```bash
helm upgrade <my-release> --install -n pizza-store --create-namespace oci://registry-1.docker.io/len4i/pizza-store-helm --version <chart-version>
```
List of versions can be found in the [dockerhub repo](https://hub.docker.com/r/len4i/pizza-store-helm/tags)
Tag of the image and helm chart version are the same  

Default securityContext:
```yaml
securityContext:
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 1000
```
> [!NOTE]
> Helm chart was created with `helm create` and has minimum changes from the default  
> In current state volumes and volumeMounts are set in the [values.yaml](./helm-chart/values.yaml)  
> While chart working with default values, it can be easily broken as require alignment between the values and parts of the template  
> For example `volumeMount` path of the db is set in the values but also as part of the configmap. both values should be aligned


## Production readiness

### Code
There are several things that have to be handled in order for this service be more production ready  
Few examples: 
- support TLS on the server
- HA database (more details in the [storage](#storage) section)
- error handling
- request validation
- api functionality
    - full CRUD capabilities
    - different table schema to support search and updates of the orders
    - api docs (for example swagger)
- proper handling of http status codes in the `httpLogger`

### Tests
For now there is only a simple [test of storage functions](./internal/storage/sqlite/sqlite_test.go)  
Things todo:
- generate mock for interface using [Mockery](https://github.com/vektra/mockery)
- add tests for `order.go` using mocked interface and `net/http/httptest` package
- more tests cases around database

### Infra
- More work on a Helm chart. _details in the [helm chart](#helm-chart-installation) section_
- add exmaples for gitops deployment (argocd / flux)
- build image with docker buildx for multiple architectures including arm


## Server
Used [Chi](https://github.com/go-chi/chi) as a web server

### Configuration
[Example](./config-example.yaml)
```yaml
storage_path: /tmp/db.sql
http_server:
  address: :8080
  read_header_timeout: 4s
  idle_timeout: 60s
graceful_shutdown_timeout: 30s
log_level: 0 # 0: info, -4: debug
```
For now only `info` and `debug` are used in logs. you can find mapping of the `log_level` values in the [slog docs](https://pkg.go.dev/log/slog#Level)

### Logger
Using 2 loggers:
1. standard slog in json format for the app logs
2. slog based logger for http requests with middleware that enriches log with http request data  
_http logger servers the only puprose to align loggin format for access logs, so that it can be collected same way as an app logs_

### Storage
Using [modernc.org/sqlite](modernc.org/sqlite)   
_It isn't the most popular implementation, but it is written in pure go, so that there is no need to deal with `cgo` dependencies at the time of the build_  
SQL lite is not designed for `HA` applications as it is writing to a local file    
In case of need should be replaces with other `database/sql` compatible driver   
> Another way is to swap the whole [storage](./internal/storage/sqlite/sqlite.go) implementation, preserving the implementation of the `Order` interface  



## Github action
[release](.github/workflows/release.yaml) workflow actions:
- create new tag for the release
- building image and pushing it to the public dockerhub repo
- packaging helm chart and pushing it to OCI registry (same dockerhub, different repository)
- create github release with chart as an artifact



## Local dev
1. 
```bash
git clone git@github.com:Len4i/pizza-store.git
```
2. 
```bash
cd pizza-store
```
3. 
```bash
go mod tidy
```
4.
```bash
export CONFIG_PATH=config-example.yaml
```
5. Make sure that you have write permissions to `storage_path` folder

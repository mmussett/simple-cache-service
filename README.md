# simple-cache-service

A simple in-memory name-value cache service

## Building from source

### Using godep to build

If you haven't got godeps installed, do:
```bash
go get github.com/tools/godep
```

then

```bash
godep restore
go build -o simple-cache-service
```

## Dockerizing

Docker script provided will build using container version of go

### Building

docker build -t mmussett/golang-simple-cache-service .

### Running

docker run -d -p 8080:8080 --name scs mmussett/golang-simple-cache-service


## Cache JSON object format

```json
{
  "id" : "<<cache key>>",
  "value" : "<<cache value>>"
}
```


## Getting from the cache

curl -vk -X GET 'http://localhost:8080/cache/1'

## Adding to the cache

curl -vk -X POST -d '{"id":"1","value":"hello"}' 'http://localhost:8080/cache'

## Updating the cache

curl -vk -X PUT -d '{"id":"1","value":"hello,world"}' 'http://localhost:8080/cache/1'

## Removing from the cache

curl -vk -X DELETE 'http://localhost:8080/cache/1'




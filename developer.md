


## Build for linux
`env GOOS=linux GOARCH=amd64 go build`

## Build for MacOs
`go build`


## Test

`go test -coverprofile=coverage.out`

## Show test coverge

`go tool cover -html=coverage.out`



## run the memcache server
`docker run -p 11211:11211 memcached:alpine`



## To do pre-commit
`ln -s $(pwd)/hooks/pre-commit .git/hooks/pre-commit`

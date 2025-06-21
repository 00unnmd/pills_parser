## pills-parser backend service

### Requirements
```
go v1.24.3 
psql v17.5
```
also need Chrome browser app for `chromedp` package work

### Run Parse
```
go run cmd/parser/main.go
```
you can configure pills and regions to request from `/pkg/utils/info.go`

### Run Server
```
go run cmd/server/main.go
```

### Build Server for deploy
```
env GOOS=linux GOARCH=amd64 go build -o app cmd/server/main.go
```
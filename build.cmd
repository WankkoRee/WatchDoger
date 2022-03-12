SET VERSION=v1.0.0

SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/client src/client/main.go
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/server src/server/main.go
tar -C build/%VERSION%/%GOOS%/%GOARCH%/ -czvf build/WatchDoger_%VERSION%_%GOOS%_%GOARCH%.tar.gz *

SET GOOS=linux
SET GOARCH=386
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/client src/client/main.go
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/server src/server/main.go
tar -C build/%VERSION%/%GOOS%/%GOARCH%/ -czvf build/WatchDoger_%VERSION%_%GOOS%_%GOARCH%.tar.gz *

SET GOOS=linux
SET GOARCH=arm64
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/client src/client/main.go
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/server src/server/main.go
tar -C build/%VERSION%/%GOOS%/%GOARCH%/ -czvf build/WatchDoger_%VERSION%_%GOOS%_%GOARCH%.tar.gz *

SET GOOS=linux
SET GOARCH=arm
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/client src/client/main.go
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/server src/server/main.go
tar -C build/%VERSION%/%GOOS%/%GOARCH%/ -czvf build/WatchDoger_%VERSION%_%GOOS%_%GOARCH%.tar.gz *

SET GOOS=windows
SET GOARCH=amd64
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/client.exe src/client/main.go
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/server.exe src/server/main.go
tar -C build/%VERSION%/%GOOS%/%GOARCH%/ -czvf build/WatchDoger_%VERSION%_%GOOS%_%GOARCH%.tar.gz *

SET GOOS=windows
SET GOARCH=386
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/client.exe src/client/main.go
go build -ldflags "-w -s -X main.version=%VERSION%" -o build/%VERSION%/%GOOS%/%GOARCH%/server.exe src/server/main.go
tar -C build/%VERSION%/%GOOS%/%GOARCH%/ -czvf build/WatchDoger_%VERSION%_%GOOS%_%GOARCH%.tar.gz *

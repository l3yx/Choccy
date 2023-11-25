cd web
npm install
npm run build
cd ..
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o choccy_linux_amd64 main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o choccy_windows_amd64.exe main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o choccy_darwin_amd64 main.go
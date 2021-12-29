set GOOS=windows
go build -o murphysec-windows.exe -tags noprint,idea .
set GOOS=linux
go build -o murphysec-linux -tags noprint,idea
set GOOS=darwin
go build -o murphysec-macos -tags noprint,idea .
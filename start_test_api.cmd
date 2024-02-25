cd tests
del tests_api.exe /q
go build tests_api.go
tests_api.exe

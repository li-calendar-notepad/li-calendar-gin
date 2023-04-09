cd web 
call npm run build 
rd ..\assets\frontend
xcopy .\dist\ ..\assets\frontend /s /i /y
cd ..\
go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/...
go build -o li-calendar.exe --ldflags="-X main.RunMode=release" main.go
echo "li-calendar.exe compilation is complete!"cd web 
npm run build 
rd ..\assets\frontend
xcopy .\dist\ ..\assets\frontend /s /i /y
cd ..\
go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/...
go build -o li-calendar.exe --ldflags="-X main.RunMode=release" main.go
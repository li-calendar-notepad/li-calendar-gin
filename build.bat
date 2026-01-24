@REM 请提前搭建好node.js环境和go环境
@REM 请将前端项目文件拉取到本目录中，并将项目文件夹重命名为web

set filenameZip=li-calendar_windows.zip
set filename=li-calendar.exe
cd web 
call npm run build 
rd ..\assets\frontend
xcopy .\dist\ ..\assets\frontend /s /i /y
cd ..\
go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/...
go build -o %filename% --ldflags="-X main.RunMode=release" main.go


REM 打包文件
powershell Compress-Archive %filename% %filenameZip%

REM 删除原始文件
del %filename%

REM 返回上一级目录，进入父文件夹
cd ..
echo "%filenameZip% compilation is complete!"
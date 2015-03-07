set NEWHOME=%HOME%
set NEWHOME=%NEWHOME:/=\%

MKDIR %NEWHOME%\site\wwwroot\gopath
MKDIR %NEWHOME%\site\wwwroot\gopath\bin
MKDIR %NEWHOME%\site\wwwroot\gopath\src

set hr=%time:~0,2%
if "%hr:~0,1%" equ " " set hr=0%hr:~1,1%
SET FOLDER=%date:~-4,4%%date:~-10,2%%date:~-7,2%_%hr%%time:~3,2%%time:~6,2%

MKDIR %NEWHOME%\site\wwwroot\gopath\src\%FOLDER%

xcopy %DEPLOYMENT_SOURCE% %NEWHOME%\site\wwwroot\gopath\src\%FOLDER% /Y
xcopy /r %DEPLOYMENT_SOURCE%\Web.Config %NEWHOME%\site\wwwroot\Web.Config /Y

xcopy %DEPLOYMENT_SOURCE%\content\* %NEWHOME%\site\wwwroot\content\ /S /Y
xcopy %DEPLOYMENT_SOURCE%\views\* %NEWHOME%\site\wwwroot\views\ /S /Y

SET GOPATH=%NEWHOME%\site\wwwroot\gopath
SET GOROOT=%NEWHOME%\site\wwwroot\go
SET PATH=%PATH%;%GOPATH%\bin;%NEWHOME%\site\wwwroot\go\bin

go get %FOLDER%

powershell "stop-process (Get-Process w3wp | Sort-Object ws | Select -first 1).Id"

DEL %NEWSITE%\site\wwwroot\app.exe
xcopy %NEWHOME%\site\wwwroot\gopath\bin\%FOLDER%.exe %NEWHOME%\site\wwwroot\ /Y
REN %NEWSITE%\site\wwwroot\%FOLDER%.exe app.exe

REM sed -i 's/GOAPPBINARY/%FOLDER%.exe/g' %NEWHOME%\site\wwwroot\Web.Config

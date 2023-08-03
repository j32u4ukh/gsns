@ECHO ON
CD ../cmd
SET ROOT=%CD%

cd dba
go build
cd ../account
go build
cd ../post_message
go build
cd ../gsns
go build
PAUSE
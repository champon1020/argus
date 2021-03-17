module github.com/champon1020/argus

go 1.16

replace github.com/champon1020/mgorm => ../argus/mgorm

require (
	github.com/champon1020/mgorm v0.0.0-20210316114102-2461c3545a03
	github.com/go-sql-driver/mysql v1.5.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/echo/v4 v4.2.1
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/sys v0.0.0-20210317225723-c4fcb01b228e // indirect
)

# gwmpd
MPD Web GUI - written in Go (multithread)

There are 2 parts:
* the backend which written in Go. It communicates with mpd and web GUI
* the frontend which written in JS (vuejs)

__Be careful. Actually, there is NO security between gui and the REST API so make local test__
__DO NOT expose REST API on internet !__

## Dependencies
* go
* gin
* cors (middleware)
* viper
* go-logging
* vuejs

## Install
```sh
go get github.com/gin-gonic/gin
go get github.com/itsjamie/gin-cors
go get github.com/op/go-logging
go get github.com/spf13/viper

git clone https://github.com/Chipsterjulien/gwmpd.git
```
Inside first console:
```sh
cd gwmpd_back
go run app.go initLogging.go loadConfig.go
```

Inside second console:
```sh
cd gwmpd_front/Gwmpd
yarn
yarn start
```

Start your browser to [http://localhost:8080](http://localhost:8080)

## Config

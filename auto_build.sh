#! /bin/sh

mkdir -p back
mkdir -p front

# build backend
go get github.com/gin-gonic/gin
go get github.com/itsjamie/gin-cors
go get github.com/op/go-logging
go get github.com/spf13/viper
go get github.com/appleboy/gin-jwt

cd gwmpd_back/
go build -o gwmpdBack
mv gwmpdBack ../back
mv cfg ../back
cd ..

# build frontend
cd gwmpd_front/Gwmpd/
yarn
yarn build
cd dist/
mv * ../../../front
cd ../../..
rm -rf gwmpd_back gwmpd_front
rm License README.md auto_build.sh

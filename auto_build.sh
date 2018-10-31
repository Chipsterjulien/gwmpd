#! /bin/sh

mkdir -p back front

# build backend
echo "Install/Update go packages"
go get -u github.com/gin-gonic/gin
go get -u github.com/itsjamie/gin-cors
go get -u github.com/op/go-logging
go get -u github.com/spf13/viper
go get -u github.com/appleboy/gin-jwt

# build backend
echo "Build backend"
cd gwmpd_back/
go build -o gwmpdBack
mv gwmpdBack ../back
mv cfg ../back
cd ..

# build frontend
echo "Build frontend"
cd gwmpd_front/Gwmpd/
yarn
yarn build
cd dist/
mv * ../../../front
cd ../../..
rm -rf gwmpd_back gwmpd_front
rm License README.md auto_build.sh

echo "End !"

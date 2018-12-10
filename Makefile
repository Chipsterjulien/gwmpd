init:
	go get -u github.com/gin-gonic/gin
	go get -u github.com/itsjamie/gin-cors
	go get -u github.com/op/go-logging
	go get -u github.com/spf13/viper
	go get -u github.com/appleboy/gin-jwt

	cd ./gwmpd_front/Gwmpd/ && yarn && cd ./

run:
	cd ./gwmpd_front/Gwmpd/ && gnome-terminal -e yarn start
	cd ./gwmpd_back/ && go build -o gwmpd_back && gnome-terminal -e ./gwmpd_back
	cd ./

build:
	cd ./gwmpd_front/Gwmpd/ && yarn build
	cp -R ./gwmpd_front/Gwmpd/dist/* ./Build/Frontend/
	cd ./gwmpd_back/ && go build -o ../Build/Backend/gwmpd_x86-64
	cd ./gwmpd_back/ && env GOOS=linux GOARCH=arm GOARM=6 go build -o ../Build/Backend/gwmpd_generic_arm_or_v6
	cd ./gwmpd_back/ && env GOOS=linux GOARCH=arm GOARM=6 go build -o ../Build/Backend/gwmpd_arm_or_v7
	cd ./gwmpd_back/ && mv * ./Build/Backend/
init:
	# Update/install backend dependancies
	go get -u github.com/gin-gonic/gin
	go get -u github.com/itsjamie/gin-cors
	go get -u github.com/op/go-logging
	go get -u github.com/spf13/viper
	go get -u github.com/appleboy/gin-jwt

	# Update/install frontend dependancies
	cd ./gwmpd_front/Gwmpd/ && yarn && cd ./

run:
	cd ./gwmpd_front/Gwmpd/ && xfce4-terminal -H -x yarn start &

	cd ./gwmpd_back/ && ./clean.sh
	cd ./gwmpd_back/ && python changePath.py dev
	cd ./gwmpd_back/ && go build -o gwmpd_back && xfce4-terminal -H -x ./gwmpd_back &
	cd ./

build:
	# Build frontend
	cd ./gwmpd_front/Gwmpd/ && yarn build
	# Build backend
	cd ./gwmpd_back/ && python changePath.py prod
	cd ./gwmpd_back/ && go build -o ../Build/Backend/gwmpd_x86-64
	cd ./gwmpd_back/ && env GOOS=linux GOARCH=arm GOARM=6 go build -o ../Build/Backend/gwmpd_generic_arm_or_v6
	cd ./gwmpd_back/ && env GOOS=linux GOARCH=arm GOARM=6 go build -o ../Build/Backend/gwmpd_arm_or_v7

	# Move files
	cp -R ./gwmpd_front/Gwmpd/dist/* ./Build/Frontend/

	# Remove useless files/directories
	rm -rf ./gwmpd_front/Gwmpd/dist
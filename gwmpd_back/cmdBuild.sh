echo "Build x86-64"
go build -o ../Build/Backend/gwmpd_x86-64
echo "Build arm-v6"
env GOOS=linux GOARCH=arm GOARM=6 go build -o ../Build/Backend/gwmpd_generic_arm_or_v6
echo "Build arm-v7"
env GOOS=linux GOARCH=arm GOARM=6 go build -o ../Build/Backend/gwmpd_arm_or_v7
echo "end"
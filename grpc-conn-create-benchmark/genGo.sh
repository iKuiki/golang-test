#!/bin/bash
source genPaths.sh

GOPATH=`go env GOPATH`

function genGrpc(){
	echo "TARGET_GRPC_PATH:	$1"
	if [ $GOPATH ];then
		if [ ! -d $GOPATH/src/$1 ];then
			mkdir -p $GOPATH/src/$1
		fi
		cd $GOPATH/src && protoc --go_out=plugins=grpc:./ $1/*.proto

		if [ -f $GOPATH/src/$1/optimize.sh ];then
			cd $GOPATH/src/$1 && ./optimize.sh
		fi
	else
		echo "gen pb.go Fail: GOPATH not exist or empty!"
	fi
}

function genPbf(){
	echo "TARGET_PBF_PATH:	$PPATH"
	if [ $GOPATH ];then
		if [ ! -d $GOPATH/src/$1 ];then
			mkdir -p $GOPATH/src/$1
		fi
		cd $GOPATH/src && protoc --go_out=./ $1/*.proto

		if [ -f $GOPATH/src/$1/optimize.sh ];then
			cd $GOPATH/src/$1 && ./optimize.sh
		fi
	else
		echo "gen pb.go Fail: GOPATH not exist or empty!"
	fi
}

for GPATH in ${GRPC_PATHS[*]}
do
	genGrpc $GPATH
done

for PPATH in ${PBF_PATHS[*]}
do
	genPbf $PPATH
done

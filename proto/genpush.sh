#!/bin/sh
git pull
#
#CRTDIR=$(pwd)
#echo ${CRTDIR}
#cd gocn
#dirs=`find . -type d`
#for dir in $dirs; do
#    echo $dir
#    # 只处理存在proto文件的目录
#    if [ "$(ls ${dir}/*.proto 2>/dev/null)" ]; then
##        echo "protoc --go_out=plugins=grpc:gen/go $dir/*.proto"
##        protoc --go_out=plugins=grpc:../../../ -I . $dir/*.proto;
##        protoc --go_out=.  --go-grpc_out=.  -I . $dir/*.proto;
#        protoc --go_out=${CRTDIR}/gocn/gen --go_opt=paths=source_relative --go-grpc_out=${CRTDIR}/gocn/gen --go-grpc_opt=paths=source_relative -I . $dir/*.proto;
#
#        # 生成python桩代码
#        # echo "python -m grpc_tools.protoc -I$dir/ --python_out=$dir/ --grpc_python_out=$dir/ $dir/*.proto"
#	      # python -m grpc_tools.protoc -I./ --python_out=./ --grpc_python_out=./  $dir/*.proto;
#    fi
#done

egoctl pb generate

git add *
git commit -am "proto"
git push
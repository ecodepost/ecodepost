#!/bin/bash

# 将此目录拷贝到临时目录执行
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
TMP=/tmp/proto-tmp
OUT="jspb"
GITHOST="git.yitum.com"
GITPATH_PREFIX="of"
rm -rf ${TMP} && cp -r ${ROOT} ${TMP} && cd ${TMP} && rm -rf ./${OUT} && cd -

if [[ ${MODE} == "" ]]; then
 MODE="local"
fi
echo "CI_COMMIT_BEFORE_SHA:"${CI_COMMIT_BEFORE_SHA}
echo "CI_COMMIT_SHA:"${CI_COMMIT_SHA}
echo "CUSTOM_DIR:"${CUSTOM_DIR}
echo "MODE:"${MODE}

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)
        arch=linux
        ;;
    Darwin*)
        arch=osx;
        shopt -s expand_aliases
        alias grep='ggrep'
        ;;
    *)
        exit 1
        ;;
esac

# 生成go.mod文件
function genGoMod() {
    echo pwd:$(pwd)
    # subpackages获取需要生成go.mod的文件夹
    subpackages=`find . -depth -type d | awk -F / 'NF>=p; {p=NF}' |\
       grep -v "\./\." | grep -v $OUT | grep -v "demo" | cut -c 3-`

    # 为每个子包单独生成go.mod
    for subpackage in ${subpackages}; do
        if [[ $subpackage == *"/v1" ]]; then
            local modulePath=`echo $subpackage | sed -e "s~/v1~~g"`
            echo "module ${GITHOST}/${GITPATH_PREFIX}/$OUT/$modulePath" >> "$OUT/$modulePath/go.mod"
        else
            local modulePath=$subpackage
            echo "module ${GITHOST}/${GITPATH_PREFIX}/$OUT/$modulePath" >> "$OUT/$modulePath/go.mod"
        fi
    done
}

# 构建target
function build() {
    #egoctl pb generate
    mkdir -p ./${OUT}/common/v1 && protoc --js_out=./${OUT}/common/v1 ./common/v1/enum_*.proto
    # genGoMod
    echo "[proto] generate success!"
}

# 推送至Git仓库
function push() {
    repo=https://${USER}:${PSWD}@git.yitum.com/of/${OUT}.git
    repoPath=${TMP}/${OUT}
    workDir=$(pwd)
    rm -rf ${repoPath}
    git clone ${repo} ${repoPath}
    ls -lh ${repoPath}

    cd ${repoPath}
    echo  ${repoPath}
    git config user.email "infra-bot@yitum.im"
    git config user.name "infra-bot"

    git fetch --all
    git checkout ${CI_COMMIT_REF_NAME}
    ls -rtlh ${workDir}
    cp -r ${workDir}/${OUT}/* ${repoPath}
    ls -rtlh ${workDir}/${OUT}/*
    git branch
    git diff --name-only
    git add .
    git commit -m "build: [${CI_COMMIT_REF_NAME}](https://git.yitum.com/of/proto/-/commits/${CI_COMMIT_SHORT_SHA})"
    git push -u origin ${CI_COMMIT_REF_NAME}
    echo -e "\n[proto] push success!"
}

build

if [[ ${MODE} == "local" ]]; then
  rm -rf ${ROOT}/${OUT} && cp -r ${TMP}/${OUT} ${ROOT}/
elif [[ ${MODE} == "production" ]]; then
  push
else
  echo "invalid MODE:${MODE}!"
fi


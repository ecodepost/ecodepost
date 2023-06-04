#!/bin/bash

# 将此目录拷贝到临时目录执行
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
TMP=/tmp/proto-tmp
OUT="pb"
PROTOBIN="buf"
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

function check() {
    if [[ $1 -ne 0 ]]; then echo -e "${RED}[${REPO}] $(date +%Y-%m-%d--%H:%M:%S) $2 fail${NC}"
        echo -e "${RED}[${REPO}] ""$4\n"${NC} "$3"
        exit $1
    else
        echo -e "${GREEN}[${REPO}] $(date +%Y-%m-%d--%H:%M:%S) $2 success\n${NC}"
    fi
}

# 构建target
function build() {
    # egoctl pb generate
    local help="proto 文件未能通过[桩代码生成]！可以查看 ${PROTOBIN} generate 使用说明。静态检查错误如下所示："
    generate_out=$(${PROTOBIN} generate $1 2>&1)
    check $? "generate" "$generate_out" "$help"

    # genGoMod
    echo "[proto] generate success!"
}

# 推送至Git仓库
function push() {
    repo=https://${GIT_USER}:${GIT_PWD}@${GITHOST}/${GITPATH_PREFIX}/${OUT}.git
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
    git branch
    git diff --name-only
    git add .
    git commit -m "build: [${CI_COMMIT_REF_NAME}](https://${GITHOST}/${GITPATH_PREFIX}/proto/-/commits/${CI_COMMIT_SHORT_SHA})"
    git push -u origin ${CI_COMMIT_REF_NAME}
    echo -e "\n[proto] push success!"
}

build

#if [[ ${MODE} == "local" ]]; then
#  rm -rf ${ROOT}/${OUT} && cp -r ${TMP}/${OUT} ${ROOT}/
#elif [[ ${MODE} == "production" ]]; then
#  push
#else
#  echo "invalid MODE:${MODE}!"
#fi

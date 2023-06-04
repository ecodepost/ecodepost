#!/bin/bash

# 使用方法：./pb-check.sh [format lint breakcheck]
# 比如：./pb-check.sh format lint 或 ./pb-check.sh

RED='\033[1;31m'
GREEN='\033[1;32m'
CYAN='\033[0;36m'
NC='\033[0m'
PROTOBIN="buf"
GIT=${GIT}

REPO="protos"
# 如果未指定GIT模式，则默认不启用GIT模式
if [[ "${GIT}" == "" ]]; then
  GIT=false
fi

echo "GIT:"${GIT}
echo "PROTOBIN:${PROTOBIN}"
${PROTOBIN} version
echo -e "${NC}"

while test $# -gt 0; do
    case "$1" in
        -h|--help)
            echo "pb-check.sh - 检查 pb"
            echo " "
            echo "pb-check.sh [options] [arguments]"
            echo " "
            echo "options:"
            echo "-h, --help                    显示帮助信息"
            exit 0
            ;;
        *)
            break
            ;;
    esac
done

function relative_path() {
    old=$(pwd)
    cd $1
    x=$(pwd)
    r=${x#*$old}
    cd $old
    echo ${r:1}
}

declare -a changedSubpackages=()
function init_changed_subpackages() {
    if [[ ${GIT} == "true" ]];then
        local changes=`git diff-tree --no-commit-id --name-only -r $CI_COMMIT_SHA`
        echo CHANGES:${changes}
        for i in ${changes}
        do
            # 如果有proto文件变更则记录变更目录
            if [[ "$i" == *"/"* ]]; then
                echo ${i}
                dir=`echo ${i} | rev | cut -d '/' -f 2- | rev`
                changedDir+=("$dir")
            fi
        done
        changedSubpackages=($(echo "${changedDir[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' '))
    else
        changedSubpackages="."
    fi
}
init_changed_subpackages $@

STATUS=0
function check {
    if [[ $1 -ne 0 ]]; then echo -e "${RED}[${REPO}] $(date +%Y-%m-%d--%H:%M:%S) $2 fail${NC}"
        echo -e "${RED}[${REPO}] ""$4\n"${NC} "$3"
        STATUS=1
    else
        echo -e "${GREEN}[${REPO}] $(date +%Y-%m-%d--%H:%M:%S) $2 success\n${NC}"
    fi
}

function lint_check {
    local help="proto 文件未能通过[静态检查]！可执行 ${PROTOBIN} lint 查看并手动修复报错。静态检查错误如下所示："
    echo -e "[${REPO}] $(date +%Y-%m-%d--%H:%M:%S) lint check now..."
    echo -e "${CYAN}[${REPO}] 当前执行指令：${PROTOBIN} lint ${NC}"
    lint_out=$(${PROTOBIN} lint 2>&1)
    check $? "lint" "$lint_out" "$help"
}

function format_check {
    local help="proto 文件未能通过[格式化检查]！可执行 ${PROTOBIN} format -w 自动修复报错。格式化检查错误如下所示："
      echo -e "[${REPO}] $(date +%Y-%m-%d--%H:%M:%S) format check now..."
      echo -e "${CYAN}[${REPO}] 当前执行指令：${PROTOBIN} format -d --exit-code $1 ${NC}"
      format_out=$(${PROTOBIN} format -d --exit-code $1  2>&1)
      check $? "format" "$format_out" "$help"
}

function breakcheck_check {
    # TODO 兼容性检查的函数可能还需要优化下，必须要对比master分支么？
    local help="proto 文件未能通过[兼容性检查]！请前往可执行 ${PROTOBIN} break 查看并手动修复报错。兼容性检查错误如下所示："
    echo -e "[${REPO}] $(date +%Y-%m-%d--%H:%M:%S) break check now..."
    echo -e "${CYAN}[${REPO}] 当前执行指令：${PROTOBIN} break check --git-branch master $1 ${NC}"
    breakcheck_out=$(${PROTOBIN} break check --git-branch master $1 2>&1)
    check $? "breakcheck" "$breakcheck_out" "$help"
}

ARGS=()
if [[ $# -ne 0 ]]; then
    ARGS=${@:1}
    if [[ -z $ARGS ]]; then
      ARGS=("lint" "format" "breakcheck")
    fi
else
    ARGS=("lint" "format" "breakcheck")
fi

echo -e "[${REPO}] $(date +%Y-%m-%d--%H:%M:%S) ARGS: ${ARGS[@]}  changedSubpackages:${changedSubpackages}\n"
for i in ${ARGS[@]}; do
    case ${i} in
        "format")
            for subpackage in ${changedSubpackages}; do
                format_check $subpackage
            done;;
        "lint")
            for subpackage in ${changedSubpackages}; do
                lint_check $subpackage
            done;;
        "breakcheck")
            for subpackage in ${changedSubpackages}; do
                breakcheck_check $subpackage
            done;;
         *)
            break;;
    esac
done

if [[ ${STATUS} -ne 0 ]]; then
    echo -e "\n${RED}[${REPO}] $(date +%Y-%m-%d--%H:%M:%S) ci executed fail\n${NC}"
    exit ${STATUS}
fi

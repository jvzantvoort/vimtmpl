#!/bin/bash

C_BUILDLIST=$(
    cat <<-END
linux  amd64
linux  386
linux  arm64
linux  arm
darwin amd64
END
)

C_HELPMSG=$(
    cat <<-HELP

Development:
    fmt
    tags

Build:
    build
    buildr
    update
    install
    package
    cleanup
    check
    list

HELP
)

C_SCRIPTPATH="$(readlink -f "$0")"
C_SCRIPTDIR="$(dirname "$C_SCRIPTPATH")"
C_BUILDDIR="${C_SCRIPTDIR}/build"
C_TAGSPATH="${C_SCRIPTDIR}/tags"
C_PACKAGESPATH="${C_SCRIPTDIR}/pkg"

C_VERSION="$(git describe --tags --abbrev=0 2>/dev/null)"
C_REVISION=$(git rev-parse --short HEAD)
C_BUILDTIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
C_URL="$(git config --get remote.origin.url)"
C_REPONAME="$(basename "${C_URL}" .git)"

C_MSG_FORMAT="%-60s [ %-7s ]\n"
#shellcheck disable=SC2059
C_MSG_LEN="$(printf "${C_MSG_FORMAT}" "x" "x")"
C_MSG_LEN="${#C_MSG_LEN}"
C_MSG_LEN="$((C_MSG_LEN - 12))"

readonly C_SCRIPTPATH
readonly C_SCRIPTDIR
readonly C_BUILDDIR

LDFLAGS=""
[[ -n "${C_VERSION}" ]] && LDFLAGS="${LDFLAGS} -X main.version=${C_VERSION}"
[[ -n "${C_REVISION}" ]] && LDFLAGS="${LDFLAGS} -X main.revision=${C_REVISION}"
[[ -n "${C_BUILDTIME}" ]] && LDFLAGS="${LDFLAGS} -X main.version.BuildTime=${C_BUILDTIME}"
LDFLAGS="${LDFLAGS} -w -s"

#shellcheck disable=SC2034
strrep() {
    local num="$1"
    local re

    re='^[0-9]+$'

    if [[ $num =~ $re ]]; then
        seq 1 "${num}" | while read -r x; do
            printf "-"
        done
    fi
}

getpad() {
    local msg="$1"
    pad="${#msg}"
    pad="$((C_MSG_LEN - pad))"
    strrep "${pad}"
}

print_msg() {
    local state=$1
    shift
    local msg="$*"
    local pad

    padstr="$(getpad "${msg}")"

    msg="${msg} $(strrep "${padstr}")"

    #shellcheck disable=SC2059
    printf "${C_MSG_FORMAT}" "$msg ${padstr}" "${state}"
}

print_title() {
    local msg
    msg="$(echo "$@" | tr "[:lower:]" "[:upper:]")"
    print_msg "       " "$msg"
}

print_subtitle() {
    local msg
    msg="$(echo "$@" | sed 's/.*/\L&/; s/[a-z]*/\u&/g')"
    print_msg "-------" "$msg"
}

print_ok() { print_msg "SUCCESS" "$@"; }
print_nok() { print_msg "FAILURE" "$@"; }
print_fatal() { print_msg "FATAL" "$@"; }
print_warning() { print_msg "WARNING" "$@"; }
print_unknown() { print_msg "UNKNOWN" "$@"; }
print_skipped() { print_msg "SKIPPED" "$@"; }

test_result() {
    local retv=$1
    local message=$2
    if [[ "${retv}" == "0" ]]; then
        print_ok "${message}"
    elif [[ "${retv}" == "127" ]]; then
        print_nok "${message} (not found)"
    else
        print_nok "${message}"
    fi
}

test_fatal() {
    local retv=$1
    local message=$2
    if [[ "${retv}" == "0" ]]; then
        print_ok "${message}"
    else
        print_fatal "${message}"
        exit 1
    fi
}

die() { test_fatal 1 "FATAL: $1"; }

list_gofiles() {
    pushd "${C_SCRIPTDIR}" >/dev/null 2>&1 || die "cannot changedir to scriptdir"
    find "." -type f -name '*.go' -not -path '*/vendor/*'
    popd >/dev/null 2>&1 || die "cannot changedir back"
}

list_binaries() {
    pushd "${C_SCRIPTDIR}" >/dev/null 2>&1 || die "cannot changedir to scriptdir"
    find cmd -maxdepth 1 -mindepth 1 -type d -printf "%f\n"
    popd >/dev/null 2>&1 || die "cannot changedir back"
}

__calcdest() {
    local goos="$1"
    local arch="$2"
    local target="$3"
    printf "%s/%s/%s/%s\n" "${C_BUILDDIR}" "${goos}" "${arch}" "${target}"
}

action_fmt() {
    list_gofiles | while read -r target; do
        goimports -w "${target}"
        test_result "$?" "reformat ${target}"
    done
    echo
    echo "Changed:"
    echo
    git status -s | awk '$1 ~ /M/ && /\.go/ { printf " - %s\n", $2 }'
    echo
}

action_tags() {
    list_gofiles | xargs gotags >"${C_TAGSPATH}"
    test_result "$?" "Generate c-tags"
}

action_build() {
    local goos="$1"
    local arch="$2"
    local dest
    local exitcode

    exitcode=0

    while read -r target; do
        dest="$(__calcdest "${goos}" "${arch}" "${target}")"

        CGO_ENABLED=0 \
            GOOS=${goos} GOARCH=${arch} \
            go build -ldflags "${LDFLAGS}" \
            -o "${dest}" "./cmd/${target}"
        retv="$?"
        test_result "$retv" "    ${target}"
        [[ "${retv}" == "0" ]] || exitcode=$((exitcode + 1))
    done < <(list_binaries)
    test_fatal "${exitcode}" "build results"
}

action_buildr() {
    local goos="$1"
    local arch="$2"
    local outdir="$3"
    local destdir
    local exitcode

    exitcode=0
    destdir="${outdir}/${goos}.${arch}"
    mkdir -p "${destdir}" || die "cannot create ${destdir}"

    while read -r target; do
        CGO_ENABLED=0 \
            GOOS=${goos} GOARCH=${arch} \
            go build -ldflags "${LDFLAGS}" \
            -o "${destdir}/${target}" "./cmd/${target}"
        retv="$?"
        test_result "$retv" "    ${target}"
        [[ "${retv}" == "0" ]] || exitcode=$((exitcode + 1))
    done < <(list_binaries)
    test_fatal "${exitcode}" "build results"
}

action_dependencies() {
    local msg

    pushd "${C_SCRIPTDIR}" >/dev/null 2>&1 || die "cannot changedir to scriptdir"
    msg="go mod init"
    if [[ -e "go.mod" ]]; then
        print_skipped "${msg}"
    else
        go mod init
        test_result "$?" "${msg}"
    fi
    msg="go mod tidy"
    go mod tidy
    test_result "$?" "${msg}"

    msg="go mod vendor"
    go mod vendor
    test_result "$?" "${msg}"

    msg="go get packages"
    go get -v -t ./...
    test_result "$?" "${msg}"

    popd >/dev/null 2>&1 || die "cannot changedir back"
}

action_install() {
    local goos="$1"
    local arch="$2"
    local bindir
    bindir="$(go env GOBIN)"

    if [[ -z "${bindir}" ]]; then
        bindir="${HOME}/go/bin"
    fi

    while read -r target; do
        local dest
        dest="$(__calcdest "${goos}" "${arch}" "${target}")"
        install -m 755 "${dest}" "${bindir}/${target}"
        test_result "$?" "    ${target}"
    done < <(list_binaries)
}

action_package() {
    local goos="$1"
    local arch="$2"
    local dest
    local destdir
    local destv
    local destf

    mkdir -p "${C_PACKAGESPATH}"

    # define a version
    destv="${C_VERSION}"
    [[ -n "${destv}" ]] || destv="${C_REVISION}"
    [[ -n "${destv}" ]] || destv="$(date +%y%m%d%H%M)"

    # define a filename
    destf="${C_REPONAME}-${goos}-${arch}-${destv}.zip"

    destdir="$(dirname "$(__calcdest "${goos}" "${arch}" "none")")"

    pushd "${destdir}" >/dev/null 2>&1 || die "cannot changedir to destdir"
    list_binaries | zip -@ "${C_PACKAGESPATH}/${destf}" >/dev/null 2>&1
    test_result "$?" "  create ${destf} archive"
    popd >/dev/null 2>&1 || die "cannot changedir back"

}

action_createpackages() {
    action_dependencies
    while read -r goos arch; do
        print_subtitle "create ${goos}/${arch} archive"
        action_build "${goos}" "${arch}"
        action_package "${goos}" "${arch}"
    done < <(echo "${C_BUILDLIST}")
}

action_cleanup() {

    [[ -d "${C_BUILDDIR}" ]] && rm -rvf "${C_BUILDDIR}"
    [[ -e "${C_PACKAGESPATH}" ]] && rm -vrf "${C_PACKAGESPATH}"
    [[ -d "${C_TAGSPATH}" ]] && rm -vf "${C_TAGSPATH}"

}

action_check() {
    pushd "${C_SCRIPTDIR}" >/dev/null 2>&1 || die "cannot changedir to scriptdir"

    go vet ./...
    test_result "$?" "go vet"

    golangci-lint run ./...
    test_result "$?" "golangci-lint"

    staticcheck ./...
    test_result "$?" "staticcheck"

    popd >/dev/null 2>&1 || die "cannot changedir back"

}

action_list() {
    pushd "${C_SCRIPTDIR}" >/dev/null 2>&1 || die "cannot changedir to scriptdir"

    find "${C_SCRIPTDIR}" -mindepth 1 -maxdepth 1 -type d \
        -not -name '.git*' -not -name build -not -name vendor \
        -printf "%f\n" | sort | while read -r target; do
        find "${target}" -type f
    done
    find "${C_SCRIPTDIR}" -mindepth 1 -maxdepth 1 -type f \
        -name '*.go'
    popd >/dev/null 2>&1 || die "cannot changedir back"

}

# }}}

# interfaces {{{
do_fmt() {
    print_title "format the sources"
    action_fmt
}
do_tags() {
    print_title "generate tags"
    action_tags
}
do_cleanup() {
    print_title "cleanup"
    action_cleanup
}
do_dependencies() {
    print_title "update dependencies"
    action_dependencies
}
do_build() {
    print_title "build ${goos}/${arch}"
    action_dependencies
    action_build "$(go env GOOS)" "$(go env GOARCH)"
}
do_buildr() {
    local goos="$2"
    local arch="$3"
    local outdir="$4"

    [[ -n "${goos}" ]] || die "missing goos"
    [[ -n "${arch}" ]] || die "missing arch"
    [[ -n "${outdir}" ]] || die "missing output directory"

    print_title "build ${goos}/${arch}"
    action_buildr "${goos}" "${arch}" "${outdir}"
}

do_install() {
    print_title "install locally"
    action_install "$(go env GOOS)" "$(go env GOARCH)"
}
do_package() {
    print_title "create packages"
    action_createpackages
}
do_usage() {
    printf "USAGE:\n\n  %s <option>\n\n" "${C_SCRIPTPATH}"
    echo "${C_HELPMSG}"
    printf "\n\n"
    exit 0
}
do_check() {
    print_title "check the sources"
    action_check
}
do_list() { action_list; }
# }}}

#------------------------------------------------------------------------------#
#                                    Main                                      #
#------------------------------------------------------------------------------#

case "$1" in
  fmt)      do_fmt;;
  tags)     do_tags;;
  build)    do_build;;
  buildr)   do_buildr "$@";;
  update)   do_dependencies;;
  install)  do_install;;
  package)  do_package;;
  cleanup)  do_cleanup;;
  check)    do_check;;
  list)     do_list;;
  help)     do_usage;;
  *)        do_usage;;
esac

#------------------------------------------------------------------------------#
#                                  The End                                     #
#------------------------------------------------------------------------------#
# vim: foldmethod=marker

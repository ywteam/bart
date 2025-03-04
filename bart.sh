#!/bin/bash
# shellcheck disable=all
# Basic Application Repository Template

bart:log() {
    local status=$?
    local args=()
    case "${1}" in
    "fatal") args+=(💀 "\033[31m${1}\033[0m") ;;
    "critical") args+=(🔥 "\033[31m${1}\033[0m") ;;
    "alert") args+=(🚨 "\033[31m${1}\033[0m") ;;
    "error") args+=(❌ "\033[31m${1}\033[0m") ;;
    "warning" | "warn") args+=(🔔 "\033[33m${1}\033[0m") ;;
    "notice") args+=(📢 "\033[33m${1}\033[0m") ;;
    "info") args+=(💬 "\033[36m${1}\033[0m") ;;
    "debug") args+=(🐞 "\033[36m${1}\033[0m") ;;
    "trace") args+=(🔍 "\033[36m${1}\033[0m") ;;
    "success") args+=(✅ "\033[32m${1}\033[0m") ;;
    *) args+=(🚀 "\033[32m${1}\033[0m") ;;
    esac
    shift
    local funcName="${FUNCNAME[1]}"
    [[ "${funcName}" =~ (ywt:wksp:log|ywt:wksp:verbose) ]] && funcName="${FUNCNAME[2]}"
    echo -e "${args[@]}" "${@}" $'\t' "${funcName}(${status})" "${SECONDS}s" 1>&2
    unset args status funcName
    return 0
}
bart:log:pipe() {
    [[ -t 0 ]] && return 0
    local line
    while IFS= read -r line || [ -n "$line" ]; do
        bart:log "${@}" "$line"
    done
    unset line
    return 0
}
bart:strings:trim() {
    bart:strings:ltrim "$(bart:strings:rtrim "${1}")"
    # local str="${1}"
    # str="${str#"${str%%[![:space:]]*}"}"
    # str="${str%"${str##*[![:space:]]}"}"
    # echo -n "${str}"
    # unset str
    # return 0
}
bart:strings:rtrim() {
    local str="${1}"
    str="${str%"${str##*[![:space:]]}"}"
    echo -n "${str}"
    unset str
    return 0
}
bart:strings:ltrim(){
    local str="${1}"
    str="${str#"${str%%[![:space:]]*}"}"
    echo -n "${str}"
    unset str
    return 0
}
bart:markdown:codeblock() {
    local langName="${1}" && shift
    local code="${1}"
    {
        case "${langName}" in
        "html")
            ! command -v pandoc && bart:log "error" "Pandoc compiler not found" && return 1
            pandoc -s -o "${code}"
            ;;
        "shell" | "bash" | "sh")
            eval "${code}"
            ;;
        "python")
            command -v python3 >/dev/null 2>&1 || {
                bart:log "error" "Python compiler not found"
                return 1
            }
            # remove shebang if present
            python3 -c "$(echo "${code}" | sed '/^#!/d')"
            ;;
        "node")
            ! command -v node && bart:log "error" "Node compiler not found" && return 1
            node -e "${code}"
            ;;
        "ruby")
            ! command -v ruby && bart:log "error" "Ruby compiler not found" && return 1
            ruby -e "${code}"
            ;;
        "php")
            ! command -v php && bart:log "error" "PHP compiler not found" && return 1
            php -r "${code}"
            ;;
        "perl")
            ! command -v perl && bart:log "error" "Perl compiler not found" && return 1
            perl -e "${code}"
            ;;
        "go")
            ! command -v go && bart:log "error" "Go compiler not found" && return 1
            go run -e "${code}"
            ;;
        "java")
            ! command -v java && bart:log "error" "Java compiler not found" && return 1
            java -e "${code}"
            ;;
        "c")
            ! command -v gcc && bart:log "error" "C compiler not found" && return 1
            gcc -e "${code}"
            ;;
        "cpp")
            ! command -v g++ && bart:log "error" "C++ compiler not found" && return 1
            g++ -e "${code}"
            ;;
        "rust")
            ! command -v rustc && bart:log "error" "Rust compiler not found" && return 1
            rustc -e "${code}"
            ;;
        esac
    } >&7
    unset langName code
    return 0
}
bart:markdown:varblock(){
    local varBlock="${1}"
    # varblocks can be nested in declaration and assigned to key with dots
    # ---
    # repo:
    #   name: your_repo_name
    #   description: A sample repository
    #   license: MIT
    #   visibility: public
    #   default_branch: main
    #   owner:
    #     name: your_name
    #     email: your_email
    # ---
    bart:log "infos" "Reading variables"
    local key value lastKey    
    
    while IFS=':' read -r key value; do
        local -a log=()
        
        value=$(bart:strings:trim "${value}")
        key=$(bart:strings:rtrim "${key}")
        # expect nested keys
        if [[ ! "${key}" =~ ^[[:space:]]+ ]]; then
            key=$(bart:strings:ltrim "${key}")
            lastKey="${key}"
            local -a breadcrumbs=("${key}")
            # DOC["${key}"]="${value}"
        else
            key="${breadcrumbs[-1]}.${lastKey}.${key}"
            # DOC["${key}"]="${value}"
        fi
        # save last key
        # breadcrumbs+=("${key}")
        bart:log "trace" "Variable" "${key}" "=" "${value}"
        unset key value lastKey        
    done <<<"${varBlock}"
    unset varBlock key value lastKey breadcrumbs
    return 0
}
bart:markdown:read() {
    local -A src=(
        ["input"]="$1"
        ["lines"]=0
        ["codes"]=0
    )
    declare -A -g -I DOC=()
    ! test -f "${src["input"]}" && bart:log "error" "File not found: ${src["input"]}" && unset src && return 1
    src[lines]=0
    local codeStartAt codeEndAt codeLang codeBlock trimedLine isVarBlock varBlock
    while IFS= read -r line || [ -n "$line" ]; do
        trimedLine=$(bart:strings:trim "$line")
        ((src[lines]++))
        # bart:log "trace" "Reading line ${src["lines"]}"
        # read variables from --- block
        if [[ $trimedLine == '---' && -n "${isVarBlock}" ]]; then
            unset isVarBlock
            bart:markdown:varblock "${varBlock}"
        elif [[ $trimedLine == '---' && -z "${isVarBlock}" ]]; then
            isVarBlock=1
            varBlock=""
        elif [[ -n "${isVarBlock}" ]]; then
            # non trimed line for variables because of spaces, and variables can be nested
            varBlock+="$line"$'\n'
            # read variable line
            # bart:echo "$line" # 7>&1
            # IFS=':' read -r key value <<<"$trimedLine"
            # key=$(bart:strings:trim "$key")
            # value=$(bart:strings:trim "$value")
            # bart:log "trace" "Variable" "${key}" "=" "${value}"
            # DOC["${key}"]="${value}"
            # unset key value        
        elif [[ $trimedLine == '```'* && -z "${codeStartAt}" ]]; then
            codeStartAt=${src["lines"]}
            codeEndAt=
            codeLang=$(echo "${line//\`/}" | awk '{print $1}')
            codeBlock=""
            ((src["codes"]++))
            [[ -z "${codeLang}" ]] && unset codeStartAt codeEndAt codeLang && continue
            bart:log "trace" "Code block ${src["codes"]} in ${codeLang} starts at line ${src["lines"]}"
        elif [[ -n "${codeStartAt}" && $trimedLine == $'```' ]]; then
            codeEndAt=${src["lines"]}
            bart:log "trace" "Code block ${src["codes"]} in ${codeLang} ends at line ${codeEndAt}"
        elif [[ -n "${codeStartAt}" && -z "${codeEndAt}" ]]; then
            codeBlock+="$trimedLine"$'\n'
        elif [[ -n "${codeStartAt}" && -n "${codeEndAt}" ]]; then
            # bart:echo -n '```'${codeLang}'' $'\n' "${codeBlock}" '```' $'\n'
            bart:log "info" "Running code block ${src["codes"]} in ${codeLang}" "${#codeBlock} characters"
            bart:markdown:codeblock "${codeLang}" "${codeBlock}" # >&7
            if [[ $? -ne 0 ]]; then
                bart:log "warning" "Failed to run code block ${src["codes"]} in ${codeLang}"
            else
                bart:log "info" "Executed code block ${src["codes"]} in ${codeLang}."
            fi
            codeStartAt=""
            codeEndAt=""
            codeLang=""
            codeBlock=""
            bart:log "info" "Executed code block ${src["codes"]} in ${codeLang}."
        else
            bart:echo "$line" # 7>&1 >&7
        fi
    done <"${src["input"]}" 7>&1
    unset src line trimedLine codeStartAt codeEndAt codeLang codeBlock
    return 0
}
bart:licence:generate() {
    local licenceName="${1}"
    [[ ! "${licenceName}" =~ (mit|lgpl-3.0|mpl-2.0|agpl-3.0|unlicense|apache-2.0|gpl-3.0) ]] && bart:log "error" "Invalid licence name: ${licenceName}" && return 1
    bart:log "info" "Generating ${licenceName} licence"
    # https://choosealicense.com/
    # https://chooser-beta.creativecommons.org/
    wget -qO- "https://api.github.com/licenses/${licenceName}" >&7
    unset licenceName
    return 0
}
bart:echo() {
    echo "${@}" >&7 # 7>&1
    return 0
}
bart:github:files() {
    # create default github files
    local files=(
        "src/bart/.github/CODE_OF_CONDUCT.md"
        "src/bart/.github/CONTRIBUTING.md"
        "src/bart/.github/ISSUE_TEMPLATE.md"
        "src/bart/.github/PULL_REQUEST_TEMPLATE.md"
        "src/bart/.github/SECURITY.md"
        "src/bart/.github/CODEOWNERS"
    )
    for file in "${files[@]}"; do
        bart:log "info" "Creating ${file}"
        test -f "${file}" && rm -f "${file}"
        touch "${file}" && chmod 644 "${file}"
        local fileDir=$(dirname "${file}")
        [[ ! -d "${fileDir}" ]] && mkdir -p "${fileDir}"
        unset fileDir        
        if [[ "${file}" == *"CODE_OF_CONDUCT.md" ]]; then
            echo "# Code of Conduct" >"${file}"
        elif [[ "${file}" == *"CONTRIBUTING.md" ]]; then
            echo "# Contributing" >"${file}"
        elif [[ "${file}" == *"ISSUE_TEMPLATE.md" ]]; then
            echo "# Issue Template" >"${file}"
        elif [[ "${file}" == *"PULL_REQUEST_TEMPLATE.md" ]]; then
            echo "# Pull Request Template" >"${file}"
        elif [[ "${file}" == *"SECURITY.md" ]]; then
            echo "# Security" >"${file}"
        elif [[ "${file}" == *"CODEOWNERS" ]]; then
            echo "* @ywteam" >"${file}"
        fi
    done
    unset files file
    return 0
}
bart() {
    ! test -p /tmp/bart.sock && mkfifo /tmp/bart.sock
    exec 7<>/tmp/bart.sock
    bart:log "info" "Running bart.sh"
    bart:markdown:read ./README.mds >README.md
    bart:log "info" "Finished processing README.mds"
    return 0
}
bart
exit 0

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
bart:strings:trim(){
    local str="${1}"
    str="${str#"${str%%[![:space:]]*}"}"
    str="${str%"${str##*[![:space:]]}"}"
    echo -n "${str}"
    unset str
    return 0
}
bart:code:run(){
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
            command -v python3 >/dev/null 2>&1 || { bart:log "error" "Python compiler not found"; return 1; }
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
bart:markdown:read() {
    local -A src=(
        ["input"]="$1"
        ["lines"]=0
        ["codes"]=0
    )
    ! test -f "${src["input"]}" && bart:log "error" "File not found: ${src["input"]}" && unset src && return 1
    src[lines]=0
    local codeStartAt codeEndAt codeLang codeBlock trimedLine
    while IFS= read -r line || [ -n "$line" ]; do
        trimedLine=$(bart:strings:trim "$line")
        ((src[lines]++))
        # bart:log "trace" "Reading line ${src["lines"]}"
        if [[ $trimedLine == '```'* && -z "${codeStartAt}" ]]; then
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
            bart:log "info" "Running code block ${src["codes"]} in ${codeLang}" "${#codeBlock} characters"
            # eval "$codeBlock" 2>&1
            bart:code:run "${codeLang}" "${codeBlock}" # >&7
            [[ $? -ne 0 ]] && bart:log "warning" "Failed to run code block ${src["codes"]} in ${codeLang}"
            codeStartAt=""
            codeEndAt=""
            codeLang=""
            codeBlock=""
            bart:log "info" "Executed code block ${src["codes"]} in ${codeLang}."
        else
            bart:echo "$line" # 7>&1 >&7
        fi
    done <"${src["input"]}" 7>&1
    unset src line
    return 0
}
bart:licence:generate(){
    local licenceName="${1}"
    local licenceFile="${2:-./LICENSE}"
    local licenceText
    # https://choosealicense.com/
    # https://chooser-beta.creativecommons.org/    
}
bart:echo() {
    echo "${@}" >&7 # 7>&1
    return 0
}
bart() {
    ! test -p /tmp/bart.sock && mkfifo /tmp/bart.sock
    exec 7<>/tmp/bart.sock
    bart:log "info" "Running bart.sh"
    bart:markdown:read ./README.mds  > README.md
    bart:log "info" "Finished processing README.mds"
    return 0
}
bart
exit 0
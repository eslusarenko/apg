# bash completion for apg

_apg() {
    local cur prev
    cur=${COMP_WORDS[COMP_CWORD]}
    prev=${COMP_WORDS[COMP_CWORD-1]}

    case $prev in
        -a) COMPREPLY=($(compgen -W '0 1' -- "$cur")); return ;;
        -M) COMPREPLY=($(compgen -W 'SNCL sncl NCL ncl N' -- "$cur")); return ;;
        --strategy) COMPREPLY=($(compgen -W 'full alphanum digits' -- "$cur")); return ;;
        -r) COMPREPLY=($(compgen -f -- "$cur")); return ;;
        -E | -n | -m | -x | -c | --length) return ;;
    esac

    COMPREPLY=($(compgen -W '-a -M -E -n -m -x -r -t -l -d -q -s -c -v -h
        --length --strategy --alnum --digits --pin --no-newline -C --copy
        --version --help' -- "$cur"))
}

complete -F _apg apg

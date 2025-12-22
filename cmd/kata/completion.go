package main

import (
	"fmt"

	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/urfave/cli/v2"
)

func completionCommand(translator i18n.Translator) *cli.Command {
	return &cli.Command{
		Name:      "completion",
		Usage:     translator.T("completion.usage"),
		ArgsUsage: translator.T("completion.args_usage"),
		Action: func(ctx *cli.Context) error {
			shell := ctx.Args().First()
			if shell == "" {
				shell = "bash" // Default
			}

			switch shell {
			case "bash":
				fmt.Println(bashCompletionScript)
			case "zsh":
				fmt.Println(zshCompletionScript)
			case "fish":
				fmt.Println(fishCompletionScript)
			default:
				return fmt.Errorf("unsupported shell: %s", shell)
			}
			return nil
		},
	}
}

const bashCompletionScript = `_cli_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
    else
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi
}

complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete kata`

const zshCompletionScript = `#compdef kata

_kata_bash_autocomplete() {
  local cur opts base
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  if [[ "$cur" == "-"* ]]; then
    opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
  else
    opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
  fi
  COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
  return 0
}

autoload -U +X bashcompinit && bashcompinit
complete -o nospace -o default -o bashdefault -F _kata_bash_autocomplete kata`

const fishCompletionScript = `function _kata_completion
    set -l args (commandline -opc)
    set -l response (kata $args[2..-1] --generate-bash-completion)
    for completion in $response
        echo $completion
    end
end

complete -c kata -f -a "(_kata_completion)"`

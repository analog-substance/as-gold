# Gold

Mine data from git and breach dumps.

```azure
Usage:
  as-gold [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  consume     Consume data from source via subcommand
  dedupe      Deduplicate entries
  help        Help about any command
  merge       Merge two or more entries by UUID
  search      Search entries by various values

Flags:
      --config string   config file (default is $HOME/.gold.yaml)
      --gold string     gold data file (default "solid-gold.json")
  -h, --help            help for as-gold
  -t, --toggle          Help message for toggle

Use "as-gold [command] --help" for more information about a command.
```


## Install

```bash
go install github.com/analog-substance/as-gold@latest 
```
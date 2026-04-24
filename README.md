# repo-switcher

fisher's `z` works, but I use multiple machines which means the cache are not the same. I have all my projects under `~/Git` so might as well use it as "cache".


## Usage

Create config at `~/.config/repo-switcher/config.yaml`

```yaml
paths:
  - ~/Git
  - /your/other/git/dir
```

Shell config (fish):

```text
if type -q repo-switcher
    repo-switcher completion fish | source

    # wrap the completions so 'r' behaves like 'repo-switcher'
    complete -c r -w repo-switcher
end

function r
    argparse 'n' -- $argv
    or return

    set -l path (command repo-switcher $argv)
    test $status -eq 0; or return

    if set -q _flag_n # cd into project
        cd $path
    else # open in IDE
        idea $path &
    end
end
```

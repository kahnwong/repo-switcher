# repo-switcher

fisher's `z` works, but I use multiple machines which means the cache are not the same. I have all my projects under `~/Git` so might as well use it as "cache".


## Usage

```text
function repo-switcher
    set path (command repo-switcher $argv)
    if test $status -eq 0
        cd $path
    end
end
```

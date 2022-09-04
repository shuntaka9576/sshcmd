# sshcmd

## Installation

```bash
brew install shuntaka9576/tap/sshcmd
```

## Usage

```bash
$ sshcmd --help
Usage: sshcmd --address=STRING --port=INT-64 --user=STRING --cmd=STRING

Execute command to ssh host

Flags:
  -h, --help                Show context-sensitive help.
  -v, --version             print the version.
  -a, --address=STRING      Specify ssh connection destination
  -p, --port=INT-64         Specify ssh host port
  -u, --user=STRING         Specify ssh user
  -c, --cmd=STRING          Specify the command to run on the ssh server
  -f, --after-cmd=STRING    Specify commands to run on the host server
```

## Examples

### SCP

```bash
# short
sshcmd \
  -a "<IP>" \
  -p <PORT> \
  -u "<USER>" \
  -c "mkdir ~/work/20220101/;cp ~/logs/hoge.log ~/work/20220101/;zip -r ~/work/20220101.zip ~/work/20220101" \
  -f "scp -r -P <PORT> <USER>@<IP>:~/work/20220101.zip ."

# long
sshcmd \
  --adress "<IP>" \
  --port <PORT> \
  --user "<USER>" \
  --cmd "mkdir ~/work/20220101/;cp ~/logs/hoge.log ~/work/20220101/;zip -r ~/work/20220101.zip ~/work/20220101" \
  --after-cmd "scp -r -P <PORT> <USER>@<IP>:~/work/20220101.zip ."
```

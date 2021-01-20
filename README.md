```
$ yay -S openswan xl2tpd
$ go install github.com/nyushi/openswan-xl2tpd-helper

$ mkdir -p ~/.config/openswan-xl2tpd-helper
$ cat > ~/.config/openswan-xl2tpd-helper/config.json
{
  "server": "<ServerAddress>",
  "psk": "<PreSharedKey>",
  "user": "<User>",
  "pass": "<Pass>"
}

$ sudo openswan-xl2tpd-helper -conf ~/.config/openswan-xl2tpd-helper/config.json -interface eth0 start

$ sudo openswan-xl2tpd-helper stop
```

```
$ yay -S openswan xl2tpd

$ mkdir -p ~/.config/openswan-xl2tpd-helper
$ cat > ~/.config/openswan-xl2tpd-helper/config.json
{
  "server": "<ServerAddress>",
  "psk": "<PreSharedKey>",
  "user": "<User>",
  "pass": "<Pass>"
}

$ sudo openswan-xl2tpd-helper -server "<ServerAddress>" -psk "<PreSharedKey>" -user "<Username>" -pass "<Password>" -interface "<Interface>" start

$ sudo openswan-xl2tpd-helper stop
```

### Source
- https://github.com/chenhw2/https-dns
  
### Thanks
- https://github.com/fardog/secureoperator
- https://github.com/jackyyf/go-shadowsocks2
- https://developers.google.com/speed/public-dns/docs/dns-over-https
  
### Docker
- https://hub.docker.com/r/chenhw2/https-dns
  
### TODO
- Currently only Block DNS TYPE:```ANY```
- More thorough tests should be written
- No caching is implemented, and probably never will
  
### Usage
```
$ docker pull chenhw2/https-dns

$ docker run -d \
    -e "ARGS=-d 8.8.8.8 --edns 0.0.0.0/0" \
    -p "5300:5300/udp" \
    -p "5300:5300/tcp" \
    chenhw2/https-dns

```
### Help
```
$ docker run --rm chenhw2/https-dns -h
NAME:
   https-dns - A DNS-protocol proxy for Google's DNS-over-HTTPS service.

USAGE:
   https-dns [global options] command [command options] [arguments...]

VERSION:
   MISSING build version [git hash]

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --listen value, -l value           Serve address (default: ":5300")
   --proxy value, -p value            Proxy (SOCKS or SHADOWSOCKS) Server for HTTP GET
   --endpoint value                   Google DNS-over-HTTPS endpoint url (default: "https://dns.google.com/resolve")
   --endpoint-ips value, --eip value  IPs of the Google DNS-over-HTTPS endpoint; if provided, endpoint lookup skip
   --dns-servers value, -d value      DNS Servers used to look up the endpoint; system default is used if absent.
   --edns value, -e value             Extension mechanisms for DNS (EDNS) is parameters of the Domain Name System (DNS) protocol.
   --no-pad, -N                       Disable padding of Google DNS-over-HTTPS requests to identical length
   --insecure, -I                     Disable SSL/TLS Certificate check (for some OS without ca-certificates)
   --udp, -U                          Listen on UDP
   --tcp, -T                          Listen on TCP
   -V value                           log level for V logs (default: 2)
   --logtostderr                      log to standard error instead of files
   --help, -h                         show help
   --version, -v                      print the version

```

# vsockcat

A tool to provide port forwarding over `AF_VSOCK`. It can be used as an SSH
`ProxyCommand` to access SSH inside a VM without any network interfaces.

## Usage (inside VM)
```
listener
```

## Usage (on host)
```
ssh -oProxyCommand="vsockcat %h %p" myvm
```

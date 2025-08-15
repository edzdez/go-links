# go-links

This project attempts to emulate the functionality of Google's internal go/ links, i.e. a self-hosted URL shortcut service.

## Build

This project relies on [Bazel](https://bazel.build/).
It should automatically download the correct version of the [Go](https://go.dev) toolchain.
You may need a C compiler (and perhaps [sqlite](https://sqlite.org/)) installed to build [the sqlite driver](https://github.com/mattn/go-sqlite3).

```shell
# Clone the repository
git clone https://github.com/edzdez/go-links && cd go-links

# Build with Bazel
bazel build :server
```

This will produce a binary at `bazel-bin/server_/server`.
For persistence, we suggest copying this binary to somewhere else.
For what it's worth, I just use the project root. :P

## Run

This project is intended to run in the background as a daemon.
Pick a domain (say `go/`) and redirect it to `0.0.0.0` in `/etc/hosts` before starting the application on port `80`.
Since `80` is a privileged port, we must grant the correct permissions to the binary:

```shell
sudo setcap 'cap_net_bind_service=+ep' /path/to/binary
```

We suggest creating a [systemd](https://systemd.io/) service file to run the server automatically at startup.

```ini
[Unit]
Description=URL Shortcut Service
After=network.target

[Service]
WorkingDirectory=/path/to/this/repository
ExecStart=/path/to/binary -port=80 -db=/path/to/db
Restart=on-failure
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=default.target
```

For development purposes, the server can be run with the following command:

```shell
bazel run :server -- -port=5200 -db=/path/to/db
```

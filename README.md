[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/waitfor)](https://goreportcard.com/report/github.com/Luzifer/waitfor)
![](https://badges.fyi/github/license/Luzifer/waitfor)
![](https://badges.fyi/github/downloads/Luzifer/waitfor)
![](https://badges.fyi/github/latest-release/Luzifer/waitfor)

# Luzifer / waitfor

`waitfor` is a small utility to check the exit code of an command to be used within a shell command.

## Usage

```console
$ waitfor --help
Usage of waitfor:
  -i, --check-interval duration    How long to wait after an unsuccessful check (default 1s)
  -c, --command-timeout duration   Stop the command execution after this time
      --log-level string           Log level (debug, info, warn, error, fatal) (default "info")
  -s, --shell string               Shell to execute with the given command (must accept -c flag) (default "/bin/bash")
      --version                    Prints current version and exits
  -w, --wait-timeout duration      Stop waiting for the command after this time
```

When a `wait-timeout` is specified and the check command did not exit with status code 0 before the timeout is reached `waitfor` will exit with status 1. This ensures a shell command connected with `&&` is not executed afterwards.

## Examples

- Wait for the VPN connection to be available before executing a command using it

    ```console
    $ waitfor 'ip a | grep 10.123.0' && echo "VPN connected"
    ```

- Ensure you are online before executing a curl command

    ```console
    $ waitfor -- ping -c 1 8.8.8.8 && curl ...
    ```

- Wait at most 5m for a file to appear before accessing it

    ```console
    $ waitfor -w 5m -- ls myfile && cat myfile
    ```

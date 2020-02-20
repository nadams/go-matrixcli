# go-matrixcli
Send and receive messages from a matrix server.

## Setup 

By default the application will look in these locations for the `config.yaml` file:

|  | Linux(and BSD) | Mac | Windows |
| ---: | :---: | :---: | :---: |
| `~/.config/matrixcli` | `~/Library/Application Support/matrixcli` | `%APPDATA%\matrixcli` |
| `./config.yaml`       | `./config.yaml`                           | `.\config.yaml`       |

```yaml
accounts:
  name: default
  homeserver: https://chat.mydomain.com
  username: bot
  password: securepassword
```

## Send a message

Simple message:
```sh
$ matrixcli send '!test:chat.mydomain.com' 'my message'
```

Piped from stdin
```sh
$ mycmd | matrixcli send '!test:chat.mydomain.com'
```

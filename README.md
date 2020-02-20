# go-matrixcli
Send and receive messages from a matrix server.

## Setup 

By default the application will look in these locations for the `config.yaml` file:

| Linux(and BSD) | Mac | Windows |
| :---: | :---: | :---: |
| `~/.config/matrixcli/config.yaml` | `~/Library/Application Support/matrixcli/config.yaml` | `%APPDATA%\matrixcli\config.yaml` |
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

Message with title
```sh
$ cmd-with-long-output | matrixcli send '!test:chat.mydomain.com' --title 'Backup stuff'
```

![Rich Text](.images/rich_text.png)

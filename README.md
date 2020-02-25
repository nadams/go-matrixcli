# go-matrixcli
Send and receive messages from a matrix server.

## Setup 

Simply log into a homeserver.

```sh
$ matrixcli login https://chat.mydomain.com
```

For all commands that need an account, you can specify an account to use with the `--account` flag. This needs to be the name your gave the account. If you don't specify an account name, the first one configured will be used.

## Send a message

Simple message:
```sh
$ matrixcli send '!someid:chat.mydomain.com' 'my message'
```

Piped from stdin
```sh
$ mycmd | matrixcli send '!someid:chat.mydomain.com'
```

Message with title
```sh
$ cmd-with-long-output | matrixcli send '!someid:chat.mydomain.com' --title 'Backup Stuff'
```

![Rich Text](.images/rich_text.png)

Channel aliases are supported
```sh
$ matrixcli send '#mychannel:chat.mydomain.com' 'my msg'
```

Full channel ids are optional
```sh
$ matrixcli send '#mychannel' 'my msg'

$ matrixcli send '!someid' 'my msg'
```

This pulls the domain from the homeserver configured with the account used.

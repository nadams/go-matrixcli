# matrixcli
Send and receive messages from a matrix server.

- [x] Send messages
- [x] Multiple account/server support
- [ ] Tail messages
- [ ] Filter incoming messages
- [ ] Join/leave rooms

## Setup 

Simply log into a homeserver.

```sh
$ matrixcli accounts login https://chat.mydomain.com
```

Commands that deal with a server need an account to work with. By default the program will select the "current" account to work with. You can set the current account by looking below. If you wish to use a different account without changing the current accout, pass in the `--account <name>` flag. The name refers to the account name seen in the account list.

## List Accounts

```sh
$ matrixcli accounts list

+------------------+---------------------------+-------------------------------------+---------+
| NAME             | HOMESERVER                | USERID                              | CURRENT |
+------------------+---------------------------+-------------------------------------+---------+
| my-account       | https://chat.mydomain.com | @my-account:chat.mydomain.com       |         |
| my-other-account | https://chat.mydomain.com | @my-other-account:chat.mydomain.com | *       |
+------------------+---------------------------+-------------------------------------+---------+
```

## Set Current Account
```sh
$ matrixcli account select my-account
```

## Remove Account

```sh
$ matrixcli accounts remove my-account
```

## Send Messages

Simple message
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

Alternate account
```sh
$ matrixcli send --account my-account '#mychannel:chat.mydomain.com' 'hello there'
```

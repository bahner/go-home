# github.com/bahner/go-myspace-client

This is go-myspace-client based on [an example from go-libp2p][src].

Now you can either run with `go run`, or build and run the binary:

```shell
go run . -identity foobar -room myTopic

# or, build and run separately
go build .
 export GO_MYSPACE_CLIENT_IDENTITY=fooBar
./go-myspace-client  -room myTopic
```

## Configuration
type `./go-myspace-client -help`. Most config settings can be set with environment variables, as follows:

```bash
export GO_MYSPACE_CLIENT_LOG_LEVEL="error"
export GO_MYSPACE_CLIENT_RENDEZVOUS="myspace"
export GO_MYSPACE_CLIENT_SERVICE_NAME="myspace"
export GO_MYSPACE_CLIENT_ROOM="mytopic"
export GO_MYSPACE_CLIENT_IDENTITY="myBase58EncodedPrivkeyGeneratedByGenerate"
```

## Identity

A `-generate` parameter to generate a text version of a secret key.
The key is text formatted privKey for your node.

This key can and should be kept safely on a PostIt note on your monitor :-)
Just don't store somewhere insecure. It's your future identity.

```bash
unset HISTFILE
 export GO_MYSPACE_CLIENT_IDENTITY=FooBarABCDEFbase58
```

or specified on the command line:

```bash
./go-myspace-client -identity FooBarABCDEFbase58
```

The first is the best. (Noticed that in most shells the empty space before the command, means that the line isn't saved in history.)

## Usage

You can join a specific chat room with the `-room` flag:

```shell
go run . -room=planet-express
```

It's usually more fun to chat with others, so open a new terminal and run the app again.
If you set a custom chat room name with the `-room` flag, make sure you use the same one
for both apps. Once the new instance starts, the two chat apps should discover each other
automatically using mDNS, and typing a message into one app will send it to any others that are open.

To quit, hit `Ctrl-C`, or type `/quit` into the input field.

## Commands

- /status [sub|topic|host]
- /discover
- /enter room
- /nick Name
- /refresh

[src]: https://github.com/libp2p/go-libp2p/tree/master/examples/pubsub/chat

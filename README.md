## Transporter

Here is the newest addition to the variety of file transfer tools already available. What makes transporter different from others is that the file transfer over the network is extremely fast.

Currently, it only provides command line interface. So, Go ahead and give it a try :)

## Usage

```
Usage of transporter:
  -R	Whether to send files recursively
  -path string
    	Absolute path of the file/directory to be transferred
  -s	Specify whether application should act as sender
```

### As Sender

***simple file-transfer***
```go
transporter -s path="<absolute-path of file/>"
```

***recursive file-transfer***
```go
transporter -R -s path="<absolute-path of directory/>"
```


### As Receiver

```go
transporter
```

## License

MIT, see [LICENSE](https://github.com/pratikju/transporter/blob/master/LICENSE.md)

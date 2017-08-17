## Transporter

Here is the newest addition to the variety of file transfer tools already available. What makes transporter different from others is that the file transfer over the network is extremely fast.

Currently, it only provides command line interface. So, Go ahead and give it a try :)

## Usage

```
Usage of transporter:
  -R	Whether to send files recursively
  -mode string
    	Choose the mode of application(sender/receiver) (default "receiver")
  -p int
    	Port to connect to (default 7080)
  -path string
    	Absolute path of the file/directory to be transferred
```

### As Sender

***simple file-transfer***
```go
transporter -mode=sender path="<absolute-path of file/>"
```

***recursive file-transfer***
```go
transporter -R -mode=sender path="<absolute-path of directory/>"
```


### As Receiver

```go
transporter
```

## License

MIT, see [LICENSE](https://github.com/pratikju/transporter/blob/master/LICENSE.md)

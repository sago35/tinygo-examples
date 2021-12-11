# wioterminal/usbcdc

First, flash the program to Wio Terminal.  

```
$ tinygo flash --target wioterminal ./wioterminal/usbcdc/
```

You can then run `./wioterminal/usbcdc/cmd/wio-client` to control the Wio Terminal via usbcdc.  
`-port` must be specified to execute the command.  

```
$ go run ./wioterminal/usbcdc/cmd/wio-client -port COM8
hello world
こんにちは 世界
hello world
こんにちは 世界
hello world
```

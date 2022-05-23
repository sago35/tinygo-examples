# xiao-ble examples

TinyGo example of XIAO BLE.  
This is a Demo using Bluetooth.  

* ble-led
    * [./xiao-ble/ble-led-client/](./xiao-ble/ble-led-client/)
    * [./xiao-ble/ble-led-client-xiao/](./xiao-ble/ble-led-client-xiao/)
    * [./xiao-ble/ble-led-server/](./xiao-ble/ble-led-server/)

[![](https://img.youtube.com/vi/HWBxuMbNUTI/0.jpg)](https://www.youtube.com/watch?v=HWBxuMbNUTI)

## How to use

First write ble-led-server to the XIAO BLE.  

```shell
$ tinygo flash --target xiao-ble --size short ./xiao-ble/ble-led-server
```

Perform the following from a Bluetooth-enabled PC.  
If successful, the LED on the XIAO BLE will flash.  

```shell
$ go run ./xiao-ble/ble-led-client/
```

This source code works perfectly with TinyGo.  
The result is the same as when run from a PC.  

```shell
$ tinygo flash --target xiao-ble --size short ./xiao-ble/ble-led-client/
```

If you have another XIAO BLE, you can try the same Demo as in the video.  

```shell
$ tinygo flash --target xiao-ble --size short ./xiao-ble/ble-led-client-xiao/
```

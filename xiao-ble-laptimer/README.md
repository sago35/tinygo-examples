# xiao-ble-laptimer examples

TinyGo example of XIAO BLE.  
This is a Demo using Bluetooth.

* [./xiao-ble-laptimer/laptimer/](./ble-laptimer/laptimer/)
* [./xiao-ble-laptimer/laptimer-xiao/](./laptimer-xiao/)

![](./xiao-ble-laptimer.jpg)


## Usage

First flash it to an nRF52840 microcontroller such as xiao-ble.  

```
$ tinygo flash --target xiao-ble --size short ./xiao-ble-laptimer/laptimer-xiao
   code    data     bss |   flash     ram
  11228     156    8012 |   11384    8168
```

Then, start a program to measure the lap time on the computer.  
If the BLE scan fails for 3 minutes, the next lap is considered to have progressed.  


```
$ go run ./xiao-ble-laptimer/laptimer
```

To change to a value other than 3 minutes, change the following settings.  


```go
	thresh := 3 * time.Minute
```

Below is the actual log file.  

```
2022/10/22 10:46:25 2562047:47:16 0 found
2022/10/22 10:53:06 lost
2022/10/22 11:26:39 00:40:14 1 found
2022/10/22 11:30:02 lost
2022/10/22 11:38:02 00:11:23 2 found
2022/10/22 11:41:19 lost
2022/10/22 11:48:16 00:10:13 3 found
2022/10/22 11:51:31 lost
2022/10/22 11:58:43 00:10:26 4 found
2022/10/22 12:01:58 lost
```

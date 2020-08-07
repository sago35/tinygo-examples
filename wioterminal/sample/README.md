# TinyGo examples

## wioterminal sample

TinyGo example of wioterminal.  
https://github.com/sago35/tinygo-examples/tree/master/wioterminal/sample  

[![](https://img.youtube.com/vi/9IpI9rUMXOs/0.jpg)](https://www.youtube.com/watch?v=9IpI9rUMXOs)

## Build

```
tinygo flash -target wioterminal -size short github.com/sago35/tinygo-examples/wioterminal/sample
```

or

```
tinygo build -o app.uf2 -target wioterminal -size short github.com/sago35/tinygo-examples/wioterminal/sample
```

or if you don't want to build it, use the uf2 file.

* [wioterminal_tinygo_sample.uf2](./wioterminal_tinygo_sample.uf2)


## Environment

```
$ tinygo version
tinygo version 0.14.0-dev linux/amd64 (using go version go1.14.4 and LLVM version 10.0.1)
```

In this example, I used the following.  

* Wio Terminal
  * LCD Screen : ILI9341 320 x 240
  * Accelerometer : LIS3DHTR
  * LED
    * machine.LED
  * Infrared Emitter
    * machine.WIO\_IR
  * Operation interface
    * machine.WIO\_KEY\_A
    * machine.WIO\_KEY\_B
    * machine.WIO\_KEY\_C
    * machine.WIO\_5S\_UP
    * machine.WIO\_5S\_LEFT
    * machine.WIO\_5S\_RIGHT
    * machine.WIO\_5S\_DOWN
    * machine.WIO\_5S\_PRESS
  * Light Sensor
    * machine.WIO\_LIGHT
  * Speaker
    * machine.WIO\_BUZZER
  * Microphone
    * machine.WIO\_MIC

## Link

* English
  * https://tinygo.org/microcontrollers/wioterminal/
  * https://wiki.seeedstudio.com/Wio-Terminal-Getting-Started/
* Japanese
  * Wio Terminal で TinyGo プログラミングを始めよう
    * https://qiita.com/sago35/items/92b22e8cbbf99d0cd3ef


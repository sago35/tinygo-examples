# TinyGo examples

## pinintterrupt

Example of combining pin interrupts and goroutine in TinyGo.  
https://github.com/sago35/tinygo-examples/tree/master/pininterrupt  

[![](https://img.youtube.com/vi/A-EA5iqDp7k/0.jpg)](https://www.youtube.com/watch?v=A-EA5iqDp7k)

info:  
pin interrupt is still in the PR stage.  

* https://github.com/tinygo-org/tinygo/pull/1094
* https://github.com/tinygo-org/tinygo/pull/1111

## Build

```
tinygo flash -target pyportal -size short github.com/sago35/tinygo-examples/pininterrupt
```

or

```
tinygo build -o app.uf2 -target pyportal -size short github.com/sago35/tinygo-examples/pininterrupt
```

## Environment

```
PyPortal
D3  : Switch (with pull-up)
D4  : LED
I2C : AE-AQM0802A (I2C 8 x 2 character LCD with st7032)
```

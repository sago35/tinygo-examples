# TinyGo examples

## wioterminal goroutines

TinyGo example of wioterminal.  
https://github.com/sago35/tinygo-examples/tree/master/wioterminal/goroutines  

[![](https://img.youtube.com/vi/-dJ-o2cH_Fk/0.jpg)](https://www.youtube.com/watch?v=-dJ-o2cH_Fk)

## Summary

This example shows the following, received via channel.  
TinyGo uses goroutine / channel to simplify the asynchronous process.  

- cnt1, which increments every 77ms
- cnt2, which increments every 500ms

## Build

```
tinygo flash -target wioterminal -size short github.com/sago35/tinygo-examples/wioterminal/goroutines
```

or

```
tinygo build -o app.uf2 -target wioterminal -size short github.com/sago35/tinygo-examples/wioterminal/goroutines
```

or if you don't want to build it, use the uf2 file.

* [wioterminal_goroutines.uf2](./wioterminal_goroutines.uf2)


## Environment

```
$ tinygo version
tinygo version 0.15.0 linux/amd64 (using go version go0.1.0 and LLVM version 10.0.1)
```

## Link

* English
  * https://tinygo.org/microcontrollers/wioterminal/
  * https://wiki.seeedstudio.com/Wio-Terminal-Getting-Started/
* Japanese
  * Wio Terminal で TinyGo プログラミングを始めよう
    * https://qiita.com/sago35/items/92b22e8cbbf99d0cd3ef



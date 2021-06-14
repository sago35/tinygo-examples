smoketest:
	tinygo build -o test.hex -size short -target pyportal    ./pininterrupt
	tinygo build -o test.hex -size short -target wioterminal ./wioterminal/buttons
	tinygo build -o test.hex -size short -target wioterminal ./wioterminal/buzzer
	tinygo build -o test.hex -size short -target wioterminal ./wioterminal/goroutines
	tinygo build -o test.hex -size short -target wioterminal ./wioterminal/gpio
	tinygo build -o test.hex -size short -target wioterminal ./wioterminal/ir
	tinygo build -o test.hex -size short -target wioterminal ./wioterminal/light_sensor
	tinygo build -o test.hex -size short -target wioterminal ./wioterminal/lis3dh
	tinygo build -o test.hex -size short -target wioterminal ./wioterminal/microphone
	tinygo build -o test.hex -size short -target wioterminal ./wioterminal/qspi_flash
	tinygo build -o test.hex -size short -target wioterminal ./wioterminal/sample

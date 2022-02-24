smoketest:
	tinygo build -o /tmp/test.hex -size short -target pyportal    ./pininterrupt
	tinygo build -o /tmp/test.hex -size short -target wioterminal ./wioterminal/buttons
	tinygo build -o /tmp/test.hex -size short -target wioterminal ./wioterminal/buzzer
	tinygo build -o /tmp/test.hex -size short -target wioterminal ./wioterminal/goroutines
	tinygo build -o /tmp/test.hex -size short -target wioterminal ./wioterminal/ir
	tinygo build -o /tmp/test.hex -size short -target wioterminal ./wioterminal/light_sensor
	tinygo build -o /tmp/test.hex -size short -target wioterminal ./wioterminal/lis3dh
	tinygo build -o /tmp/test.hex -size short -target wioterminal ./wioterminal/microphone
	tinygo build -o /tmp/test.hex -size short -target wioterminal ./wioterminal/qspi_flash
	tinygo build -o /tmp/test.hex -size short -target wioterminal ./wioterminal/sample
	tinygo build -o /tmp/test.hex -size short -target wioterminal ./wioterminal/usbcdc
	tinygo build -o /tmp/test.hex -size short -target wioterminal ./deviceid
	go build -o /tmp/test.exe ./wioterminal/usbcdc/cmd/wio-client

fmt-check:
	@unformatted=$$(gofmt -l `find . -name "*.go"`); [ -z "$$unformatted" ] && exit 0; echo "Unformatted:"; for fn in $$unformatted; do echo "  $$fn"; done; exit 1

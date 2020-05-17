package st7032

import (
	"fmt"
	"machine"
	"time"
)

type Device struct {
	bus  *machine.I2C
	buf  []byte
	addr uint8

	displayFunction uint8
	displayControl  uint8
	displayMode     uint8

	cols     uint8
	lines    uint8
	currline uint8
	numlines uint8
}

func New(i2c *machine.I2C, address uint8) *Device {
	return &Device{
		bus:  i2c,
		buf:  make([]byte, 2),
		addr: address,
	}
}

func (d *Device) Configure() {
	d.cols = 8
	d.lines = 2
	d.currline = 0
	d.numlines = d.lines

	d.displayFunction = LCD_8BITMODE | LCD_1LINE | LCD_5x8DOTS

	if d.lines > 1 {
		d.displayFunction |= LCD_2LINE
	}

	// power on and external reset
	// wait > 40ms
	time.Sleep(40 * time.Millisecond)

	d.normalFunctionSet()

	d.extendFunctionSet()
	d.command(LCD_EX_SETBIASOSC | LCD_BIAS_1_5 | LCD_OSC_183HZ)             // 1/5bias, OSC=183Hz@3.0V
	d.command(LCD_EX_CONTRASTSETL)                                          // Contrast set
	d.command(LCD_EX_POWICONCONTRASTH | LCD_ICON_OFF | LCD_BOOST_ON | 0x02) // Power/ICON control/Contrast set
	d.command(LCD_EX_FOLLOWERCONTROL | LCD_FOLLOWER_ON | LCD_RAB_2_00)      // internal follower circuit is turn on
	time.Sleep(300 * time.Millisecond)                                      // Wait time >200ms (for power stable)
	d.normalFunctionSet()

	d.displayControl = 0x00
	d.setDisplayControl(LCD_DISPLAYON | LCD_CURSOROFF | LCD_BLINKOFF)

	d.Clear()

	d.displayMode = 0x00
	//d.setEntryMode(LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT)
}

func (d *Device) Clear() {
	d.command(LCD_CLEARDISPLAY)
	time.Sleep((3 + 2) * time.Millisecond)
}

func (d *Device) setDisplayControl(setBit uint8) {
	d.displayControl |= setBit
	d.command(LCD_DISPLAYCONTROL | d.displayControl)
}

func resetDisplayControl(resetBit uint8) {
	//_displaycontrol &= ~resetBit;
	//command(LCD_DISPLAYCONTROL | _displaycontrol)
}

func (d *Device) setEntryMode(setBit uint8) {
	d.displayMode |= setBit
	d.command(LCD_ENTRYMODESET | d.displayMode)
}

func (d *Device) resetEntryMode(resetBit uint8) {
	//_displaymode &= ~resetBit;
	//command(LCD_ENTRYMODESET | _displaymode)
}

func (d *Device) normalFunctionSet() {
	d.command(LCD_FUNCTIONSET | d.displayFunction)
}

func (d *Device) extendFunctionSet() {
	d.command(LCD_FUNCTIONSET | d.displayFunction | LCD_EX_INSTRUCTION)
}

func (d *Device) SetCursor(col, row uint8) {
	var row_offsets = []uint8{0x00, 0x40, 0x14, 0x54}

	if row > d.numlines {
		row = d.numlines - 1 // we count rows starting w/0
	}

	d.command(LCD_SETDDRAMADDR | (col + row_offsets[row]))
}

func (d *Device) Print(str string) {
	for _, s := range []byte(str) {
		d.write(uint8(s))
	}
	//d.bus.WriteRegister(d.addr, 0x40, []byte(str))
	//time.Sleep(26300 * 5 * time.Nanosecond) // >26.3us
}

func (d *Device) SetContrast(contrast uint8) {
	if contrast > 63 {
		contrast = 63
	}
	if contrast < 1 {
		contrast = 1
	}
	d.extendFunctionSet()
	d.command(LCD_EX_CONTRASTSETL | (contrast & 0x0F))                                         // Contrast set
	d.command(LCD_EX_POWICONCONTRASTH | LCD_ICON_ON | LCD_BOOST_ON | ((contrast >> 4) & 0x03)) // Power/ICON control/Contrast set
	d.normalFunctionSet()
}

/*********** mid level commands, for sending data/cmds */
func (d *Device) writeByte(reg uint8, data byte) {
	d.buf[0] = reg
	d.buf[1] = data
	err := d.bus.Tx(uint16(d.addr), d.buf, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (d *Device) command(value uint8) {
	//fmt.Printf("command: %08X\r\n", value)
	d.writeByte(0x00, value)
	//time.Sleep(26300 * time.Nanosecond) // >26.3us
	time.Sleep(26300 * 5 * time.Nanosecond) // >26.3us
	//time.Sleep(4 * time.Millisecond) // >26.3us
}

func (d *Device) write(value uint8) uint8 {
	d.writeByte(0x40, value)
	time.Sleep(26300 * time.Nanosecond) // >26.3us
	//time.Sleep(4 * time.Millisecond) // >26.3us

	return 1
}

//class ST7032 : public Print {
//public:
//	ST7032(int i2c_addr = ST7032_I2C_DEFAULT_ADDR);
//
//    void begin(uint8_t cols, uint8_t rows, uint8_t charsize = LCD_5x8DOTS);
//
//    void setContrast(uint8_t cont);
//    void setIcon(uint8_t addr, uint8_t bit);
//    void clear();
//    void home();
//
//    void noDisplay();
//    void display();
//    void noBlink();
//    void blink();
//    void noCursor();
//    void cursor();
//    void scrollDisplayLeft();
//    void scrollDisplayRight();
//    void leftToRight();
//    void rightToLeft();
//    void autoscroll();
//    void noAutoscroll();
//
//    void createChar(uint8_t location, uint8_t charmap[]);
//    void setCursor(uint8_t col, uint8_t row);
//    virtual size_t write(uint8_t value);
//    void command(uint8_t value);
//
//private:
//    void setDisplayControl(uint8_t setBit);
//    void resetDisplayControl(uint8_t resetBit);
//    void setEntryMode(uint8_t setBit);
//    void resetEntryMode(uint8_t resetBit);
//    void normalFunctionSet();
//    void extendFunctionSet();
//
////  void send(uint8_t, uint8_t);
///*
//    uint8_t _rs_pin; // LOW: command.   HIGH: character.
//    uint8_t _rw_pin; // LOW: write to LCD.  HIGH: read from LCD.
//    uint8_t _enable_pin; // activated by a HIGH pulse.
//    uint8_t _data_pins[8];
//*/
//    uint8_t _displayfunction;
//    uint8_t _displaycontrol;
//    uint8_t _displaymode;
////  uint8_t _iconfunction;
//
//    uint8_t _initialized;
//
//    uint8_t _numlines;
//    uint8_t _currline;
//
//	uint8_t _i2c_addr;
//};

// https://github.com/FaBoPlatform/FaBoLCDmini-AQM0802A-Library/blob/master/src/FaBoLCDmini_AQM0802A.cpp

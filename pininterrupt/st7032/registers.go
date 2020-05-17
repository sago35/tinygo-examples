package st7032

// Registers
const (
	// commands
	LCD_CLEARDISPLAY   = 0x01
	LCD_RETURNHOME     = 0x02
	LCD_ENTRYMODESET   = 0x04
	LCD_DISPLAYCONTROL = 0x08
	LCD_CURSORSHIFT    = 0x10
	LCD_FUNCTIONSET    = 0x20
	LCD_SETCGRAMADDR   = 0x40
	LCD_SETDDRAMADDR   = 0x80

	LCD_EX_SETBIASOSC       = 0x10 // Bias selection / Internal OSC frequency adjust
	LCD_EX_SETICONRAMADDR   = 0x40 // Set ICON RAM address
	LCD_EX_POWICONCONTRASTH = 0x50 // Power / ICON control / Contrast set(high byte)
	LCD_EX_FOLLOWERCONTROL  = 0x60 // Follower control
	LCD_EX_CONTRASTSETL     = 0x70 // Contrast set(low byte)

	// flags for display entry mode
	LCD_ENTRYRIGHT          = 0x00
	LCD_ENTRYLEFT           = 0x02
	LCD_ENTRYSHIFTINCREMENT = 0x01
	LCD_ENTRYSHIFTDECREMENT = 0x00

	// flags for display on/off control
	LCD_DISPLAYON  = 0x04
	LCD_DISPLAYOFF = 0x00
	LCD_CURSORON   = 0x02
	LCD_CURSOROFF  = 0x00
	LCD_BLINKON    = 0x01
	LCD_BLINKOFF   = 0x00

	// flags for display/cursor shift
	LCD_DISPLAYMOVE = 0x08
	LCD_CURSORMOVE  = 0x00
	LCD_MOVERIGHT   = 0x04
	LCD_MOVELEFT    = 0x00

	// flags for function set
	LCD_8BITMODE       = 0x10
	LCD_4BITMODE       = 0x00
	LCD_2LINE          = 0x08
	LCD_1LINE          = 0x00
	LCD_5x10DOTS       = 0x04
	LCD_5x8DOTS        = 0x00
	LCD_EX_INSTRUCTION = 0x01 // IS: instruction table select

	// flags for Bias selection
	LCD_BIAS_1_4 = 0x08 // bias will be 1/4
	LCD_BIAS_1_5 = 0x00 // bias will be 1/5

	// flags Power / ICON control / Contrast set(high byte)
	LCD_ICON_ON   = 0x08 // ICON display on
	LCD_ICON_OFF  = 0x00 // ICON display off
	LCD_BOOST_ON  = 0x04 // booster circuit is turn on
	LCD_BOOST_OFF = 0x00 // booster circuit is turn off
	LCD_OSC_122HZ = 0x00 // 122Hz@3.0V
	LCD_OSC_131HZ = 0x01 // 131Hz@3.0V
	LCD_OSC_144HZ = 0x02 // 144Hz@3.0V
	LCD_OSC_161HZ = 0x03 // 161Hz@3.0V
	LCD_OSC_183HZ = 0x04 // 183Hz@3.0V
	LCD_OSC_221HZ = 0x05 // 221Hz@3.0V
	LCD_OSC_274HZ = 0x06 // 274Hz@3.0V
	LCD_OSC_347HZ = 0x07 // 347Hz@3.0V

	// flags Follower control
	LCD_FOLLOWER_ON  = 0x08 // internal follower circuit is turn on
	LCD_FOLLOWER_OFF = 0x00 // internal follower circuit is turn off
	LCD_RAB_1_00     = 0x00 // 1+(Rb/Ra)=1.00
	LCD_RAB_1_25     = 0x01 // 1+(Rb/Ra)=1.25
	LCD_RAB_1_50     = 0x02 // 1+(Rb/Ra)=1.50
	LCD_RAB_1_80     = 0x03 // 1+(Rb/Ra)=1.80
	LCD_RAB_2_00     = 0x04 // 1+(Rb/Ra)=2.00
	LCD_RAB_2_50     = 0x05 // 1+(Rb/Ra)=2.50
	LCD_RAB_3_00     = 0x06 // 1+(Rb/Ra)=3.00
	LCD_RAB_3_75     = 0x07 // 1+(Rb/Ra)=3.75
)

//https://github.com/tomozh/arduino_ST7032/blob/master/ST7032.h

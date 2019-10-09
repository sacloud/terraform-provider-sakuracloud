package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/mitchellh/go-vnc"
	"github.com/sacloud/libsacloud/sacloud"
)

const keyLeftShift uint32 = 0xFFE1

// SendCommandOption is the Option value of VNC VNCSendCommand
type SendCommandOption struct {
	UseUSKeyboard  bool
	Debug          bool
	ProgressWriter io.Writer
}

// NewSendCommandOption returns new SendCommandOption
func NewSendCommandOption() *SendCommandOption {
	return &SendCommandOption{
		UseUSKeyboard:  false,
		Debug:          false,
		ProgressWriter: ioutil.Discard,
	}
}

type vncClientConnWrapper struct {
	*vnc.ClientConn
	w io.Writer
}

func newVNCClientConnWrapper(c *vnc.ClientConn, debug bool, w io.Writer) *vncClientConnWrapper {
	if w == nil || !debug {
		w = ioutil.Discard
	}

	return &vncClientConnWrapper{
		ClientConn: c,
		w:          w,
	}

}

func (c *vncClientConnWrapper) log(format string, a ...interface{}) {
	fmt.Fprintf(c.w, format, a...)
}

// VNCSendCommand sends command over VNC connection
func VNCSendCommand(vncProxyInfo *sacloud.VNCProxyResponse, command string, option *SendCommandOption) error {
	host := vncProxyInfo.ActualHost()

	fmt.Fprintf(option.ProgressWriter, "Connecting to VM via VNC...\n")
	// Connect to VNC
	nc, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, vncProxyInfo.Port))
	if err != nil {
		return fmt.Errorf("Error connecting to VNC: %s", err)
	}
	defer nc.Close()

	fmt.Fprintf(option.ProgressWriter, "Handshaking with VNC...\n")
	// Connect to VNC
	auth := []vnc.ClientAuth{&vnc.PasswordAuth{Password: vncProxyInfo.Password}}
	c, err := vnc.Client(nc, &vnc.ClientConfig{Auth: auth, Exclusive: false})
	if err != nil {
		return fmt.Errorf("Error handshaking with VNC: %s", err)
	}
	defer c.Close()

	wrapper := newVNCClientConnWrapper(c, option.Debug, option.ProgressWriter)

	fmt.Fprintf(option.ProgressWriter, "Sending command...\n")
	vncSendString(wrapper, command, option.UseUSKeyboard)
	fmt.Fprintf(option.ProgressWriter, "Done\n")
	return nil
}

func vncSendString(c *vncClientConnWrapper, original string, useUSKeyboard bool) {
	// Scancodes reference: https://github.com/qemu/qemu/blob/master/ui/vnc_keysym.h
	special := make(map[string]uint32)
	special["<bs>"] = 0xFF08
	special["<del>"] = 0xFFFF
	special["<enter>"] = 0xFF0D
	special["<esc>"] = 0xFF1B
	special["<f1>"] = 0xFFBE
	special["<f2>"] = 0xFFBF
	special["<f3>"] = 0xFFC0
	special["<f4>"] = 0xFFC1
	special["<f5>"] = 0xFFC2
	special["<f6>"] = 0xFFC3
	special["<f7>"] = 0xFFC4
	special["<f8>"] = 0xFFC5
	special["<f9>"] = 0xFFC6
	special["<f10>"] = 0xFFC7
	special["<f11>"] = 0xFFC8
	special["<f12>"] = 0xFFC9
	special["<return>"] = 0xFF0D
	special["<tab>"] = 0xFF09
	special["<up>"] = 0xFF52
	special["<down>"] = 0xFF54
	special["<left>"] = 0xFF51
	special["<right>"] = 0xFF53
	special["<spacebar>"] = 0x020
	special["<insert>"] = 0xFF63
	special["<home>"] = 0xFF50
	special["<end>"] = 0xFF57
	special["<pageUp>"] = 0xFF55
	special["<pageDown>"] = 0xFF56
	special["<leftAlt>"] = 0xFFE9
	special["<leftCtrl>"] = 0xFFE3
	special["<leftShift>"] = 0xFFE1
	special["<rightAlt>"] = 0xFFEA
	special["<rightCtrl>"] = 0xFFE4
	special["<rightShift>"] = 0xFFE2
	special["<leftWin>"] = 0xFF5B
	special["<rightWin>"] = 0xFF5C

	shiftedChars := "!\"#$%&'()=~|{`+*}<>?"
	if useUSKeyboard {
		shiftedChars = "~!@#$%^&*()_+{}|:\"<>?"
	}

	// TODO(mitchellh): Ripe for optimizations of some point, perhaps.
	for len(original) > 0 {
		var keyCode uint32
		keyShift := false

		if strings.HasPrefix(original, "<leftAltOn>") {
			keyCode = special["<leftAlt>"]
			original = original[len("<leftAltOn>"):]
			c.log("Special code '<leftAltOn>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftCtrlOn>") {
			keyCode = special["<leftCtrl>"]
			original = original[len("<leftCtrlOn>"):]
			c.log("Special code '<leftCtrlOn>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftShiftOn>") {
			keyCode = special["<leftShift>"]
			original = original[len("<leftShiftOn>"):]
			c.log("Special code '<leftShiftOn>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftAltOff>") {
			keyCode = special["<leftAlt>"]
			original = original[len("<leftAltOff>"):]
			c.log("Special code '<leftAltOff>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftCtrlOff>") {
			keyCode = special["<leftCtrl>"]
			original = original[len("<leftCtrlOff>"):]
			c.log("Special code '<leftCtrlOff>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftShiftOff>") {
			keyCode = special["<leftShift>"]
			original = original[len("<leftShiftOff>"):]
			c.log("Special code '<leftShiftOff>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightAltOn>") {
			keyCode = special["<rightAlt>"]
			original = original[len("<rightAltOn>"):]
			c.log("Special code '<rightAltOn>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightCtrlOn>") {
			keyCode = special["<rightCtrl>"]
			original = original[len("<rightCtrlOn>"):]
			c.log("Special code '<rightCtrlOn>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightShiftOn>") {
			keyCode = special["<rightShift>"]
			original = original[len("<rightShiftOn>"):]
			c.log("Special code '<rightShiftOn>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightAltOff>") {
			keyCode = special["<rightAlt>"]
			original = original[len("<rightAltOff>"):]
			c.log("Special code '<rightAltOff>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightCtrlOff>") {
			keyCode = special["<rightCtrl>"]
			original = original[len("<rightCtrlOff>"):]
			c.log("Special code '<rightCtrlOff>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightShiftOff>") {
			keyCode = special["<rightShift>"]
			original = original[len("<rightShiftOff>"):]
			c.log("Special code '<rightShiftOff>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftWinOn>") {
			keyCode = special["<leftWin>"]
			original = original[len("<leftWinOn>"):]
			c.log("Special code '<leftWinOn>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftWinOff>") {
			keyCode = special["<leftWin>"]
			original = original[len("<leftWinOff>"):]
			c.log("Special code '<leftWinOff>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}
		if strings.HasPrefix(original, "<rightWinOn>") {
			keyCode = special["<rightWin>"]
			original = original[len("<rightWinOn>"):]
			c.log("Special code '<rightWinOn>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightWinOff>") {
			keyCode = special["<rightWin>"]
			original = original[len("<rightWinOff>"):]
			c.log("Special code '<rightWinOff>' found, replacing with: %d\n", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<wait>") {
			c.log("Special code '<wait>' found, sleeping one second\n")
			time.Sleep(1 * time.Second)
			original = original[len("<wait>"):]
			continue
		}

		if strings.HasPrefix(original, "<wait5>") {
			c.log("Special code '<wait5>' found, sleeping 5 seconds\n")
			time.Sleep(5 * time.Second)
			original = original[len("<wait5>"):]
			continue
		}

		if strings.HasPrefix(original, "<wait10>") {
			c.log("Special code '<wait10>' found, sleeping 10 seconds\n")
			time.Sleep(10 * time.Second)
			original = original[len("<wait10>"):]
			continue
		}

		if strings.HasPrefix(original, "<wait") && strings.HasSuffix(original, ">") {
			re := regexp.MustCompile(`<wait([0-9hms]+)>$`)
			dstr := re.FindStringSubmatch(original)
			if len(dstr) > 1 {
				c.log("Special code %s found, sleeping\n", dstr[0])
				if dt, err := time.ParseDuration(dstr[1]); err == nil {
					time.Sleep(dt)
					original = original[len(dstr[0]):]
					continue
				}
			}
		}

		for specialCode, specialValue := range special {
			if strings.HasPrefix(original, specialCode) {
				c.log("Special code '%s' found, replacing with: %d\n", specialCode, specialValue)
				keyCode = specialValue
				original = original[len(specialCode):]
				break
			}
		}

		if keyCode == 0 {
			r, size := utf8.DecodeRuneInString(original)
			original = original[size:]
			keyCode = uint32(r)
			keyShift = unicode.IsUpper(r) || strings.ContainsRune(shiftedChars, r)

			c.log("Sending char '%c', code %d, shift %v\n", r, keyCode, keyShift)
		}

		if keyShift {
			c.KeyEvent(keyLeftShift, true)
		}

		c.KeyEvent(keyCode, true)
		time.Sleep(time.Second / 10)
		c.KeyEvent(keyCode, false)
		time.Sleep(time.Second / 10)

		if keyShift {
			c.KeyEvent(keyLeftShift, false)
		}

		// qemu is picky, so no matter what, wait a small period
		time.Sleep(100 * time.Millisecond)
	}
}

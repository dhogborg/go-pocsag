package pocsag

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dhogborg/go-pocsag/internal/datatypes"
	"github.com/dhogborg/go-pocsag/internal/utils"

	"github.com/fatih/color"
)

const (
	POCSAG_PREAMBLE     uint32 = 0x7CD215D8
	POCSAG_IDLE         uint32 = 0x7A89C197
	POCSAG_BATCH_LEN    int    = 512
	POCSAG_CODEWORD_LEN int    = 32
)

type CodewordType string

const (
	CodewordTypeAddress CodewordType = "ADDRESS"
	CodewordTypeMessage CodewordType = "MESSAGE"
	CodewordTypeIdle    CodewordType = "IDLE"
)

type MessageType string

const (
	MessageTypeAuto            MessageType = "auto"
	MessageTypeAlphanumeric    MessageType = "alpha"
	MessageTypeBitcodedDecimal MessageType = "bcd"
)

// ParsePOCSAG takes bits decoded from the stream and parses them for
// batches of codewords then prints them using the specefied message type.
func ParsePOCSAG(bits []datatypes.Bit, messagetype MessageType) []*Message {

	pocsag := &POCSAG{}

	batches, err := pocsag.ParseBatches(bits)
	if err != nil {
		println(err.Error())
		return []*Message{}
	}

	if DEBUG && LEVEL > 1 {
		for i, batch := range batches {
			println("")
			println("Batch: ", i)
			batch.Print()
		}
	}

	return pocsag.ParseMessages(batches)
}

type POCSAG struct{}

// ParseBatches takes bits decoded from the stream and parses them for
// batches of codewords.
func (p *POCSAG) ParseBatches(bits []datatypes.Bit) ([]*Batch, error) {

	batches := []*Batch{}

	var start = -1
	var batchno = -1
	// synchornize with the decoded bits
	for a := 0; a < len(bits)-32; a += 1 {

		bytes := utils.MSBBitsToBytes(bits[a:a+32], 8)

		if isPreamble(bytes) {

			batchno += 1
			start = a + 32

			// for file output as bin data
			batchbits := bits[a : a+POCSAG_BATCH_LEN+32]
			stream := utils.MSBBitsToBytes(batchbits, 8)

			if DEBUG && LEVEL > 2 {
				out, err := os.Create(fmt.Sprintf("batches/batch-%d.bin", batchno))
				if err != nil {
					return nil, err
				}
				out.Write(stream)
			}

			batch, err := NewBatch(bits[start : start+POCSAG_BATCH_LEN])
			if err != nil {
				println(err.Error())
			} else {
				batches = append(batches, batch)
			}

		}

	}

	if start < 0 {
		return nil, fmt.Errorf("could not obtain message sync")
	}

	return batches, nil

}

// ParseMessages takes a bundle of codeword from a series of batches and
// compiles them into messages.
// A message starts with an address codeword and a bunch of message codewords follows
// until either the batch ends or an idle codeword appears.
func (p *POCSAG) ParseMessages(batches []*Batch) []*Message {

	messages := []*Message{}

	var message *Message
	for _, b := range batches {

		for _, codeword := range b.Codewords {

			switch codeword.Type {
			// append current and begin new message
			case CodewordTypeAddress:
				if message != nil {
					messages = append(messages, message)
				}
				message = NewMessage(codeword)

			// append current but dont start new
			case CodewordTypeIdle:
				if message != nil {
					messages = append(messages, message)
				}
				message = nil

			case CodewordTypeMessage:
				if message != nil {
					message.AddPayload(codeword)
				} else {
					red.Println("Message codeword without sync!")
				}

			default:
				red.Println("Unknown codeword encounterd")
			}
		}
	}

	if message != nil {
		messages = append(messages, message)
	}

	return messages
}

// Message construct holds refernces to codewords.
// The Payload is a seies of codewords of message type.
type Message struct {
	Timestamp  time.Time
	Reciptient *Codeword
	Payload    []*Codeword
}

// NewMessage creates a new message construct ready to accept payload codewords
func NewMessage(reciptient *Codeword) *Message {
	return &Message{
		Timestamp:  time.Now(),
		Reciptient: reciptient,
		Payload:    []*Codeword{},
	}
}

func (m *Message) Print(messagetype MessageType) {
	green.Println("-- Message --------------")
	green.Println("Reciptient: ", m.ReciptientString())

	if !m.IsValid() {
		red.Println("This message has parity check errors. Contents might be corrupted")
	}

	println("")
	print(m.PayloadString(messagetype))
	println("")
	println("")

}

func (m *Message) Write(path string, messagetype MessageType) {
	if !os.IsPathSeparator(path[len(path)-1]) {
		path += "/"
	}

	now := time.Now()
	timestr := now.Format("20060102_15.04.05")
	file, err := os.Create(path + m.ReciptientString() + "_" + timestr + ".txt")
	defer file.Close()

	if err != nil {
		println("error creating file: " + err.Error())
		return
	}

	file.WriteString("Time: " + now.String() + "\n")
	file.WriteString("Reciptient: " + m.ReciptientString() + "\n")
	file.WriteString("-------------------\n")
	file.WriteString(m.PayloadString(messagetype) + "\n")

}

// AddPayload codeword to a message. Must be codeword of CodewordTypeMessage type
// to make sense.
func (m *Message) AddPayload(codeword *Codeword) {
	m.Payload = append(m.Payload, codeword)
}

// ReciptientString returns the reciptient address as a hexadecimal representation,
// with the function bits as 0 or 1.
func (m *Message) ReciptientString() string {
	bytes := utils.MSBBitsToBytes(m.Reciptient.Payload[0:17], 8)
	addr := uint(bytes[1])
	addr += uint(bytes[0]) << 8

	return fmt.Sprintf("%X%d%d", addr,
		m.Reciptient.Payload[18].Int(),
		m.Reciptient.Payload[19].Int())
}

// IsValid returns true if no parity bit check errors occurs in the message payload
// or the reciptient address.
func (m *Message) IsValid() bool {

	if !m.Reciptient.ValidParity {
		return false
	}

	for _, c := range m.Payload {
		if !c.ValidParity {
			return false
		}
	}
	return true
}

// PayloadString can try to decide to print the message as bitcoded decimal ("bcd") or
// as an alphanumeric string. There is not always a clear indication which is correct,
// so we can force either type by setting messagetype to something other than Auto.
func (m *Message) PayloadString(messagetype MessageType) string {

	bits := m.concactenateBits()

	alphanum := m.AlphaPayloadString(bits)
	bcd := utils.BitcodedDecimals(bits)

	var decided = MessageTypeAuto
	if messagetype == MessageTypeAuto {
		decided = m.estimateMessageType(alphanum, bcd)
	} else {
		decided = messagetype
	}

	switch decided {
	case MessageTypeAlphanumeric:
		return alphanum
	case MessageTypeBitcodedDecimal:
		return bcd
	default:
		return alphanum
	}

}

// AlphaPayloadString takes bits in LSB to MSB order and decodes them as
// 7 bit bytes that will become ASCII text.
// Characters outside of ASCII can occur, so we substitude the most common.
func (m *Message) AlphaPayloadString(bits []datatypes.Bit) string {

	str := string(utils.LSBBitsToBytes(bits, 7))

	// translate to utf8
	charmap := map[string]string{
		"[":  "Ä",
		"\\": "Ö",
		"]":  "Å",
		"{":  "ä",
		"|":  "ö",
		"}":  "å",
		"~":  "ß",
	}
	for b, s := range charmap {
		str = strings.Replace(str, b, s, -1)
	}

	return str
}

// concactenateBits to a single bitstream
func (m *Message) concactenateBits() []datatypes.Bit {
	msgsbits := []datatypes.Bit{}

	for _, cw := range m.Payload {
		if cw.Type == CodewordTypeMessage {
			msgsbits = append(msgsbits, cw.Payload...)
		}
	}
	return msgsbits
}

// estimateMessageType tries to figure out if a message is in alpha-numeric format
// or Bitcoded decimal format. There is not always a clear indication which is correct,
// so we try to guess based on some assumptions:
// 1) Alpha numeric messages contains mostly printable charaters.
// 2) BCD messages are usually shorter.
func (m *Message) estimateMessageType(persumed_alpha, persumed_bcd string) MessageType {

	// MessageTypeAuto...
	// Start guessing

	odds_a := 0
	odds_b := 0

	specials_a := 0
	specials_b := 0
	const alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 *.,-()<>\n\r"
	const bcdnum = "0123456789"

	for _, char := range persumed_alpha {
		r := rune(char)
		if strings.IndexRune(alphanum, r) < 0 {
			specials_a++
		}
	}

	for _, char := range persumed_bcd {
		r := rune(char)
		if strings.IndexRune(bcdnum, r) < 0 {
			specials_b++
		}
	}

	partspecial_a := float32(specials_a) / float32(len(persumed_alpha))
	partspecial_b := float32(specials_b) / float32(len(persumed_bcd))

	if partspecial_a < 0.2 {
		odds_a += 2
	}

	if partspecial_a >= 0.2 {
		odds_a += 1
	}

	if partspecial_a >= 0.50 {
		odds_b += 2
	}

	// special charaters are uncommon in bcd messages
	if partspecial_b == 0 {
		odds_b += 2
	}

	if partspecial_b >= 0.1 {
		odds_a += 1
	}

	if len(persumed_alpha) > 25 {
		odds_a += 2
	}

	if len(persumed_alpha) > 40 {
		odds_a += 3
	}

	if DEBUG {
		red.Printf("odds: %d/%d\nspecial: %d/%d (%0.0f%%)\n\n", odds_a, odds_b, specials_a, specials_b, (partspecial_a * 100))
	}

	if odds_a > odds_b {
		return MessageTypeAlphanumeric
	} else {
		return MessageTypeBitcodedDecimal
	}

}

//-----------------------------
// Batch
// Contains codewords. We dont care about frames, we keep the 16 codewords in a single list.
type Batch struct {
	Codewords []*Codeword
}

func NewBatch(bits []datatypes.Bit) (*Batch, error) {
	if len(bits) != POCSAG_BATCH_LEN {
		return nil, fmt.Errorf("invalid number of bits in batch: ", len(bits))
	}

	words := []*Codeword{}
	for a := 0; a < len(bits); a = a + POCSAG_CODEWORD_LEN {
		word, err := NewCodeword(bits[a : a+POCSAG_CODEWORD_LEN])
		if err != nil {
			println(err.Error())
		} else {
			words = append(words, word)
		}
	}

	b := &Batch{
		Codewords: words,
	}

	return b, nil
}

// Print will print a list with the codewords of this bach. FOr debugging.
func (b *Batch) Print() {
	for _, w := range b.Codewords {
		w.Print()
	}
}

//-----------------------------
// Codeword contains the actual data. There are two codewords per frame,
// and there are 8 frames per batch.
// Type can be either Address or Message, and a special Idle codeword will occur
// from time to time.
// Payload is a stream of bits, and ValidParity bit is set on creation for later
// reference.
type Codeword struct {
	Type        CodewordType
	Payload     []datatypes.Bit
	ParityBits  []datatypes.Bit
	EvenParity  datatypes.Bit
	ValidParity bool
}

// NewCodeword takes 32 bits, creates a new codeword construct, sets the type and checks for parity errors.
func NewCodeword(bits []datatypes.Bit) (*Codeword, error) {
	if len(bits) != 32 {
		return nil, fmt.Errorf("invalid number of bits for codeword: ", len(bits))
	}

	mtype := CodewordTypeAddress
	if bits[0] == true {
		mtype = CodewordTypeMessage
	}

	bytes := utils.MSBBitsToBytes(bits, 8)
	if isIdle(bytes) {
		mtype = CodewordTypeIdle
	}

	c := &Codeword{
		Type:        mtype,
		Payload:     bits[1:21],
		ParityBits:  bits[21:31],
		EvenParity:  bits[31],
		ValidParity: utils.ParityCheck(bits[0:31], bits[31]),
	}

	return c, nil
}

// Print the codeword contents and type to terminal. For debugging.
func (c *Codeword) Print() {

	payload := ""
	var color *color.Color = blue

	switch c.Type {

	case CodewordTypeAddress:
		payload = c.Adress()
		color = red

	case CodewordTypeMessage:
		payload = ""
		color = green

	default:
		color = blue

	}

	parity := utils.TernaryStr(c.ValidParity, "", "*")

	color.Printf("%s %s %s\n", c.Type, payload, parity)
}

func (c *Codeword) Adress() string {
	bytes := utils.MSBBitsToBytes(c.Payload[0:17], 8)
	addr := uint(bytes[1])
	addr += uint(bytes[0]) << 8

	return fmt.Sprintf("%X:%s%s", addr,
		utils.TernaryStr(bool(c.Payload[18]), "1", "0"),
		utils.TernaryStr(bool(c.Payload[19]), "1", "0"))

}

// Utilities
// isPreamble matches 4 bytes to the POCSAG preamble 0x7CD215D8
func isPreamble(bytes []byte) bool {

	var a uint32 = 0
	a += uint32(bytes[0]) << 24
	a += uint32(bytes[1]) << 16
	a += uint32(bytes[2]) << 8
	a += uint32(bytes[3])

	return a == POCSAG_PREAMBLE
}

// isIdle matches 4 bytes to the POCSAG idle codeword 0x7A89C197
func isIdle(bytes []byte) bool {

	var a uint32 = 0
	a += uint32(bytes[0]) << 24
	a += uint32(bytes[1]) << 16
	a += uint32(bytes[2]) << 8
	a += uint32(bytes[3])

	return a == POCSAG_IDLE

}

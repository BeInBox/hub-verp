package verp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"time"
)

var Inc = 0
var publicIp string
var Token string

type MiniVerp struct {
	User  int
	Inc   int
	Time  time.Time
	Ip    string
	Chain string
	MX    string
}

func NewMinVerp(userId int) *MiniVerp {
	var m = MiniVerp{User: userId}
	m.Ip = publicIp
	m.Inc = GetInc()
	m.Time = time.Now()
	return &m

}
func (o *MiniVerp) Encode() (chainOut string) {
	// Encode Vars
	//fmt.Println("Pad  Timestamp")
	t := PadLeft(fmt.Sprint(o.Time.Unix()), "0", 10)
	//fmt.Println("Pad  Ip")
	longIp := PadLeft(fmt.Sprint(ip2Long(o.Ip)), "0", 10)
	//fmt.Println("Pad  Customer")
	customer := PadLeft(fmt.Sprint(o.User), "0", 8)
	//fmt.Println("Pad  Inc")
	inc := PadLeft(fmt.Sprint(o.Inc), "0", 4)
	// Create Struct
	chainOut = t + inc + customer + longIp
	//fmt.Println("Chain", chainOut)
	n := new(big.Int)
	n, _ = n.SetString(chainOut, 10)
	//fmt.Println("Int", n)
	hexOut := fmt.Sprintf("%x", n)
	//fmt.Println("Hex", hexOut)
	hexOut = PadRight(hexOut, "X", 30)
	//fmt.Println("HexPadded", hexOut)
	o.Chain = hexOut
	return hexOut
}
func DecodeMiniVerp(chain string) *MiniVerp {
	// Encode Vars
	var mini MiniVerp

	tmpMini := chain[0 : len(chain)-4]
	//tmpMini := strings.Replace(chain, "X", "", -1)
	n := new(big.Int)
	n.SetString(tmpMini, 16)
	nS := fmt.Sprint(n)

	timed, _ := strconv.ParseInt(nS[0:10], 10, 64)
	mini.Time = time.Unix(timed, 0)
	mini.Inc, _ = strconv.Atoi(nS[10:14])
	mini.User, _ = strconv.Atoi(nS[14:22])
	ipLong, _ := strconv.Atoi(nS[22:])
	mini.Ip = long2Ip(int64(ipLong))
	return &mini
}
func PadLeft(chain string, padValue string, length int) string {
	oLen := len(chain)
	var chainOut string
	if oLen < length {
		for i := 1; i < (length - oLen + 1); i++ {
			chainOut += padValue
		}
		chainOut += chain
	} else {
		chainOut = chain
	}
	return chainOut
}
func PadRight(chain string, padValue string, length int) string {
	oLen := len(chain)
	if oLen < length {
		for i := 1; i < (length - oLen + 1); i++ {
			chain += padValue
		}

	}
	return chain
}
func GetInc() int {
	Inc += 1
	if Inc > 9999 {
		Inc = 0
	}
	return Inc
}
func long2Ip(ipInt int64) string {

	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

func GetPublicIp() (ipV4 string) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.beinbox.io/ping/simple", nil)
	if err != nil {
		return ""
	}
	req.Header = http.Header{"Authorization": []string{Token}}
	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	return string(b)
}
func ip2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

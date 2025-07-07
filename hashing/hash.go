package hashing

import (
	"math/big"
	"math/rand"
	// "time"

	"github.com/google/uuid"
)

// Get hash by method chaining

const charSet62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// func init() {
// 	rand.Seed(time.Now().UnixNano())
// }

// random number from 0 to hi - 1
func getRandomNumber(hi int) int {
	return rand.Intn(hi)
}

type hash struct {
	uuidH uuid.UUID
	number *big.Int
	encode string
}

func (h *hash) getUUID() *hash {
	h.uuidH = uuid.New()
	return h
}

func (h *hash) convertToBigInt() *hash {
	h.number = new(big.Int).SetBytes(h.uuidH[:])
	return h
}

func (h *hash) convertToBase62() *hash {
	encode := ""
	zero := big.NewInt(0)
	base := big.NewInt(62)
	remainder := new(big.Int)
	for h.number.Cmp(zero) == 1 {
		h.number, remainder = new(big.Int).DivMod(h.number, base, remainder)
		c := string(charSet62[remainder.Int64()])
		encode = c + encode
	}
	h.encode = encode
	return h
}

func (h *hash) getShortCode() string {
	shortCode := ""
	length := getRandomNumber(5) + 3
	index := getRandomNumber(len(h.encode) - length)
	for i := index; i < index + length; i++ {
		// index := getRandomNumber(length)
		shortCode = shortCode + string(h.encode[i])
	}

	// return h.encode[index : index + length + 1]

	// h.uuidH = uuid.Nil
	// h.number = nil
	return shortCode
}

func GetUUIDHash() string {
	// uuidH := getUUID()
	// number := convertToBigInt(uuidH)
	// number := convertToBase62(convertToBigInt(getUUID()))
	// shortCode := getShortCode(number)

	// h := hash{}
	
	// // Error because of getUUID() receives pointer and hash{} is value
	// return Hash{}.getUUID().convertToBigInt().convertToBase62().getShortCode()

	// Now hash{} a pointer
	return (&hash{}).getUUID().convertToBigInt().convertToBase62().getShortCode()
}

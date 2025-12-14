package shortener

import "crypto/md5"

func GenerateShortCode(longURL string) string {
	hash := md5.Sum([]byte(longURL))
	base62Code := encodeBase62(hash)
	return base62Code[:6] // 6 chars = 62^6 = 56 billion combos
}

package passphrase

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewParams() (p params) {
	var memory uint32 = 64 * 1024
	var iterations uint32 = 3
	var parallelism uint8 = 2
	var saltLength uint32 = 16
	var keyLength uint32 = 32
	p = params{memory, iterations, parallelism, saltLength, keyLength}
	return p
}

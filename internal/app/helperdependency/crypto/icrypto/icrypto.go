package icrypto

type Crypter interface {
	Encrypter
	Decrypter
}

type Encrypter interface {
	Hash(plain string) (string, error)
}

type Decrypter interface {
	Compare(hashed string, plain string) error
}

/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    rsa
 * @Date:    2021/10/25 4:20 下午
 * @package: encrypt
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"strings"
)

/*
 * 生成RSA公钥和私钥并保存在对应的目录文件下
 * 参数bits: 指定生成的秘钥的长度, 单位: bit
 */

func GenRsaKey(path string, bitNum int) {
	// 1. 生成私钥文件
	// GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bitNum)
	if err != nil {
		log.Panicf("Init rsa key err=%v", err)
	}
	// 2. MarshalPKCS1PrivateKey将rsa私钥序列化为ASN.1 PKCS#1 DER编码
	derPrivateStream := x509.MarshalPKCS1PrivateKey(privateKey)

	// 3. Block代表PEM编码的结构, 对其进行设置
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPrivateStream,
	}

	// 4. 创建文件
	var rsaPath, pubPath string
	if strings.HasSuffix(path, "/") {
		rsaPath = path + "id_rsa"
		pubPath = path + "id_rsa.pub"
	} else {
		rsaPath = path + "/id_rsa"
		pubPath = path + "/id_rsa.pub"
	}
	privateFile, err := os.Create(rsaPath)
	if err != nil {
		log.Panicf("Init rsa key err=%v", err)
	}
	defer privateFile.Close()

	// 5. 使用pem编码, 并将数据写入文件中
	err = pem.Encode(privateFile, &block)
	if err != nil {
		log.Panicf("Init rsa key err=%v", err)
	}

	// 1. 生成公钥文件
	publicKey := privateKey.PublicKey
	derPublicStream, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		log.Panicf("Init rsa key err=%v", err)
	}

	block = pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPublicStream,
	}

	publicFile, err := os.Create(pubPath)
	defer publicFile.Close()
	if err != nil {
		log.Panicf("Init rsa key err=%v", err)
	}

	// 2. 编码公钥, 写入文件
	err = pem.Encode(publicFile, &block)
	if err != nil {
		log.Panicf("Init rsa key err=%v", err)
	}
}

func readRsaFile(path string) ([]byte, error) {
	// 根据文件名读出内容
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}

	// 从数据中解析出pem块
	block, _ := pem.Decode(buf)
	if block == nil {
		return nil, errors.New("pem.Decode block is nil")
	}
	return block.Bytes, nil
}

func ReadPrivateKey(path string) (*rsa.PrivateKey, error) {
	// 根据文件名读出内容
	bytes, err := readRsaFile(path)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PrivateKey(bytes)
}

func ReadPublicKey(path string) (*rsa.PublicKey, error) {
	// 根据文件名读出内容
	bytes, err := readRsaFile(path)
	if err != nil {
		return nil, err
	}
	pub, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		return nil, err
	}
	if publicKey, ok := pub.(*rsa.PublicKey); ok {
		return publicKey, nil
	}
	return nil, errors.New("PublicKey is not rsa")
}

// RSAEncrypt RSA公钥加密
func RSAEncrypt(plaintext []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	srcLen := len(plaintext)
	stage := publicKey.Size() - 11
	if srcLen <= stage {
		return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
	}
	var resp []byte
	for i := 0; i < srcLen; i += stage {
		end := i + stage
		if end >= srcLen {
			end = srcLen - 1
		}
		bts := plaintext[i:end]
		its, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, bts)
		if err != nil {
			return nil, err
		}
		resp = append(resp, its...)
	}
	return resp, nil
}

// RSADecrypt  RSA私钥解密
func RSADecrypt(cipher []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	cipherLen := len(cipher)
	stage := privateKey.Size()
	if cipherLen <= stage {
		return rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipher)
	}

	var resp []byte
	for i := 0; i < cipherLen; i += stage {
		end := i + stage
		if end >= cipherLen {
			end = cipherLen - 1
		}
		bts := cipher[i:end]
		its, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, bts)
		if err != nil {
			return nil, err
		}
		resp = append(resp, its...)
	}
	return resp, nil
}

func defaultPrivateKey() (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(id_rsa))
	if block == nil {
		return nil, errors.New("pem.Decode block is nil")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func defaultPublicKey() (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(id_rsa_pub))
	if block == nil {
		return nil, errors.New("pem.Decode block is nil")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if publicKey, ok := pub.(*rsa.PublicKey); ok {
		return publicKey, nil
	}
	return nil, errors.New("PublicKey is not rsa")
}

func DefaultEncrypt(plaintext []byte) ([]byte, error) {
	publicKey, err := defaultPublicKey()
	if err != nil {
		return nil, err
	}
	return RSAEncrypt(plaintext, publicKey)
}

func DefaultDecrypt(cipher []byte) ([]byte, error) {
	privateKey, err := defaultPrivateKey()
	if err != nil {
		return nil, err
	}
	return RSADecrypt(cipher, privateKey)
}

const (
	id_rsa = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA3sLnb5DBZ6Yt7X3cACcER+e7U1jLx4+Pl3UlSp9Op2WznLx4
rxYaBGJ/lFbzWK7FMvcvbf6yJsB+4oKqeY1c9/Wqsch4YPzjWy6FCJvX2pkYdp5y
4d7N1YrCl5CeqIhGAsSRIqKixqcXm5vdATOcT1wuCuqMVpLD8EfK6vst/jwgBYMC
Z+4DYYU4eilX0IJhZF48e9z8DVlKzzIDBMfuIVzwDaWSIHBTtCQkefSFuharUEnS
3Kdv5squ7/o2cNQyth7lT4kLmCa4IinG2ywc57vDX2DxHgFw6ZIpmwTrAzGAoxTe
JaUtZvM5uX2w6ubklXPSZI/RV+p/pyptTjJsmQIDAQABAoIBAF3trQy76vE7fw4v
Db76tLFlIvXH8VUaUZ+5g5nthorLNWsXhYO+PAYxSj5QU5fHSdttoxAsXw48CMSV
+C/8zYC4k9sW/rtWpr9h5DJ3FBNWjpwlv6dB/WTXd6nVDzFdFLhCDjiefyhoeGni
1NOW4YgNBFUSBU7T387HVuayNviFtEy6kxjFLR84K9QApSAw2Svtk2k/OhQgUhIU
dSmiwuJm/32E/zXnAcA24Sxo0YFUSP2BdWRn75afp7BkW6CxiiuMyM3OsPIy1xgF
vS41XYvet6ckMuxxXaL6vqP4FyYerWAA7LdZtwZU7YLJgk0zkDyerw7sJ6lZx4PA
4QflUikCgYEA7FP76BWpKxp1y4eooW2hbe+rFDr3yyAk5PhcZHofRqspIEPqP7ZY
JtJB9fg3kX3THxHjxEIEJk9jbh8xsTsumsByACPAcIUalO6uN+/hbMMaD+AlIIee
H7cVVcJYQfqHRlSLcmPuNTkjJanl+rxzxqZR2waK+lYG1efy4pop2bcCgYEA8U3S
LfLPhvR4LDDOx1Ag88g8NtiKTKs5Ht+mGi1tV23PcnCNa+tE2WvfnKloCp2hQk0D
VQLmizZJfpZPDR7Ywniaa2sw1jUMT5FDXtu0aRh+UTnw0DFCjyH+hr7EA3wwIjIX
/If6UJy4pFvcz4ViQwG55eYrUSEdLybzdomoLC8CgYBeqHz7xsFNnHmWHi4zCoBg
UmEJ51ydJbDhbYFEVY31dlUwsUNAZb8FPa5h4RwQ0H8hsL60O2UCG0ZOM3xp6rSO
E4sV9zv7VbUB9mAd051NPRBRT4xPbUKunVyUTxWex8QrLW31UYV2F/619UlClv0g
kjmzKvm7r31pbFBi5zDgHQKBgQCkYILg0hsUr8x9LfJuS+Nmrex24CO1/p22rluU
UYW+nQtHxexQq8AG7DkzmyNIuAB2DchLTOKHyr9eAD5xjaXTNBzdN/PYt+JfAXGH
WNsZDJpf8rxc1nGk119votwcE6kmYkF8wZR+/YO6UumkZBR/2NkuBxFS/t/Gkx4e
jb+F+QKBgQCKv9PSREXaXzhnw6FRte63+FynN/H/kF5qBH/cZ0cANBdusIevbJdi
sqdztSDkAC5YNCn+I8uoD+cUEnWstuOtxg5NGSQFPtZBoMoVkDFwDTewPpI0UHeV
aBjONBN2CbVTe5FIaZfS+8e08JxtTblObKW3r2UaDtbjuZqP/bQ/Zw==
-----END RSA PRIVATE KEY-----`

	id_rsa_pub = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3sLnb5DBZ6Yt7X3cACcE
R+e7U1jLx4+Pl3UlSp9Op2WznLx4rxYaBGJ/lFbzWK7FMvcvbf6yJsB+4oKqeY1c
9/Wqsch4YPzjWy6FCJvX2pkYdp5y4d7N1YrCl5CeqIhGAsSRIqKixqcXm5vdATOc
T1wuCuqMVpLD8EfK6vst/jwgBYMCZ+4DYYU4eilX0IJhZF48e9z8DVlKzzIDBMfu
IVzwDaWSIHBTtCQkefSFuharUEnS3Kdv5squ7/o2cNQyth7lT4kLmCa4IinG2ywc
57vDX2DxHgFw6ZIpmwTrAzGAoxTeJaUtZvM5uX2w6ubklXPSZI/RV+p/pyptTjJs
mQIDAQAB
-----END RSA PUBLIC KEY-----
`
)

package ethereum

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"

	"github.com/tendermint/tendermint/crypto/armor"
	// dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/crypto/bcrypt"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/xsalsa20symmetric"

	"github.com/wangfeiping/wallet/wallet/adapter/types"
)

const (
	blockTypePrivKey        = "TENDERMINT PRIVATE KEY"
	infoSuffix              = "info"
	BcryptSecurityParameter = 12
)

// localInfo is the public information about a locally stored key
type LocalInfo struct {
	Name         string `json:"name"`
	PubKey       string `json:"pubkey"`
	PrivKeyArmor string `json:"privkey"`
	Address      string `json:"address"`
}

var _ types.Adapter = (*AdapterEthereum)(nil)

type AdapterEthereum struct {
}

// CreateAccount returns the account info that created with name, password and mnemonic input.
func (e *AdapterEthereum) CreateAccount(rootDir, name,
	password, mnemonic string) (ko types.KeyOutput, err error) {
	if name == "" {
		err = types.ErrMissingName
		return
	}
	if password == "" {
		err = types.ErrMissingPassword
		return
	}
	//generate wallet with mnemonic
	if mnemonic == "" {
		err = types.ErrMissingSeed
		return
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		err = fmt.Errorf("mnemonic is invalid")
		return
	}
	//convert mnemonic string to seed byte
	seed := bip39.NewSeed(mnemonic, "")

	//dpath for the key base path derive:  m / purpose' / coin_type' / account' / change / address_index (m/44'/60'/0'/0/0)
	var dpath accounts.DerivationPath
	dpath, err = accounts.ParseDerivationPath(`m/44'/60'/0'/0/0`)
	if err != nil {
		return
	}

	//fetch the masterKey for the wallet
	var masterKey *hdkeychain.ExtendedKey
	masterKey, err = hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	//masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return
	}

	key := masterKey

	for _, n := range dpath {
		key, err = key.Child(n)
		if err != nil {
			return
		}
	}

	//generate the privateKey and pubKey
	var privateKey *btcec.PrivateKey
	privateKey, err = key.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	if err != nil {
		return
	}
	//then the pubKey
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = fmt.Errorf("the public key structure conversion failed")
		return
	}
	//the address, pubKey, PrivKey with hexString format
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	pubKeyHex := hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA))[4:]

	//Armor and encrypt the privateKey
	privateKeyByte := crypto.FromECDSA(privateKeyECDSA)
	var saltBytes, encBytes []byte
	saltBytes, encBytes, err = encryptPrivKey(privateKeyByte, password)
	header := map[string]string{
		"kdf":  "bcrypt",
		"salt": fmt.Sprintf("%X", saltBytes),
	}
	priKeyAmor := armor.EncodeArmor(blockTypePrivKey, header, encBytes)
	//priKeyHex := hexutil.Encode(crypto.FromECDSA(privateKeyECDSA))[2:]

	//gather the local info into struct
	LInfo := &LocalInfo{
		Name:         name,
		PubKey:       pubKeyHex,
		PrivKeyArmor: priKeyAmor,
		Address:      address,
	}

	//write the local info by key
	// key1 := []byte(fmt.Sprintf("%s.%s", name, infoSuffix))
	// var serializeInfo []byte
	// serializeInfo, err = json.Marshal(LInfo)
	// if err != nil {
	// 	return
	// }

	//init a go level DB to store the key and Info, specify the ethkeys
	// db, err := dbm.NewGoLevelDB("keys", filepath.Join(rootDir, "ethkeys"))
	// if err != nil {
	// 	return err.Error()
	// }

	// db.SetSync(key1, serializeInfo)
	// // store a pointer to the infokey by address for fast lookup
	// addrKey := []byte(fmt.Sprintf("%s.%s", address, "addr"))
	// db.SetSync(addrKey, key1)
	// //Close the db to release the lock
	// db.Close()
	//fetch the result
	ko.Name = LInfo.Name
	ko.Type = "local"
	ko.Address = LInfo.Address
	ko.PubKey = LInfo.PubKey
	// ko.PrivKeyArmor = hexutil.Encode(crypto.FromECDSA(privateKeyECDSA))[2:]
	ko.Denom = "ETH"
	ko.Seed = mnemonic
	return
}

// encrypt the given privKey with the passphrase using a randomly
// generated salt and the xsalsa20 cipher. returns the salt and the
// encrypted priv key.
func encryptPrivKey(privKey []byte, passphrase string) (
	saltBytes []byte, encBytes []byte, err error) {
	saltBytes = tcrypto.CRandBytes(16)
	var key []byte
	key, err = bcrypt.GenerateFromPassword(saltBytes, []byte(passphrase), BcryptSecurityParameter)
	if err != nil {
		return
	}
	key = tcrypto.Sha256(key) // get 32 bytes
	privKeyBytes := privKey
	encBytes = xsalsa20symmetric.EncryptSymmetric(privKeyBytes, key)
	return
}

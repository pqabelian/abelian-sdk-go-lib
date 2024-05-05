package main

/*
#include <memory.h>
*/
import "C"

import (
	core "abelian.info/sdk/core"
	pb "abelian.info/sdk/proto"

	"errors"
	"unsafe"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func packToOutData(data []byte, outData []byte) error {
	size := len(data)
	if size > len(outData) {
		return errors.New("outData is too small")
	}

	C.memcpy(unsafe.Pointer(&outData[0]), unsafe.Pointer(&size), 4)
	if size > 0 {
		C.memcpy(unsafe.Pointer(&outData[4]), unsafe.Pointer(&data[0]), C.size_t(size))
	}

	return nil
}

func packToRetData(data []byte) *C.char {
	size := int32(len(data))
	retData := C.malloc(C.size_t(size) + 4)

	C.memcpy(retData, unsafe.Pointer(&size), 4)
	if size > 0 {
		C.memcpy(unsafe.Pointer(uintptr(retData)+4), unsafe.Pointer(&data[0]), C.size_t(size))
	}

	return (*C.char)(retData)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func ignore(args interface{}) {
}

func unmarshalArgs(argsData []byte, args protoreflect.ProtoMessage) {
	err := proto.Unmarshal(argsData, args)
	panicIf(err)
}

func marshalResultAndPackToRetData(result protoreflect.ProtoMessage) *C.char {
	resultData, err := proto.Marshal(result)
	panicIf(err)
	return (*C.char)(packToRetData(resultData))
}

//export GenerateSafeCryptoSeed
func GenerateSafeCryptoSeed(argsData []byte) *C.char {
	// Unmarshal args.
	args := &pb.GenerateSafeCryptoSeedArgs{}
	unmarshalArgs(argsData, args)

	// Prepare data.
	ignore(args)

	// Do real work.
	cryptoSeed, err := core.GenerateSafeCryptoSeed()
	panicIf(err)

	// Marshal result and return it.
	result := &pb.GenerateSafeCryptoSeedResult{
		CryptoSeed: cryptoSeed.Slice(),
	}
	return marshalResultAndPackToRetData(result)
}

//export GenerateCryptoKeysAndAddress
func GenerateCryptoKeysAndAddress(argsData []byte) *C.char {
	// Unmarshal args.
	args := &pb.GenerateCryptoKeysAndAddressArgs{}
	unmarshalArgs(argsData, args)

	// Prepare data.
	cryptoSeed := core.AsBytes(args.GetCryptoSeed())

	// Do real work.
	ckaa, err := core.GenerateCryptoKeysAndAddress(cryptoSeed)
	panicIf(err)

	// Marshal result and return it.
	result := &pb.GenerateCryptoKeysAndAddressResult{
		SpendSecretKey:    ckaa.SpendSecretKey.Slice(),
		SerialNoSecretKey: ckaa.SerialNoSecretKey.Slice(),
		ViewSecretKey:     ckaa.ViewSecretKey.Slice(),
		CryptoAddress:     ckaa.CryptoAddress.Data(),
	}
	return marshalResultAndPackToRetData(result)
}

//export GetAbelAddressFromCryptoAddress
func GetAbelAddressFromCryptoAddress(argsData []byte) *C.char {
	// Unmarshal args.
	args := &pb.GetAbelAddressFromCryptoAddressArgs{}
	unmarshalArgs(argsData, args)

	// Prepare data.
	cryptoAddress := core.NewCryptoAddress(args.GetCryptoAddress())
	chainID := args.GetChainID()

	// Do real work.
	abelAddress := core.NewAbelAddressFromCryptoAddress(cryptoAddress, int8(chainID))

	// Marshal result and return it.
	result := &pb.GetAbelAddressFromCryptoAddressResult{
		AbelAddress: abelAddress.Data(),
	}
	return marshalResultAndPackToRetData(result)
}

//export GetCryptoAddressFromAbelAddress
func GetCryptoAddressFromAbelAddress(argsData []byte) *C.char {
	// Unmarshal args.
	args := &pb.GetCryptoAddressFromAbelAddressArgs{}
	unmarshalArgs(argsData, args)

	// Prepare data.
	abelAddress := core.NewAbelAddress(args.GetAbelAddress())

	// Do real work.
	cryptoAddress := abelAddress.GetCryptoAddress()

	// Marshal result and return it.
	result := &pb.GetCryptoAddressFromAbelAddressResult{
		CryptoAddress: cryptoAddress.Data(),
	}
	return marshalResultAndPackToRetData(result)
}

//export GetShortAbelAddressFromAbelAddress
func GetShortAbelAddressFromAbelAddress(argsData []byte) *C.char {
	// Unmarshal args.
	args := &pb.GetShortAbelAddressFromAbelAddressArgs{}
	unmarshalArgs(argsData, args)

	// Prepare data.
	abelAddress := core.NewAbelAddress(args.GetAbelAddress())

	// Do real work.
	shortAbelAddress := abelAddress.GetShortAbelAddress()

	// Marshal result and return it.
	result := &pb.GetShortAbelAddressFromAbelAddressResult{
		ShortAbelAddress: shortAbelAddress.Data(),
	}
	return marshalResultAndPackToRetData(result)
}

//export DecodeFingerprintFromTxVoutScript
func DecodeFingerprintFromTxVoutScript(argsData []byte) *C.char {
	// Unmarshal args.
	args := &pb.DecodeFingerprintFromTxVoutScriptArgs{}
	unmarshalArgs(argsData, args)

	// Prepare data.
	txVoutScript := core.AsBytes(args.GetTxVoutScript())

	// Do real work.
	coinAddress, err := core.DecodeCoinAddressFromTxOutData(txVoutScript)
	panicIf(err)
	fingerprint := coinAddress.Fingerprint()

	// Marshal result and return it.
	result := &pb.DecodeFingerprintFromTxVoutScriptResult{
		Fingerprint: fingerprint.Slice(),
	}
	return marshalResultAndPackToRetData(result)
}

//export DecodeCoinValueFromTxVoutScript
func DecodeCoinValueFromTxVoutScript(argsData []byte) *C.char {
	// Unmarshal args.
	args := &pb.DecodeCoinValueFromTxVoutScriptArgs{}
	unmarshalArgs(argsData, args)

	// Prepare data.
	txVoutScript := core.AsBytes(args.GetTxVoutScript())
	viewSecretKey := core.NewCryptoKey(args.GetViewSecretKey())

	// Do real work.
	coinValue, err := core.DecodeValueFromTxOutData(txVoutScript, viewSecretKey)
	panicIf(err)

	// Marshal result and return it.
	result := &pb.DecodeCoinValueFromTxVoutScriptResult{
		CoinValue: coinValue,
	}
	return marshalResultAndPackToRetData(result)
}

//export GenerateUnsignedRawTxData
func GenerateUnsignedRawTxData(argsData []byte) *C.char {
	// Unmarshal args.
	args := &pb.GenerateUnsignedRawTxDataArgs{}
	unmarshalArgs(argsData, args)

	// Prepare data.
	txDesc := newTxDescFromArgs(args)

	// Do real work.
	unsignedRawTx, err := core.GenerateUnsignedRawTx(txDesc)
	panicIf(err)
	signerShortAddresses := make([][]byte, 0, len(unsignedRawTx.Signers))
	for _, signer := range unsignedRawTx.Signers {
		signerShortAddresses = append(signerShortAddresses, signer.Data())
	}

	// Marshal result and return it.
	result := &pb.GenerateUnsignedRawTxDataResult{
		Data:                 unsignedRawTx.Bytes,
		SignerShortAddresses: signerShortAddresses,
	}
	return marshalResultAndPackToRetData(result)
}

func newTxDescFromArgs(args *pb.GenerateUnsignedRawTxDataArgs) *core.TxDesc {
	// Create txInDescs.
	txInDescMessages := args.GetTxInDescs()
	txInDescs := make([]*core.TxInDesc, 0, len(txInDescMessages))
	for _, txInDescMessage := range txInDescMessages {
		txInDesc := newTxInDescFromMessage(txInDescMessage)
		txInDescs = append(txInDescs, txInDesc)
	}

	// Create txOutDescs.
	txOutDescMessages := args.GetTxOutDescs()
	txOutDescs := make([]*core.TxOutDesc, 0, len(txOutDescMessages))
	for _, txOutDescMessage := range txOutDescMessages {
		txOutDesc := newTxOutDescFromMessage(txOutDescMessage)
		txOutDescs = append(txOutDescs, txOutDesc)
	}

	// Create txRingBlockDescs.
	txRingBlockDescMessages := args.GetTxRingBlockDescs()
	txRingBlockDescs := make(map[int64]*core.TxBlockDesc)
	for _, txRingBlockDescMessage := range txRingBlockDescMessages {
		txRingBlockDesc := newTxRingBlockDescFromMessage(txRingBlockDescMessage)
		txRingBlockDescs[txRingBlockDesc.Height] = txRingBlockDesc
	}

	// Create txFee.
	txFee := args.GetTxFee()

	// Create and return txDesc.
	return core.NewTxDesc(txInDescs, txOutDescs, txFee, txRingBlockDescs)
}

func newTxInDescFromMessage(txInDescMessage *pb.TxInDescMessage) *core.TxInDesc {
	return &core.TxInDesc{
		TxOutData:  core.AsBytes(txInDescMessage.GetScript()),
		CoinValue:  txInDescMessage.GetValue(),
		Owner:      core.NewShortAbelAddress(txInDescMessage.GetShortAbelAddress()),
		Height:     txInDescMessage.GetHeight(),
		TxHash:     core.AsBytes(txInDescMessage.GetTxid()),
		TxOutIndex: uint8(txInDescMessage.GetIndex()),
	}
}

func newTxOutDescFromMessage(txOutDescMessage *pb.TxOutDescMessage) *core.TxOutDesc {
	return &core.TxOutDesc{
		AbelAddress: core.NewAbelAddress(txOutDescMessage.GetAbelAddress()),
		CoinValue:   txOutDescMessage.GetValue(),
	}
}

func newTxRingBlockDescFromMessage(txRingBlockDescMessage *pb.BlockDescMessage) *core.TxBlockDesc {
	return &core.TxBlockDesc{
		BinData: core.AsBytes(txRingBlockDescMessage.GetBinData()),
		Height:  txRingBlockDescMessage.GetHeight(),
	}
}

//export GenerateSignedRawTxData
func GenerateSignedRawTxData(argsData []byte) *C.char {
	// Unmarshal args.
	args := &pb.GenerateSignedRawTxDataArgs{}
	unmarshalArgs(argsData, args)

	// Prepare data.
	signerShortAddresses := args.GetSignerShortAddresses()
	signers := make([]*core.ShortAbelAddress, 0, len(signerShortAddresses))
	for _, signerShortAddress := range signerShortAddresses {
		signers = append(signers, core.NewShortAbelAddress(signerShortAddress))
	}

	unsignedRawTx := &core.UnsignedRawTx{
		Bytes:   core.AsBytes(args.GetUnsignedRawTxData()),
		Signers: signers,
	}

	signerKeys := make([]*core.CryptoKeysAndAddress, 0, len(unsignedRawTx.Signers))
	for _, cryptoSeed := range args.GetSignerCryptoSeeds() {
		ckaa, err := core.GenerateCryptoKeysAndAddress(cryptoSeed)
		panicIf(err)
		signerKeys = append(signerKeys, ckaa)
	}

	// Do real work.
	signedRawTx, err := core.GenerateSignedRawTx(unsignedRawTx, signerKeys)
	panicIf(err)

	// Marshal result and return it.
	result := &pb.GenerateSignedRawTxDataResult{
		Data: signedRawTx.Bytes,
		Txid: signedRawTx.Txid,
	}
	return marshalResultAndPackToRetData(result)
}

//export GenerateCoinSerialNumber
func GenerateCoinSerialNumber(argsData []byte) *C.char {
	// Unmarshal args.
	args := &pb.GenerateCoinSerialNumberArgs{}
	unmarshalArgs(argsData, args)

	// Prepare data.
	coinID := core.NewCoinID(core.AsBytes(args.GetTxid()), uint8(args.GetIndex()))
	serialNoSecretKey := core.NewCryptoKey(args.GetSerialNoSecretKey())
	ringBlockDescs := make(map[int64]*core.TxBlockDesc)
	for _, ringBlockDescMessage := range args.GetRingBlockDescs() {
		ringBlockDesc := &core.TxBlockDesc{
			BinData: core.AsBytes(ringBlockDescMessage.GetBinData()),
			Height:  ringBlockDescMessage.GetHeight(),
		}
		ringBlockDescs[ringBlockDesc.Height] = ringBlockDesc
	}

	// Do real work.
	serialNumbers, err := core.DecodeCoinSerialNumbers([]*core.CoinID{coinID}, []*core.CryptoKey{serialNoSecretKey}, ringBlockDescs)
	panicIf(err)

	// Marshal result and return it.
	result := &pb.GenerateCoinSerialNumberResult{
		SerialNumber: serialNumbers[0].Slice(),
	}
	return marshalResultAndPackToRetData(result)
}

package transaction

type UTXOInfo struct {
	TXID   []byte
	Index  int64
	Output TXOutput
}

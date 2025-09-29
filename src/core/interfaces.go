package core

// ITransactionFactory mendefinisikan antarmuka untuk membuat berbagai jenis transaksi.
// Ini memisahkan tanggung jawab pembuatan transaksi dari struct Blockchain itu sendiri.

type ITransactionFactory interface {
	NewUTXOTransaction(from, to string, amount int, UTXOSet *UTXOSet) (*Transaction, error)
	NewContractCreationTransaction(from string, code []byte, UTXOSet *UTXOSet) (*Transaction, error)
	NewContractCallTransaction(from, contractAddress, function, args string, UTXOSet *UTXOSet) (*Transaction, error)
	NewStakeTransaction(from string, amount int, UTXOSet *UTXOSet) (*Transaction, error)
	NewInferenceRequestTransaction(from, modelName, inputData string, reward int, UTXOSet *UTXOSet) (*Transaction, error)
}

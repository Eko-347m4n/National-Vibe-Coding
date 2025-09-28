package core

// transaction.go akan berisi definisi struktur Transaction (UTXO).
// Termasuk input, output, dan logika penandatanganan.

// Transaction merepresentasikan sebuah transaksi
type Transaction struct {
	ID      []byte     // ID transaksi (hash dari data transaksi)
	Inputs  []TXInput  // Input transaksi
	Outputs []TXOutput // Output transaksi
}

// TXInput merepresentasikan sebuah input transaksi
// Setiap input harus merujuk ke sebuah output dari transaksi sebelumnya (UTXO)
type TXInput struct {
	TxID      []byte // ID transaksi sebelumnya yang outputnya akan digunakan
	OutIndex  int    // Indeks output di dalam transaksi tersebut
	Signature []byte // Tanda tangan digital (placeholder)
	PubKey    []byte // Kunci publik (placeholder)
}

// TXOutput merepresentasikan sebuah output transaksi
// Setiap output berisi 'koin' yang belum dibelanjakan
type TXOutput struct {
	Value      int    // Nilai koin
	PubKeyHash []byte // Hash dari kunci publik pemilik
}

// TODO: Implementasikan metode untuk menandatangani dan memverifikasi transaksi
// TODO: Implementasikan logika untuk mengunci dan membuka output (locking/unlocking script)

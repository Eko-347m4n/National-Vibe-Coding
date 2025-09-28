# Referensi CLI

Berikut adalah daftar lengkap perintah yang tersedia di `swatantra-node`.

### Perintah Dasar & Wallet

- `createwallet`
  - Membuat pasangan kunci publik-privat baru dan menyimpannya.
- `listaddresses`
  - Menampilkan semua alamat yang tersimpan di file wallet.
- `getbalance -address <ADDRESS>`
  - Menampilkan saldo UTXO dari sebuah alamat.

### Operasi Blockchain

- `createblockchain -address <ADDRESS>`
  - Membuat blockchain baru dan mengirim hadiah genesis ke alamat yang ditentukan.
- `printchain`
  - Mencetak semua blok yang ada di blockchain dari yang terbaru hingga genesis.
- `send -from <FROM> -to <TO> -amount <AMOUNT>`
  - Mengirim sejumlah koin dari satu alamat ke alamat lain.

### Smart Contracts

- `deploycontract -from <FROM> -file <FILE_PATH>`
  - Menerbitkan smart contract baru (.lua) ke blockchain.
- `callcontract -from <FROM> -contract <CONTRACT_ADDR> -function <FUNC> -args "<ARGS>"`
  - Memanggil sebuah fungsi di dalam smart contract yang sudah ada. Argumen dipisahkan koma.

### AI Oracle & Staking

- `registernode -from <FROM> -contract <REGISTRY_ADDR>`
  - Mendaftarkan alamat `-from` sebagai inference node ke dalam kontrak registri.
- `publishmodel -from <FROM> -contract <REGISTRY_ADDR> -name <NAME> -file <FILE>`
  - Mempublikasikan hash dari sebuah file model ke kontrak registri model.
- `requestinference -from <FROM> -contract <MARKET_ADDR> -model <NAME> -input <DATA>`
  - Membuat pekerjaan inferensi baru di papan pekerjaan.
- `submitresponse -from <FROM> -contract <MARKET_ADDR> -jobid <ID> -result <DATA>`
  - Mengirimkan hasil untuk sebuah pekerjaan inferensi.
- `getjob -contract <MARKET_ADDR> -jobid <ID>`
  - Mengambil detail dan semua respons dari sebuah pekerjaan inferensi.

### Jaringan

- `startnode -miner <ADDRESS>`
  - Memulai node dan (opsional) mengaktifkan mode mining, dengan hadiah dikirim ke alamat yang ditentukan.

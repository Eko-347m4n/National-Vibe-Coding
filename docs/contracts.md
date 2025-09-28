# Referensi Smart Contract

Ekosistem Swatantra menggunakan beberapa smart contract inti untuk mengelola operasinya. Berikut adalah penjelasan singkatnya.

## `oracle_registry.lua`

Kontrak ini berfungsi sebagai buku catatan on-chain untuk semua node yang ingin berpartisipasi sebagai *inference node*.

### Fungsi Utama
- **`register(node_address)`**: Mendaftarkan `node_address` ke dalam daftar node aktif.
- **`get_nodes()`**: Mengembalikan daftar semua alamat node yang terdaftar, dipisahkan oleh koma.

## `model_registry.lua`

Kontrak ini melacak model-model AI yang "resmi" dan dapat digunakan di dalam jaringan. Ia menyimpan hash dan lokasi penyimpanan dari setiap model.

### Fungsi Utama
- **`publish_model(name, model_hash, location)`**: Mencatat model baru. `name` adalah ID unik, `model_hash` adalah hash Keccak256 dari file model, dan `location` adalah tempat file tersebut bisa diunduh (misal: CID IPFS).
- **`get_model(name)`**: Mengambil `model_hash` dan `location` untuk model dengan nama yang diberikan.

## `inference_market.lua`

Berfungsi sebagai "papan pekerjaan" terdesentralisasi untuk tugas-tugas inferensi.

### Fungsi Utama
- **`request_inference(model_name, input_data, reward)`**: Membuat pekerjaan baru, meminta inferensi pada `model_name` dengan `input_data` dan menawarkan `reward`.
- **`submit_response(job_id, result, node_address)`**: Dipanggil oleh *inference node* untuk mengirimkan `result` mereka untuk `job_id` tertentu.
- **`get_job(job_id)`**: Mengambil semua detail pekerjaan, termasuk semua respons yang telah dikirimkan oleh para node.

package cli

import (
	// Impor paket cmd untuk mendaftarkan semua perintah yang tersedia.
	// Underscore `_` digunakan untuk memastikan fungsi init() di dalam paket cmd
	// dan sub-paketnya dieksekusi, yang akan memanggil AddCommand().
	_ "swatantra-node/src/cli/cmd"
	"swatantra-node/src/cli/cmd"
)

// CLI sekarang hanya berfungsi sebagai placeholder untuk fungsi Run.
type CLI struct{}

// Run adalah titik masuk utama yang mendelegasikan eksekusi ke command runner baru.
func (cli *CLI) Run() {
	cmd.Execute()
}

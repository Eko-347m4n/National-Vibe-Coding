package vm

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

// TestSimpleContractExecution menguji alur eksekusi kontrak yang paling dasar.
func TestSimpleContractExecution(t *testing.T) {
	// Kode kontrak Lua sederhana untuk diuji.
	// 1. Menyimpan sebuah nilai.
	// 2. Membaca nilai tersebut kembali.
	// 3. Mengembalikan nilai yang sudah ditambah 1.
	contractCode := `
        log("Menjalankan kontrak sederhana...")
        
        -- Simpan nilai awal
        set_value("counter", 10)
        
        -- Baca nilai yang baru disimpan
        local val = get_value("counter")
        log("Nilai yang dibaca dari state: " .. tostring(val))
        
        -- Lakukan operasi dan simpan kembali
        local new_val = val + 1
        set_value("counter", new_val)

        log("Eksekusi kontrak selesai.")
    `

	// Buat VM baru
	vm := NewVM()
	defer vm.Close()

	// Eksekusi kode
	err := vm.Execute(contractCode)
	if err != nil {
		t.Fatalf("Eksekusi kontrak gagal: %v", err)
	}

	// Verifikasi state setelah eksekusi
	expectedValue := lua.LNumber(11)
	actualValue, ok := vm.State["counter"]

	if !ok {
		t.Fatal("State 'counter' tidak ditemukan setelah eksekusi kontrak.")
	}

	if actualValue != expectedValue {
		t.Fatalf("Nilai state tidak sesuai. Harusnya %s, tapi dapat %s", expectedValue, actualValue)
	}

	t.Logf("Verifikasi state berhasil. Nilai 'counter' adalah: %s", actualValue)
}

// TestHostFunctions menguji fungsionalitas host functions secara individual
func TestHostFunctions(t *testing.T) {
	vm := NewVM()
	defer vm.Close()

	// Uji set_value dan get_value
	err := vm.L.DoString(`set_value("test_key", "test_value")`)
	if err != nil {
		t.Fatalf("Gagal menjalankan set_value: %v", err)
	}

	val, ok := vm.State["test_key"]
	if !ok || val.(lua.LString) != "test_value" {
		t.Fatal("set_value tidak berhasil menyimpan state dengan benar.")
	}
	t.Log("Fungsi host set_value berhasil diuji.")

	err = vm.L.DoString(`v = get_value("test_key")`)
	if err != nil {
		t.Fatalf("Gagal menjalankan get_value: %v", err)
	}

	retrievedVal := vm.L.GetGlobal("v")
	if retrievedVal.(lua.LString) != "test_value" {
		t.Fatal("get_value tidak berhasil mengambil state dengan benar.")
	}
	t.Log("Fungsi host get_value berhasil diuji.")
}

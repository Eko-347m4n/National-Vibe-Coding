package vm

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

// VM mewakili instance dari Virtual Machine untuk smart contract
type VM struct {
	L         *lua.LState
	GasLimit  uint64
	GasUsed   uint64
	State     map[string]string // State sementara (key-value string) untuk satu eksekusi
}

// NewVM membuat instance VM baru
func NewVM() *VM {
	L := lua.NewState()

	vm := &VM{
		L:     L,
		State: make(map[string]string),
	}

	// Daftarkan fungsi Go yang bisa dipanggil dari Lua (host functions)
	vm.registerHostFunctions()

	return vm
}

// Close menutup state Lua untuk membersihkan memori
func (vm *VM) Close() {
	vm.L.Close()
}

// registerHostFunctions mendaftarkan fungsi-fungsi Go ke environment Lua
func (vm *VM) registerHostFunctions() {
	// Fungsi untuk menyimpan state (key-value)
	vm.L.SetGlobal("set_value", vm.L.NewFunction(vm.hostSetValue))
	// Fungsi untuk membaca state
	vm.L.SetGlobal("get_value", vm.L.NewFunction(vm.hostGetValue))
	// Fungsi untuk logging (berguna untuk debugging)
	vm.L.SetGlobal("log", vm.L.NewFunction(vm.hostLog))
}

// Execute menjalankan kode smart contract
func (vm *VM) Execute(contractCode string) error {
	// Di masa depan, kita akan menambahkan gas limit di sini
	// vm.L.SetMx(int(vm.GasLimit))

	return vm.L.DoString(contractCode)
}

// --- Host Functions ---

// hostSetValue memungkinkan kontrak Lua untuk menyimpan data di state-nya
// Penggunaan di Lua: set_value("my_key", "my_value")
func (vm *VM) hostSetValue(L *lua.LState) int {
	key := L.ToString(1)   // Argumen pertama
	value := L.ToString(2) // Argumen kedua (sekarang hanya string)

	if key == "" {
		L.Push(lua.LString("Error: key tidak boleh kosong."))
		return 1 // Mengembalikan 1 nilai (pesan error)
	}

	vm.State[key] = value
	return 0 // Tidak mengembalikan apa-apa
}

// hostGetValue memungkinkan kontrak Lua untuk membaca data dari state-nya
// Penggunaan di Lua: local val = get_value("my_key")
func (vm *VM) hostGetValue(L *lua.LState) int {
	key := L.ToString(1)
	value, ok := vm.State[key]
	if !ok {
		L.Push(lua.LNil)
	} else {
		L.Push(lua.LString(value))
	}
	return 1 // Mengembalikan 1 nilai (value atau nil)
}

// hostLog memungkinkan kontrak Lua untuk mencetak log
// Penggunaan di Lua: log("Hello from contract!")
func (vm *VM) hostLog(L *lua.LState) int {
	message := L.ToString(1)
	fmt.Printf("[CONTRACT LOG] %s\n", message)
	return 0
}

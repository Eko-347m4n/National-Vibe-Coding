document.addEventListener('DOMContentLoaded', () => {

    const getContractAddress = () => {
        const addr = document.getElementById('contract-address').value;
        if (!addr) {
            alert('Silakan masukkan alamat kontrak terlebih dahulu di bagian Setup.');
            return null;
        }
        return addr;
    };

    // --- Daftarkan Aset ---
    document.getElementById('create-asset-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const assetId = document.getElementById('asset-id').value;
        const assetType = document.getElementById('asset-type').value;
        if (!assetId || !assetType) {
            alert('ID dan Tipe Aset tidak boleh kosong.');
            return;
        }

        const args = `${assetId},${assetType},[ALAMAT_ANDA]`;
        const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "create_asset" -args "${args}"`;
        document.getElementById('asset-command').textContent = command;
        document.getElementById('asset-output').textContent = "Jalankan perintah di atas. Outputnya adalah ID Internal aset baru Anda (misal: 1). Simpan ID ini!";
    });

    // --- Catat Peristiwa ---
    document.getElementById('log-event-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const internalId = document.getElementById('event-asset-id').value;
        const eventName = document.getElementById('event-name').value;
        const eventDetails = document.getElementById('event-details').value;
        if (!internalId || !eventName || !eventDetails) {
            alert('ID Internal, Nama Peristiwa, dan Detail tidak boleh kosong.');
            return;
        }

        const args = `${internalId},${eventName},[ALAMAT_ANDA],${eventDetails}`;
        const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "log_event" -args "${args}"`;
        document.getElementById('event-command').textContent = command;
        document.getElementById('event-output').textContent = `Jalankan perintah di atas untuk mencatat peristiwa '${eventName}'.`;
    });

    // --- Dapatkan Riwayat ---
    document.getElementById('get-history-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const internalId = document.getElementById('history-asset-id').value;
        if (!internalId) {
            alert('ID Internal Aset tidak boleh kosong.');
            return;
        }

        const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "get_asset_history" -args "${internalId}"`;
        document.getElementById('history-command').textContent = command;
        document.getElementById('history-output').textContent = "Jalankan perintah di atas. Output JSON akan muncul di log, berisi detail lengkap dan riwayat peristiwa aset.";
    });
});

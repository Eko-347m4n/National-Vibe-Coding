document.addEventListener('DOMContentLoaded', () => {
    // Inisialisasi ClipboardJS
    new ClipboardJS('.copy-btn');

    // --- Helper Functions ---
    const getContractAddress = () => {
        const addr = document.getElementById('contract-address').value.trim();
        if (!addr) {
            showError('contract-address', 'Alamat kontrak tidak boleh kosong.');
            return null;
        }
        hideError('contract-address');
        return addr;
    };

    const showError = (inputId, message) => {
        const errorEl = document.getElementById(`${inputId}-error`);
        const inputEl = document.getElementById(inputId);
        if (errorEl) {
            errorEl.textContent = message;
            errorEl.style.display = 'block';
        }
        if (inputEl) {
            inputEl.style.borderColor = 'var(--error-color)';
            inputEl.focus();
        }
    };

    const hideError = (inputId) => {
        const errorEl = document.getElementById(`${inputId}-error`);
        const inputEl = document.getElementById(inputId);
        if (errorEl) {
            errorEl.style.display = 'none';
        }
        if (inputEl) {
            inputEl.style.borderColor = 'var(--border-color)';
        }
    };

    const showLoader = (btnId) => {
        const btn = document.getElementById(btnId);
        if (btn) {
            btn.disabled = true;
            btn.querySelector('.loader').style.display = 'inline-block';
        }
    };

    const hideLoader = (btnId) => {
        const btn = document.getElementById(btnId);
        if (btn) {
            btn.disabled = false;
            btn.querySelector('.loader').style.display = 'none';
        }
    };
    
    const displayCommand = (commandElId, outputElId, command, output) => {
        const commandEl = document.getElementById(commandElId);
        const outputEl = document.getElementById(outputElId);
        
        commandEl.textContent = command;
        outputEl.textContent = output;

        // Make parent wrappers visible
        commandEl.closest('.code-block-wrapper').style.display = 'block';
        outputEl.style.display = 'block';
    };

    // --- Event Listeners ---
    document.querySelectorAll('input[type="text"]').forEach(input => {
        input.addEventListener('input', () => hideError(input.id));
    });

    document.querySelectorAll('.copy-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            const originalText = btn.textContent;
            btn.textContent = 'Disalin!';
            setTimeout(() => {
                btn.textContent = originalText;
            }, 1500);
        });
    });

    // --- Daftarkan Aset ---
    document.getElementById('create-asset-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const assetId = document.getElementById('asset-id').value.trim();
        const assetType = document.getElementById('asset-type').value.trim();

        if (!assetId) {
            showError('asset-id', 'ID Aset tidak boleh kosong.');
            return;
        }
        if (!assetType) {
            showError('asset-type', 'Tipe Aset tidak boleh kosong.');
            return;
        }

        showLoader('create-asset-btn');
        setTimeout(() => {
            const args = `${assetId},${assetType},[ALAMAT_ANDA]`;
            const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "create_asset" -args "${args}"`;
            const output = "Jalankan perintah di atas. Outputnya adalah ID Internal aset baru Anda (misal: 1). Simpan ID ini!";
            displayCommand('asset-command', 'asset-output', command, output);
            hideLoader('create-asset-btn');
        }, 500); // Simulasi loading
    });

    // --- Catat Peristiwa ---
    document.getElementById('log-event-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const internalId = document.getElementById('event-asset-id').value.trim();
        const eventName = document.getElementById('event-name').value.trim();
        const eventDetails = document.getElementById('event-details').value.trim();

        if (!internalId) {
            showError('event-asset-id', 'ID Internal Aset tidak boleh kosong.');
            return;
        }
        if (!eventName) {
            showError('event-name', 'Nama Peristiwa tidak boleh kosong.');
            return;
        }
        if (!eventDetails) {
            showError('event-details', 'Detail Peristiwa tidak boleh kosong.');
            return;
        }

        showLoader('log-event-btn');
        setTimeout(() => {
            const args = `${internalId},${eventName},[ALAMAT_ANDA],${eventDetails}`;
            const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "log_event" -args "${args}"`;
            const output = `Jalankan perintah di atas untuk mencatat peristiwa '${eventName}'.`;
            displayCommand('event-command', 'event-output', command, output);
            hideLoader('log-event-btn');
        }, 500);
    });

    // --- Dapatkan Riwayat ---
    document.getElementById('get-history-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const internalId = document.getElementById('history-asset-id').value.trim();
        if (!internalId) {
            showError('history-asset-id', 'ID Internal Aset tidak boleh kosong.');
            return;
        }

        showLoader('get-history-btn');
        setTimeout(() => {
            const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "get_asset_history" -args "${internalId}"`;
            const output = "Jalankan perintah di atas. Output JSON akan muncul di log, berisi detail lengkap dan riwayat peristiwa aset.";
            displayCommand('history-command', 'history-output', command, output);
            hideLoader('get-history-btn');
        }, 500);
    });
});
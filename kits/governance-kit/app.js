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

    // --- Anggaran ---
    document.getElementById('get-budget-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        showLoader('get-budget-btn');
        setTimeout(() => {
            const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "get_budget" -args ""`;
            const output = "Jalankan perintah di atas di terminal Anda. Output JSON akan muncul di log, contoh: {\"total\":10000, \"spent\":0, \"remaining\":10000}";
            displayCommand('budget-command', 'budget-output', command, output);
            hideLoader('get-budget-btn');
        }, 500); // Simulasi loading
    });

    // --- Buat Proposal ---
    document.getElementById('create-proposal-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const title = document.getElementById('proposal-title').value.trim();
        const desc = document.getElementById('proposal-desc').value.trim();

        if (!title) {
            showError('proposal-title', 'Judul proposal tidak boleh kosong.');
            return;
        }
        if (!desc) {
            showError('proposal-desc', 'Deskripsi tidak boleh kosong.');
            return;
        }

        showLoader('create-proposal-btn');
        setTimeout(() => {
            const args = `${title},${desc}`;
            const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "create_proposal" -args "${args}"`;
            const output = "Jalankan perintah di atas. Outputnya adalah ID proposal baru Anda.";
            displayCommand('proposal-command', 'proposal-output', command, output);
            hideLoader('create-proposal-btn');
        }, 500);
    });

    // --- Vote ---
    const handleVote = (voteType, btnId) => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const proposalId = document.getElementById('vote-proposal-id').value.trim();
        if (!proposalId) {
            showError('vote-proposal-id', 'ID Proposal tidak boleh kosong.');
            return;
        }

        showLoader(btnId);
        setTimeout(() => {
            const args = `${proposalId},${voteType},[ALAMAT_ANDA]`;
            const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "vote" -args "${args}"`;
            const output = `Jalankan perintah di atas untuk memberikan suara ${voteType}.`;
            displayCommand('vote-command', 'vote-output', command, output);
            hideLoader(btnId);
        }, 500);
    };

    document.getElementById('vote-for-btn').addEventListener('click', () => handleVote('FOR', 'vote-for-btn'));
    document.getElementById('vote-against-btn').addEventListener('click', () => handleVote('AGAINST', 'vote-against-btn'));

    // --- Analisis Anomali ---
    document.getElementById('check-anomaly-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const txData = document.getElementById('tx-data').value.trim();
        if (!txData) {
            showError('tx-data', 'Data transaksi tidak boleh kosong.');
            return;
        }

        showLoader('check-anomaly-btn');
        setTimeout(() => {
            const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "check_anomaly" -args "${txData}"`;
            const output = "Jalankan perintah di atas. Outputnya adalah hasil analisis dari AI Oracle (simulasi).";
            displayCommand('anomaly-command', 'anomaly-output', command, output);
            hideLoader('check-anomaly-btn');
        }, 500);
    });
});

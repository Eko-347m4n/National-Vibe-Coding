document.addEventListener('DOMContentLoaded', () => {

    const getContractAddress = () => {
        const addr = document.getElementById('contract-address').value;
        if (!addr) {
            alert('Silakan masukkan alamat kontrak terlebih dahulu di bagian Setup.');
            return null;
        }
        return addr;
    };

    // --- Anggaran ---
    document.getElementById('get-budget-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "get_budget" -args ""`;
        document.getElementById('budget-command').textContent = command;
        document.getElementById('budget-output').textContent = "Jalankan perintah di atas di terminal Anda. Output JSON akan muncul di log, contoh: {\"total\":10000, \"spent\":0, \"remaining\":10000}";
    });

    // --- Buat Proposal ---
    document.getElementById('create-proposal-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const title = document.getElementById('proposal-title').value;
        const desc = document.getElementById('proposal-desc').value;
        if (!title || !desc) {
            alert('Judul dan deskripsi proposal tidak boleh kosong.');
            return;
        }

        const args = `${title},${desc}`;
        const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "create_proposal" -args "${args}"`;
        document.getElementById('proposal-command').textContent = command;
        document.getElementById('proposal-output').textContent = "Jalankan perintah di atas. Outputnya adalah ID proposal baru Anda.";
    });

    // --- Vote ---
    const handleVote = (voteType) => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const proposalId = document.getElementById('vote-proposal-id').value;
        if (!proposalId) {
            alert('ID Proposal tidak boleh kosong.');
            return;
        }

        const args = `${proposalId},${voteType},[ALAMAT_ANDA]`;
        const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "vote" -args "${args}"`;
        document.getElementById('vote-command').textContent = command;
        document.getElementById('vote-output').textContent = `Jalankan perintah di atas untuk memberikan suara ${voteType}.`;
    };

    document.getElementById('vote-for-btn').addEventListener('click', () => handleVote('FOR'));
    document.getElementById('vote-against-btn').addEventListener('click', () => handleVote('AGAINST'));

    // --- Analisis Anomali ---
    document.getElementById('check-anomaly-btn').addEventListener('click', () => {
        const contractAddr = getContractAddress();
        if (!contractAddr) return;

        const txData = document.getElementById('tx-data').value;
        if (!txData) {
            alert('Data transaksi tidak boleh kosong.');
            return;
        }

        const command = `./swatantra-node callcontract -from [ALAMAT_ANDA] -contract ${contractAddr} -function "check_anomaly" -args "${txData}"`;
        document.getElementById('anomaly-command').textContent = command;
        document.getElementById('anomaly-output').textContent = "Jalankan perintah di atas. Outputnya adalah hasil analisis dari AI Oracle (simulasi).";
    });
});

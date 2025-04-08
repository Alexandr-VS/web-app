function toggleFileInput() {
    const fileInput = document.getElementById('filename');
    const pseudoRandOption = document.getElementById('pseudoRand');
    const packetSizeInput = document.getElementById('packetSize');
    const packetSizeLabel = document.getElementById('packetSizeLabel');

    if (pseudoRandOption.checked) {
        fileInput.style.display = 'none';
        packetSizeInput.style.display = 'block';
        packetSizeLabel.style.display = 'block';

    } else {
        fileInput.style.display = 'block';
        packetSizeInput.style.display = 'none';
        packetSizeLabel.style.display = 'none';
    }
}

function validateIP(input) {
    const ipPattern = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
    if (ipPattern.test(input.value)) {
        input.classList.add('valid');
        input.classList.remove('invalid');
    } else {
        input.classList.add('invalid');
        input.classList.remove('valid');
    }
}

function validateMAC(input) {
    const macPattern = /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/;
    if (macPattern.test(input.value)) {
        input.classList.add('valid');
        input.classList.remove('invalid');
    } else {
        input.classList.add('invalid');
        input.classList.remove('valid');
    }
}

function validatePort(input){
    const port = parseInt(input.value, 10);
    if (!isNaN(port) && Number.isInteger(port) && port >= 0 && port <= 65535) {
        input.classList.add('valid');
        input.classList.remove('invalid')
    } else {
        input.classList.add('invalid');
        input.classList.remove('valid');
    }
}

function checkInputsOnLoad() {
    const ipSrcInput = document.getElementById('ip-src');
    const ipDstInput = document.getElementById('ip-dst');
    const macSrcInput = document.getElementById('mac-src');
    const macDstInput = document.getElementById('mac-dst');
    const portSrcInput = document.getElementById('src-port');
    const portDstInput = document.getElementById('dst-port');

    validateIP(ipSrcInput);
    validateIP(ipDstInput);
    validateMAC(macSrcInput);
    validateMAC(macDstInput);
    validatePort(portSrcInput);
    validatePort(portDstInput);
}

document.getElementById('interval').addEventListener('input', function (event) {
    this.value = this.value
        .replace(/,/g, '.') // Заменяем запятую на точку
        .trim(); // Удаляем пробелы в начале и в конце
});

document.addEventListener('DOMContentLoaded', () => {
    checkInputsOnLoad();
    toggleFileInput();
});
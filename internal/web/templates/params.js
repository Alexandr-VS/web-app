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
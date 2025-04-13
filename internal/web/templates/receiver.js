// Функция для получения пакетов с сервера
async function fetchPackets() {
    try {
        const response = await fetch('/receiver');
        if (!response.ok) {
            throw new Error('Ошибка сети');
        }
        const packets = await response.json();

        const packetsList = document.getElementById('packets-list');
        packetsList.innerHTML = ''; // Очищаем текущий список

        packets.forEach(packet => {
            const ul = document.createElement('ul');
            ul.textContent = `Пакет #${packet.counter}: Время задержки пакета: ${packet.delay}`;
            packetsList.appendChild(ul);
        });
    } catch (error) {
        console.error('Ошибка при получении пакетов:', error);
    }
}
// Обновление списка пакетов
setInterval(fetchPackets, 3000);

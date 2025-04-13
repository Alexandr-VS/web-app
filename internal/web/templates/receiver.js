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

// Функция для проверки завершения приема пакетов
async function checkCompletion() {
    try {
        const response = await fetch('/check-completion');
        if (!response.ok) {
            throw new Error('Ошибка сети при проверке завершения');
        }

        const { completed } = await response.json();
        if (completed) {
            // Отображаем кнопку "Показать отчет", если прием завершен
            document.getElementById('reportButton').style.display = 'block';
            // Останавливаем дальнейшую проверку состояния
            clearInterval(completionInterval);
        }
    } catch (error) {
        console.error('Ошибка при проверке завершения:', error);
    }
}

// Обновление списка пакетов каждые 3 секунды
setInterval(fetchPackets, 3000);

// Проверка на завершение приема каждые 5 секунд
const completionInterval = setInterval(checkCompletion, 5000);

// Обработка нажатия на кнопку отчета
document.getElementById('reportButton').addEventListener('click', function() {
    window.location.href = '/report'; // Переход на страницу отчета
});
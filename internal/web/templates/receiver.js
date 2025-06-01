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


        if (packets.length === 0) {
            const tr = document.createElement('tr');
            const td = document.createElement('td');
            td.colSpan = 5; // Количество колонок в таблице
            td.textContent = 'Пакеты не получены';
            tr.appendChild(td);
            packetsList.appendChild(tr);
            return;
        }

        packets.forEach(packet => {
            const tr = document.createElement('tr');

            const tdCounter = document.createElement('td');
            tdCounter.textContent = packet.counter;
            tr.appendChild(tdCounter);

            const tdTotalDelay = document.createElement('td');
            tdTotalDelay.textContent = packet.totalDelay;
            tr.appendChild(tdTotalDelay);

            const tdForwardDelay = document.createElement('td');
            tdForwardDelay.textContent = packet.forwardDelay;
            tr.appendChild(tdForwardDelay);

            const tdBackwardDelay = document.createElement('td');
            tdBackwardDelay.textContent = packet.backwardDelay;
            tr.appendChild(tdBackwardDelay);

            const tdInterArrival = document.createElement('td');
            tdInterArrival.textContent = packet.interArrival;
            tr.appendChild(tdInterArrival);

            packetsList.appendChild(tr);
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
// Основные функции JavaScript
document.addEventListener('DOMContentLoaded', function() {
    // Автоскрытие уведомлений
    setTimeout(function() {
        const alerts = document.querySelectorAll('.alert');
        alerts.forEach(alert => {
            const bsAlert = new bootstrap.Alert(alert);
            bsAlert.close();
        });
    }, 5000);

    // Подтверждение удаления
    const deleteLinks = document.querySelectorAll('a[href*="/delete/"]');
    deleteLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            if (!confirm('Вы уверены, что хотите удалить этот элемент?')) {
                e.preventDefault();
            }
        });
    });

    // Динамическая загрузка участков при выборе цеха
    const hallSelect = document.getElementById('hall_id');
    const areaSelect = document.getElementById('area_id');
    const teamSelect = document.getElementById('work_team_id');

    if (hallSelect && areaSelect) {
        hallSelect.addEventListener('change', function() {
            const hallId = this.value;
            if (hallId) {
                fetch(`/api/areas/${hallId}`)
                    .then(response => response.json())
                    .then(data => {
                        areaSelect.innerHTML = '<option value="">Выберите участок</option>';
                        data.forEach(area => {
                            const option = document.createElement('option');
                            option.value = area.id;
                            option.textContent = area.name;
                            areaSelect.appendChild(option);
                        });
                        
                        // Очистить список бригад
                        if (teamSelect) {
                            teamSelect.innerHTML = '<option value="">Выберите бригаду</option>';
                        }
                    })
                    .catch(error => console.error('Ошибка:', error));
            } else {
                areaSelect.innerHTML = '<option value="">Выберите участок</option>';
                if (teamSelect) {
                    teamSelect.innerHTML = '<option value="">Выберите бригаду</option>';
                }
            }
        });
    }

    // Динамическая загрузка бригад при выборе участка
    if (areaSelect && teamSelect) {
        areaSelect.addEventListener('change', function() {
            const areaId = this.value;
            if (areaId) {
                fetch(`/api/teams/${areaId}`)
                    .then(response => response.json())
                    .then(data => {
                        teamSelect.innerHTML = '<option value="">Выберите бригаду</option>';
                        data.forEach(team => {
                            const option = document.createElement('option');
                            option.value = team.id;
                            option.textContent = team.name;
                            teamSelect.appendChild(option);
                        });
                    })
                    .catch(error => console.error('Ошибка:', error));
            } else {
                teamSelect.innerHTML = '<option value="">Выберите бригаду</option>';
            }
        });
    }

    // Анимация появления карточек
    const cards = document.querySelectorAll('.card');
    cards.forEach(card => {
        card.classList.add('fade-in');
    });
});

// Функция для форматирования дат
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU');
}

// Функция для поиска в таблицах
function searchTable(inputId, tableId) {
    const input = document.getElementById(inputId);
    const table = document.getElementById(tableId);
    
    if (input && table) {
        input.addEventListener('keyup', function() {
            const filter = this.value.toLowerCase();
            const rows = table.getElementsByTagName('tr');
            
            for (let i = 1; i < rows.length; i++) {
                const row = rows[i];
                const cells = row.getElementsByTagName('td');
                let found = false;
                
                for (let j = 0; j < cells.length; j++) {
                    if (cells[j].textContent.toLowerCase().includes(filter)) {
                        found = true;
                        break;
                    }
                }
                
                row.style.display = found ? '' : 'none';
            }
        });
    }
}

-- Заполнение таблицы testing_laboratory
INSERT INTO "testing_laboratory" ("id", "name") VALUES
(1, 'Лаборатория контроля качества'),
(2, 'Лаборатория испытаний материалов');

-- Заполнение таблицы production_halls
INSERT INTO "production_halls" ("id", "name") VALUES
(1, 'Сборочный цех №1'),
(2, 'Механический цех'),
(3, 'Испытательный цех');

-- Заполнение таблицы production_area
INSERT INTO "production_area" ("id", "name", "hall_id") VALUES
(1, 'Участок сборки А', 1),
(2, 'Участок токарной обработки', 2),
(3, 'Участок испытаний готовой продукции', 3);

-- Заполнение таблицы category_item
INSERT INTO "category_item" ("id", "name", "attribute") VALUES
(1, 'Электроника', 'Высокоточная'),
(2, 'Механические узлы', 'Промышленные'),
(3, 'Корпусные детали', NULL);

-- Заполнение таблицы category_engineer
INSERT INTO "category_engineer" ("id", "name", "attribute") VALUES
(1, 'Инженер-конструктор', 'Высшая категория'),
(2, 'Инженер-технолог', 'Первая категория'),
(3, 'Техник', NULL);

-- Заполнение таблицы type_item
INSERT INTO "type_item" ("id", "name", "category_id") VALUES
(1, 'Микросхема управления', 1),
(2, 'Редуктор малый', 2),
(3, 'Передняя панель', 3),
(4, 'Двигатель постоянного тока', 1);

-- Заполнение таблицы item
INSERT INTO "item" ("id", "name", "type_id", "hall_id", "status") VALUES
(1, 'Изделие Альфа-001', 1, 1, 'in_progress'),
(2, 'Изделие Бета-002', 2, 2, 'testing'),
(3, 'Изделие Гамма-003', 3, 1, 'completed'),
(4, 'Изделие Дельта-004', 4, 2, 'in_progress');

-- Заполнение таблицы employee
INSERT INTO "employee" ("id", "name", "hire_date", "current_position") VALUES
(1, 'Иванов Иван Иванович', '2020-05-15', 'Сборщик'),
(2, 'Петров Петр Петрович', '2019-08-20', 'Токарь'),
(3, 'Сидоров Сидор Сидорович', '2021-01-10', 'Инженер-технолог'),
(4, 'Кузнецова Мария Олеговна', '2018-03-01', 'Инженер-конструктор'),
(5, 'Васильев Василий Васильевич', '2022-06-01', 'Начальник участка'),
(6, 'Смирнов Алексей Игоревич', '2020-11-11', 'Лаборант'),
(7, 'Орлова Ольга Сергеевна', '2021-07-23', 'Мастер цеха'),
(8, 'Зайцев Захар Захарович', '2019-02-18', 'Сварщик');


-- Заполнение таблицы worker
INSERT INTO "worker" ("employee_id", "category") VALUES
(1, 'Сборщик'),
(2, 'Токарь'),
(8, 'Сварщик');

-- Заполнение таблицы engineer
INSERT INTO "engineer" ("employee_id", "category_id") VALUES
(3, 2), -- Сидоров, Инженер-технолог
(4, 1), -- Кузнецова, Инженер-конструктор
(5, 2), -- Васильев, Начальник участка (будем считать его технологом для примера)
(7, 3); -- Орлова, Мастер цеха (будем считать ее техником для примера)

-- Заполнение таблицы worker_boss (Начальники бригад из числа рабочих)
-- Предположим, Иванов - бригадир
INSERT INTO "worker_boss" ("worker_id") VALUES
(1);

-- Заполнение таблицы work_team
INSERT INTO "worker_boss" ("worker_id") VALUES (2) ON CONFLICT (worker_id) DO NOTHING;
-- Для этого сначала добавим Петрова в worker_boss, если он еще не там
INSERT INTO "work_team" ("id", "name", "worker_boss_id", "area_id") VALUES
(1, 'Бригада сборщиков №1', 1, 1),
(2, 'Бригада токарей', 2, 2); -- Предположим, Петров тоже бригадир для второй бригады (нужно добавить его в worker_boss)


-- Заполнение таблицы work_team_member
INSERT INTO "work_team_member" ("work_team_id", "worker_id", "join_date") VALUES
(1, 1, '2020-05-20'); -- Иванов в своей бригаде

-- Заполнение таблицы work_type
INSERT INTO "work_type" ("id", "work_name", "area_id", "work_team_id") VALUES
(1, 'Сборка основного модуля', 1, 1),
(2, 'Токарная обработка вала', 2, 2),
(3, 'Сварка корпуса', 1, 1); -- Предположим, сварка тоже на участке сборки А бригадой №1

-- Заполнение таблицы item_work_type
INSERT INTO "item_work_type" ("id", "seq_number", "item_id", "work_type_id", "start_date", "end_date") VALUES
(1, 1, 1, 1, '2023-01-11', '2023-01-20'),
(2, 1, 2, 2, '2023-02-16', '2023-02-25'),
(3, 2, 1, 3, '2023-01-21', '2023-01-25'); -- Сварка для Изделия Альфа-001

-- Заполнение таблицы ready_item
INSERT INTO "ready_item" ("id", "item_id", "start_date", "completion_date", "counter") VALUES
(1, 3, '2023-01-01', '2023-04-01', 1);

-- Заполнение таблицы lab_equip
INSERT INTO "lab_equip" ("id", "lab_id", "name") VALUES
(1, 1, 'Осциллограф цифровой'),
(2, 1, 'Микрометр'),
(3, 2, 'Разрывная машина');

-- Заполнение таблицы area_boss
INSERT INTO "area_boss" ("area_id", "engineer_id") VALUES
(1, 5); -- Васильев - начальник участка сборки А

-- Заполнение таблицы hall_bosses
INSERT INTO "hall_bosses" ("hall_id", "engineer_id") VALUES
(1, 4); -- Кузнецова - начальник сборочного цеха №1

-- Заполнение таблицы masters
INSERT INTO "masters" ("id", "area_id", "engineer_id") VALUES
(1, 1, 7); -- Орлова - мастер на участке сборки А

-- Заполнение таблицы lab_worker
INSERT INTO "lab_worker" ("employee_id", "lab_id") VALUES
(6, 1); -- Смирнов работает в Лаборатории контроля качества

-- Заполнение таблицы item_tests
INSERT INTO "item_tests" ("id", "item_id", "lab_worker_id", "lab_equip_id", "test_date", "result") VALUES
(1, 2, 6, 1, '2023-02-28', 'Соответствует ТУ'), -- Тестирование Изделия Бета-002
(2, 3, 6, 2, '2023-03-25', 'Годен'); -- Тестирование Изделия Гамма-003

-- Заполнение таблицы lab_hall
INSERT INTO "lab_hall" ("lab_id", "hall_id") VALUES
(1, 3); -- Лаборатория контроля качества находится в Испытательном цехе

-- Заполнение таблицы employee_movement
INSERT INTO "employee_movement" ("id", "employee_id", "position_from", "position_to", "date") VALUES
(1, 1, 'Стажер', 'Сборщик', '2020-08-15'),
(2, 3, 'Техник', 'Инженер-технолог', '2022-01-10');

INSERT INTO "areas_items" ("area_id", "item_id") VALUES
(1, 1),
(2, 2),
(3, 3),
(1, 4);
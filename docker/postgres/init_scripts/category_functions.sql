-- Функции для работы с категориями

-- Функция для получения всех доступных категорий изделий
CREATE OR REPLACE FUNCTION get_item_categories()
RETURNS item_category_enum[] AS $$
BEGIN
    RETURN ARRAY['civil_aircraft', 'transport_aircraft', 'military_aircraft', 
                 'glider', 'helicopter', 'hang_glider', 
                 'artillery_rocket', 'aviation_rocket', 'naval_rocket', 'other']::item_category_enum[];
END;
$$ LANGUAGE plpgsql;

-- Функция для получения всех доступных категорий инженеров
CREATE OR REPLACE FUNCTION get_engineer_categories()
RETURNS engineer_category_enum[] AS $$
BEGIN
    RETURN ARRAY['engineer', 'technologist', 'technician']::engineer_category_enum[];
END;
$$ LANGUAGE plpgsql;

-- Функция для получения всех доступных категорий рабочих
CREATE OR REPLACE FUNCTION get_worker_categories()
RETURNS worker_category_enum[] AS $$
BEGIN
    RETURN ARRAY['assembler', 'turner', 'locksmith', 'welder']::worker_category_enum[];
END;
$$ LANGUAGE plpgsql;

-- Функция для получения всех доступных статусов изделий
CREATE OR REPLACE FUNCTION get_item_statuses()
RETURNS item_status_enum[] AS $$
BEGIN
    RETURN ARRAY['in_progress', 'testing', 'completed']::item_status_enum[];
END;
$$ LANGUAGE plpgsql;

-- Функция для проверки, является ли категория изделия ракетой
CREATE OR REPLACE FUNCTION is_rocket_category(category item_category_enum)
RETURNS boolean AS $$
BEGIN
    RETURN category IN ('artillery_rocket', 'aviation_rocket', 'naval_rocket');
END;
$$ LANGUAGE plpgsql;

-- Функция для проверки, является ли категория изделия самолетом
CREATE OR REPLACE FUNCTION is_aircraft_category(category item_category_enum)
RETURNS boolean AS $$
BEGIN
    RETURN category IN ('civil_aircraft', 'transport_aircraft', 'military_aircraft');
END;
$$ LANGUAGE plpgsql;

-- Функция для получения русского названия категории изделия
CREATE OR REPLACE FUNCTION get_item_category_name_ru(category item_category_enum)
RETURNS varchar(255) AS $$
BEGIN
    CASE category
        WHEN 'civil_aircraft' THEN RETURN 'Гражданские самолеты';
        WHEN 'transport_aircraft' THEN RETURN 'Транспортные самолеты';
        WHEN 'military_aircraft' THEN RETURN 'Военные самолеты';
        WHEN 'glider' THEN RETURN 'Планеры';
        WHEN 'helicopter' THEN RETURN 'Вертолеты';
        WHEN 'hang_glider' THEN RETURN 'Дельтопланы';
        WHEN 'artillery_rocket' THEN RETURN 'Артиллерийские ракеты';
        WHEN 'aviation_rocket' THEN RETURN 'Авиационные ракеты';
        WHEN 'naval_rocket' THEN RETURN 'Военно-морские ракеты';
        WHEN 'other' THEN RETURN 'Прочие изделия';
        ELSE RETURN 'Неизвестная категория';
    END CASE;
END;
$$ LANGUAGE plpgsql;

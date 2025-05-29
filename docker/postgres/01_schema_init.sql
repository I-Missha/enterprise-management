
-- Создание ENUM типов для категорий
CREATE TYPE item_category_enum AS ENUM (
    'civil_aircraft',      -- гражданские самолеты
    'transport_aircraft',  -- транспортные самолеты
    'military_aircraft',   -- военные самолеты
    'glider',             -- планеры
    'helicopter',         -- вертолеты
    'hang_glider',        -- дельтопланы
    'artillery_rocket',   -- артиллерийские ракеты
    'aviation_rocket',    -- авиационные ракеты
    'naval_rocket',       -- военно-морские ракеты
    'other'               -- прочие изделия
);

CREATE TYPE engineer_category_enum AS ENUM (
    'engineer',    -- инженеры
    'technologist', -- технологи
    'technician'   -- техники
);

CREATE TYPE worker_category_enum AS ENUM (
    'assembler',   -- сборщики
    'turner',      -- токари
    'locksmith',   -- слесари
    'welder'       -- сварщики
);

CREATE TYPE item_status_enum AS ENUM (
    'in_progress', -- в процессе изготовления
    'testing',     -- на испытаниях
    'completed'    -- завершено
);

CREATE TABLE "testing_laboratory" (
    "id" SERIAL PRIMARY KEY,
    "name" varchar(255) NOT NULL
);

CREATE TABLE "production_hall" (
   "id" SERIAL PRIMARY KEY,
   "name" varchar(255) NOT NULL,
   "shop_manager_id" integer -- начальник цеха
);

CREATE TABLE "production_area" (
    "id" SERIAL PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "hall_id" integer NOT NULL,
    "area_manager_id" integer -- начальник участка
);

CREATE TABLE "category_item" (
    "id" SERIAL PRIMARY KEY,
    "name" item_category_enum NOT NULL,
    "attribute" varchar(255)
);

CREATE TABLE "category_engineer" (
    "id" SERIAL PRIMARY KEY,
    "name" engineer_category_enum NOT NULL,
    "attribute" varchar(255)
);

CREATE TABLE "type_item" (
    "id" SERIAL PRIMARY KEY,
    "name" varchar(255) NOT NULL UNIQUE,
    "category_id" integer NOT NULL
);

CREATE TABLE "item" (
    "id" SERIAL PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "type_id" integer NOT NULL,
    "hall_id" integer NOT NULL,
    "status" item_status_enum NOT NULL
);

CREATE TABLE "item_work_type" (
    "id" SERIAL PRIMARY KEY,
    "seq_number" integer NOT NULL,
    "item_id" integer NOT NULL,
    "work_type_id" integer NOT NULL,
    "start_date" date,
    "end_date" date,
    UNIQUE ("item_id", "work_type_id")
);

CREATE TABLE "completed_item" (
    "id" SERIAL PRIMARY KEY,
    "item_id" integer NOT NULL,
    "production_start_date" date NOT NULL,
    "production_completion_date" date, -- NULL for items still in progress
    "assembled_by_team_id" integer NOT NULL,
    "assembled_in_hall_id" integer NOT NULL,
    "final_area_id" integer NOT NULL,
    "quantity_produced" integer NOT NULL DEFAULT 1,
    "status" varchar(50) NOT NULL DEFAULT 'completed', -- completed, testing, shipped
    "notes" text,
    "created_at" date NOT NULL,
    "updated_at" date NOT NULL
);

-- Таблица для испытаний завершённых изделий
CREATE TABLE "completed_item_test" (
    "id" SERIAL PRIMARY KEY,
    "completed_item_id" integer NOT NULL,
    "lab_id" integer NOT NULL,
    "test_start_date" date NOT NULL,
    "test_completion_date" date,
    "test_result" varchar(255),
    "test_status" varchar(50) NOT NULL DEFAULT 'in_progress', -- in_progress, passed, failed
    "conducted_by_worker_id" integer NOT NULL,
    "notes" text,
    "created_at" date NOT NULL,
    "updated_at" date NOT NULL,
    UNIQUE ("completed_item_id", "lab_id", "test_start_date")
);

-- Таблица для связи испытаний с оборудованием
CREATE TABLE "test_equipment_usage" (
    "id" SERIAL PRIMARY KEY,
    "completed_item_test_id" integer NOT NULL,
    "lab_equip_id" integer NOT NULL,
    "usage_date" date NOT NULL,
    "duration_hours" decimal(5,2),
    "notes" text,
    "created_at" date NOT NULL
);

CREATE TABLE "lab_equip" (
    "id" SERIAL PRIMARY KEY,
    "lab_id" integer NOT NULL,
    "name" varchar(255) NOT NULL
);


CREATE TABLE "employee" (
    "id" SERIAL PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "hire_date" date NOT NULL,
    "current_position" varchar(255)
);

CREATE TABLE "worker" (
    "employee_id" integer PRIMARY KEY,
    "hall_id" integer NOT NULL,
    "area_id" integer NOT NULL,
    "work_team_id" integer NOT NULL,
    "category" worker_category_enum NOT NULL
);

CREATE TABLE "engineer" (
    "employee_id" integer PRIMARY KEY,
    "hall_id" integer NOT NULL,
    "area_id" integer NOT NULL,
    "category_id" integer NOT NULL
);

CREATE TABLE "worker_boss" (
    "worker_id" integer PRIMARY KEY
);

CREATE TABLE "work_team" (
    "id" SERIAL PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "area_id" integer NOT NULL,
    "hall_id" integer NOT NULL,
    "team_leader_id" integer -- бригадир
);


CREATE TABLE "work_type" (
    "id" SERIAL PRIMARY KEY,
    "work_name" varchar(255) NOT NULL,
    "area_id" integer NOT NULL,
    "work_team_id" integer NOT NULL
);

CREATE TABLE "area_boss" (
    "area_id" integer PRIMARY KEY,
    "engineer_id" integer NOT NULL
);

CREATE TABLE "hall_bosses" (
    "hall_id" integer PRIMARY KEY,
    "engineer_id" integer NOT NULL
);

CREATE TABLE "masters" (
    "id" SERIAL PRIMARY KEY,
    "area_id" integer NOT NULL,
    "engineer_id" integer NOT NULL,
    UNIQUE ("area_id", "engineer_id")
);

CREATE TABLE "lab_worker" (
    "employee_id" integer PRIMARY KEY,
    "lab_id" integer NOT NULL
);

CREATE TABLE "item_tests" (
    "id" SERIAL PRIMARY KEY,
    "item_id" integer NOT NULL,
    "lab_worker_id" integer NOT NULL,
    "lab_equip_id" integer NOT NULL,
    "test_date" date NOT NULL,
    "result" varchar(255),
    UNIQUE ("item_id", "lab_worker_id", "lab_equip_id", "test_date")
);

CREATE TABLE "lab_hall" (
   "lab_id" integer NOT NULL,
   "hall_id" integer NOT NULL,
   PRIMARY KEY ("hall_id", "lab_id")
);

CREATE TABLE "areas_items" (
    "area_id" integer NOT NULL,
    "item_id" integer NOT NULL,
    PRIMARY KEY ("area_id", "item_id")
);

-- Внешние ключи для новых таблиц завершённых изделий и испытаний
ALTER TABLE "completed_item" ADD CONSTRAINT "fk_completed_item"
    FOREIGN KEY ("item_id") REFERENCES "item" ("id");

ALTER TABLE "completed_item" ADD CONSTRAINT "fk_completed_item_team"
    FOREIGN KEY ("assembled_by_team_id") REFERENCES "work_team" ("id");

ALTER TABLE "completed_item" ADD CONSTRAINT "fk_completed_item_hall"
    FOREIGN KEY ("assembled_in_hall_id") REFERENCES "production_hall" ("id");

ALTER TABLE "completed_item" ADD CONSTRAINT "fk_completed_item_area"
    FOREIGN KEY ("final_area_id") REFERENCES "production_area" ("id");

ALTER TABLE "completed_item_test" ADD CONSTRAINT "fk_completed_item_test_completed"
    FOREIGN KEY ("completed_item_id") REFERENCES "completed_item" ("id");

ALTER TABLE "completed_item_test" ADD CONSTRAINT "fk_completed_item_test_lab"
    FOREIGN KEY ("lab_id") REFERENCES "testing_laboratory" ("id");

ALTER TABLE "completed_item_test" ADD CONSTRAINT "fk_completed_item_test_worker"
    FOREIGN KEY ("conducted_by_worker_id") REFERENCES "lab_worker" ("employee_id");

ALTER TABLE "test_equipment_usage" ADD CONSTRAINT "fk_test_equip_usage_test"
    FOREIGN KEY ("completed_item_test_id") REFERENCES "completed_item_test" ("id");

ALTER TABLE "test_equipment_usage" ADD CONSTRAINT "fk_test_equip_usage_equip"
    FOREIGN KEY ("lab_equip_id") REFERENCES "lab_equip" ("id");

ALTER TABLE "worker" ADD CONSTRAINT "fk_worker_hall"
    FOREIGN KEY ("hall_id") REFERENCES "production_hall" ("id");

ALTER TABLE "worker" ADD CONSTRAINT "fk_worker_area"
    FOREIGN KEY ("area_id") REFERENCES "production_area" ("id");

ALTER TABLE "worker" ADD CONSTRAINT "fk_worker_team"
    FOREIGN KEY ("work_team_id") REFERENCES "work_team" ("id");

ALTER TABLE "engineer" ADD CONSTRAINT "fk_engineer_hall"
    FOREIGN KEY ("hall_id") REFERENCES "production_hall" ("id");

ALTER TABLE "engineer" ADD CONSTRAINT "fk_engineer_area"
    FOREIGN KEY ("area_id") REFERENCES "production_area" ("id");

ALTER TABLE "areas_items" ADD CONSTRAINT "fk_area_item"
    FOREIGN KEY ("area_id") REFERENCES "production_area" ("id");

ALTER TABLE "areas_items" ADD CONSTRAINT "fk_item_area"
    FOREIGN KEY ("item_id") REFERENCES "item" ("id");

ALTER TABLE "production_area" ADD CONSTRAINT "fk_area_hall"
    FOREIGN KEY ("hall_id") REFERENCES "production_hall" ("id");

ALTER TABLE "type_item" ADD CONSTRAINT "fk_type_category"
    FOREIGN KEY ("category_id") REFERENCES "category_item" ("id");

ALTER TABLE "item" ADD CONSTRAINT "fk_item_type"
    FOREIGN KEY ("type_id") REFERENCES "type_item" ("id");

ALTER TABLE "item" ADD CONSTRAINT "fk_item_hall"
    FOREIGN KEY ("hall_id") REFERENCES "production_hall" ("id");

ALTER TABLE "item_work_type" ADD CONSTRAINT "fk_iwt_item"
    FOREIGN KEY ("item_id") REFERENCES "item" ("id");

ALTER TABLE "item_work_type" ADD CONSTRAINT "fk_iwt_work_type"
    FOREIGN KEY ("work_type_id") REFERENCES "work_type" ("id");



ALTER TABLE "lab_equip" ADD CONSTRAINT "fk_equip_lab"
    FOREIGN KEY ("lab_id") REFERENCES "testing_laboratory" ("id");

ALTER TABLE "worker" ADD CONSTRAINT "fk_worker_employee"
    FOREIGN KEY ("employee_id") REFERENCES "employee" ("id");

ALTER TABLE "engineer" ADD CONSTRAINT "fk_engineer_employee"
    FOREIGN KEY ("employee_id") REFERENCES "employee" ("id");

ALTER TABLE "engineer" ADD CONSTRAINT "fk_engineer_category"
    FOREIGN KEY ("category_id") REFERENCES "category_engineer" ("id");

ALTER TABLE "worker_boss" ADD CONSTRAINT "fk_boss_worker"
    FOREIGN KEY ("worker_id") REFERENCES "worker" ("employee_id");

-- ALTER TABLE "work_team" ADD CONSTRAINT "fk_team_boss"
--     FOREIGN KEY ("worker_boss_id") REFERENCES "worker_boss" ("worker_id");

ALTER TABLE "work_team" ADD CONSTRAINT "fk_team_area"
    FOREIGN KEY ("area_id") REFERENCES "production_area" ("id");

-- ALTER TABLE "work_team_member" ADD CONSTRAINT "fk_member_team"
--     FOREIGN KEY ("work_team_id") REFERENCES "work_team" ("id");
--
-- ALTER TABLE "work_team_member" ADD CONSTRAINT "fk_member_worker"
--     FOREIGN KEY ("worker_id") REFERENCES "worker" ("employee_id");

ALTER TABLE "work_type" ADD CONSTRAINT "fk_work_area"
    FOREIGN KEY ("area_id") REFERENCES "production_area" ("id");

ALTER TABLE "work_type" ADD CONSTRAINT "fk_work_team"
    FOREIGN KEY ("work_team_id") REFERENCES "work_team" ("id");

ALTER TABLE "area_boss" ADD CONSTRAINT "fk_area_boss_area"
    FOREIGN KEY ("area_id") REFERENCES "production_area" ("id");

ALTER TABLE "area_boss" ADD CONSTRAINT "fk_area_boss_engineer"
    FOREIGN KEY ("engineer_id") REFERENCES "engineer" ("employee_id");

ALTER TABLE "hall_bosses" ADD CONSTRAINT "fk_hall_boss_hall"
    FOREIGN KEY ("hall_id") REFERENCES "production_hall" ("id");

ALTER TABLE "hall_bosses" ADD CONSTRAINT "fk_hall_boss_engineer"
    FOREIGN KEY ("engineer_id") REFERENCES "engineer" ("employee_id");

ALTER TABLE "masters" ADD CONSTRAINT "fk_master_area"
    FOREIGN KEY ("area_id") REFERENCES "production_area" ("id");

ALTER TABLE "masters" ADD CONSTRAINT "fk_master_engineer"
    FOREIGN KEY ("engineer_id") REFERENCES "engineer" ("employee_id");

ALTER TABLE "lab_worker" ADD CONSTRAINT "fk_lab_worker_employee"
    FOREIGN KEY ("employee_id") REFERENCES "employee" ("id");

ALTER TABLE "lab_worker" ADD CONSTRAINT "fk_lab_worker_lab"
    FOREIGN KEY ("lab_id") REFERENCES "testing_laboratory" ("id");

ALTER TABLE "item_tests" ADD CONSTRAINT "fk_test_item"
    FOREIGN KEY ("item_id") REFERENCES "item" ("id");

ALTER TABLE "item_tests" ADD CONSTRAINT "fk_test_worker"
    FOREIGN KEY ("lab_worker_id") REFERENCES "lab_worker" ("employee_id");

ALTER TABLE "item_tests" ADD CONSTRAINT "fk_test_equip"
    FOREIGN KEY ("lab_equip_id") REFERENCES "lab_equip" ("id");

ALTER TABLE "lab_hall" ADD CONSTRAINT "fk_lab_hall_lab"
    FOREIGN KEY ("lab_id") REFERENCES "testing_laboratory" ("id");

ALTER TABLE "lab_hall" ADD CONSTRAINT "fk_lab_hall_hall"
    FOREIGN KEY ("hall_id") REFERENCES "production_hall" ("id");
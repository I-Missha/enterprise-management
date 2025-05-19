CREATE TABLE "testing_laboratory" (
    "id" integer PRIMARY KEY,
    "name" varchar(255) NOT NULL
);

CREATE TABLE "production_halls" (
   "id" integer PRIMARY KEY,
   "name" varchar(255) NOT NULL
);

CREATE TABLE "production_area" (
    "id" integer PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "hall_id" integer NOT NULL
);

CREATE TABLE "category_item" (
    "id" integer PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "attribute" varchar(255)
);

CREATE TABLE "category_engineer" (
    "id" integer PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "attribute" varchar(255)
);

CREATE TABLE "type_item" (
    "id" integer PRIMARY KEY,
    "name" varchar(255) NOT NULL UNIQUE,
    "category_id" integer NOT NULL
);

CREATE TABLE "item" (
    "id" integer PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "type_id" integer NOT NULL,
    "hall_id" integer NOT NULL,
    "status" varchar(50) CHECK (status IN ('in_progress', 'testing', 'completed'))
);

CREATE TABLE "item_work_type" (
    "id" integer PRIMARY KEY,
    "seq_number" integer NOT NULL,
    "item_id" integer NOT NULL,
    "work_type_id" integer NOT NULL,
    "start_date" date,
    "end_date" date,
    UNIQUE ("item_id", "work_type_id")
);

CREATE TABLE "ready_item" (
    "id" integer PRIMARY KEY,
    "item_id" integer,
    "start_date" date NOT NULL,
    "completion_date" date NOT NULL,
    "counter" integer NOT NULL DEFAULT 1
);

CREATE TABLE "lab_equip" (
    "id" integer PRIMARY KEY,
    "lab_id" integer NOT NULL,
    "name" varchar(255) NOT NULL
);

CREATE TABLE "employee" (
    "id" integer PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "hire_date" date NOT NULL,
    "current_position" varchar(255)
);

CREATE TABLE "worker" (
    "employee_id" integer PRIMARY KEY,
    "category" varchar(255) NOT NULL
);

CREATE TABLE "engineer" (
    "employee_id" integer PRIMARY KEY,
    "category_id" integer NOT NULL
);

CREATE TABLE "worker_boss" (
    "worker_id" integer PRIMARY KEY
);

CREATE TABLE "work_team" (
    "id" integer PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "worker_boss_id" integer NOT NULL,
    "area_id" integer NOT NULL
);

-- Таблица для связи рабочих с бригадами
CREATE TABLE "work_team_member" (
    "work_team_id" integer NOT NULL,
    "worker_id" integer NOT NULL,
    "join_date" date NOT NULL,
    PRIMARY KEY ("work_team_id", "worker_id")
);

CREATE TABLE "work_type" (
    "id" integer PRIMARY KEY,
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
    "id" integer PRIMARY KEY,
    "area_id" integer NOT NULL,
    "engineer_id" integer NOT NULL,
    UNIQUE ("area_id", "engineer_id")
);

CREATE TABLE "lab_worker" (
    "employee_id" integer PRIMARY KEY,
    "lab_id" integer NOT NULL
);

CREATE TABLE "item_tests" (
    "id" integer PRIMARY KEY,
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

CREATE TABLE "employee_movement" (
    "id" integer PRIMARY KEY,
    "employee_id" integer NOT NULL,
    "position_from" varchar(255),
    "position_to" varchar(255),
    "date" date NOT NULL
);

CREATE TABLE "areas_items" (
    "area_id" integer NOT NULL,
    "item_id" integer NOT NULL,
    PRIMARY KEY ("area_id", "item_id")
);

ALTER TABLE "areas_items" ADD CONSTRAINT "fk_area_item"
    FOREIGN KEY ("area_id") REFERENCES "production_area" ("id");

ALTER TABLE "areas_items" ADD CONSTRAINT "fk_item_area"
    FOREIGN KEY ("item_id") REFERENCES "item" ("id");

ALTER TABLE "production_area" ADD CONSTRAINT "fk_area_hall"
    FOREIGN KEY ("hall_id") REFERENCES "production_halls" ("id");

ALTER TABLE "type_item" ADD CONSTRAINT "fk_type_category"
    FOREIGN KEY ("category_id") REFERENCES "category_item" ("id");

ALTER TABLE "item" ADD CONSTRAINT "fk_item_type"
    FOREIGN KEY ("type_id") REFERENCES "type_item" ("id");

ALTER TABLE "item" ADD CONSTRAINT "fk_item_hall"
    FOREIGN KEY ("hall_id") REFERENCES "production_halls" ("id");

ALTER TABLE "item_work_type" ADD CONSTRAINT "fk_iwt_item"
    FOREIGN KEY ("item_id") REFERENCES "item" ("id");

ALTER TABLE "item_work_type" ADD CONSTRAINT "fk_iwt_work_type"
    FOREIGN KEY ("work_type_id") REFERENCES "work_type" ("id");

ALTER TABLE "ready_item" ADD CONSTRAINT "fk_ready_item"
    FOREIGN KEY ("item_id") REFERENCES "item" ("id");

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

ALTER TABLE "work_team" ADD CONSTRAINT "fk_team_boss"
    FOREIGN KEY ("worker_boss_id") REFERENCES "worker_boss" ("worker_id");

ALTER TABLE "work_team" ADD CONSTRAINT "fk_team_area"
    FOREIGN KEY ("area_id") REFERENCES "production_area" ("id");

ALTER TABLE "work_team_member" ADD CONSTRAINT "fk_member_team"
    FOREIGN KEY ("work_team_id") REFERENCES "work_team" ("id");

ALTER TABLE "work_team_member" ADD CONSTRAINT "fk_member_worker"
    FOREIGN KEY ("worker_id") REFERENCES "worker" ("employee_id");

ALTER TABLE "work_type" ADD CONSTRAINT "fk_work_area"
    FOREIGN KEY ("area_id") REFERENCES "production_area" ("id");

ALTER TABLE "work_type" ADD CONSTRAINT "fk_work_team"
    FOREIGN KEY ("work_team_id") REFERENCES "work_team" ("id");

ALTER TABLE "area_boss" ADD CONSTRAINT "fk_area_boss_area"
    FOREIGN KEY ("area_id") REFERENCES "production_area" ("id");

ALTER TABLE "area_boss" ADD CONSTRAINT "fk_area_boss_engineer"
    FOREIGN KEY ("engineer_id") REFERENCES "engineer" ("employee_id");

ALTER TABLE "hall_bosses" ADD CONSTRAINT "fk_hall_boss_hall"
    FOREIGN KEY ("hall_id") REFERENCES "production_halls" ("id");

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
    FOREIGN KEY ("hall_id") REFERENCES "production_halls" ("id");

ALTER TABLE "employee_movement" ADD CONSTRAINT "fk_movement_employee"
    FOREIGN KEY ("employee_id") REFERENCES "employee" ("id");
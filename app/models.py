from app import db
from sqlalchemy import Enum
import enum

# Определение ENUM типов
class ItemCategoryEnum(enum.Enum):
    civil_aircraft = 'civil_aircraft'
    transport_aircraft = 'transport_aircraft'
    military_aircraft = 'military_aircraft'
    glider = 'glider'
    helicopter = 'helicopter'
    hang_glider = 'hang_glider'
    artillery_rocket = 'artillery_rocket'
    aviation_rocket = 'aviation_rocket'
    naval_rocket = 'naval_rocket'
    other = 'other'

class EngineerCategoryEnum(enum.Enum):
    engineer = 'engineer'
    technologist = 'technologist'
    technician = 'technician'

class WorkerCategoryEnum(enum.Enum):
    assembler = 'assembler'
    turner = 'turner'
    locksmith = 'locksmith'
    welder = 'welder'

class ItemStatusEnum(enum.Enum):
    in_progress = 'in_progress'
    testing = 'testing'
    completed = 'completed'

# Модели таблиц
class TestingLaboratory(db.Model):
    __tablename__ = 'testing_laboratory'
    
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(255), nullable=False)
    
    # Связи
    lab_equipment = db.relationship('LabEquip', backref='laboratory', lazy=True)
    lab_workers = db.relationship('LabWorker', backref='laboratory', lazy=True)

class ProductionHall(db.Model):
    __tablename__ = 'production_hall'
    
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(255), nullable=False)
    shop_manager_id = db.Column(db.Integer, db.ForeignKey('engineer.employee_id'), nullable=True)
    
    # Связи
    areas = db.relationship('ProductionArea', backref='hall', lazy=True)
    workers = db.relationship('Worker', backref='hall', lazy=True)
    engineers = db.relationship('Engineer', backref='hall', lazy=True)
    shop_manager = db.relationship('Engineer', foreign_keys=[shop_manager_id], backref='managed_halls', lazy=True)

class ProductionArea(db.Model):
    __tablename__ = 'production_area'
    
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(255), nullable=False)
    hall_id = db.Column(db.Integer, db.ForeignKey('production_hall.id'), nullable=False)
    area_manager_id = db.Column(db.Integer, db.ForeignKey('engineer.employee_id'), nullable=True)
    
    # Связи
    workers = db.relationship('Worker', backref='area', lazy=True)
    engineers = db.relationship('Engineer', backref='area', lazy=True)
    work_teams = db.relationship('WorkTeam', backref='area', lazy=True)
    area_manager = db.relationship('Engineer', foreign_keys=[area_manager_id], backref='managed_areas', lazy=True)

class CategoryItem(db.Model):
    __tablename__ = 'category_item'
    
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(Enum(ItemCategoryEnum), nullable=False)
    attribute = db.Column(db.String(255))
    
    # Связи
    types = db.relationship('TypeItem', backref='category', lazy=True)

class CategoryEngineer(db.Model):
    __tablename__ = 'category_engineer'
    
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(Enum(EngineerCategoryEnum), nullable=False)
    attribute = db.Column(db.String(255))
    
    # Связи
    engineers = db.relationship('Engineer', backref='category', lazy=True)

class TypeItem(db.Model):
    __tablename__ = 'type_item'
    
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(255), nullable=False, unique=True)
    category_id = db.Column(db.Integer, db.ForeignKey('category_item.id'), nullable=False)
    
    # Связи
    items = db.relationship('Item', backref='type', lazy=True)

class Item(db.Model):
    __tablename__ = 'item'
    
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(255), nullable=False)
    type_id = db.Column(db.Integer, db.ForeignKey('type_item.id'), nullable=False)
    hall_id = db.Column(db.Integer, db.ForeignKey('production_hall.id'), nullable=False)
    status = db.Column(Enum(ItemStatusEnum), nullable=False)
    
    # Связи
    work_types = db.relationship('ItemWorkType', backref='item', lazy=True)
    completed_items = db.relationship('CompletedItem', backref='item', lazy=True)

class Employee(db.Model):
    __tablename__ = 'employee'
    
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(255), nullable=False)
    hire_date = db.Column(db.Date, nullable=False)
    current_position = db.Column(db.String(255))
    
    # Связи
    worker = db.relationship('Worker', backref='employee', uselist=False)
    engineer = db.relationship('Engineer', backref='employee', uselist=False)
    lab_worker = db.relationship('LabWorker', backref='employee', uselist=False)

class Worker(db.Model):
    __tablename__ = 'worker'
    
    employee_id = db.Column(db.Integer, db.ForeignKey('employee.id'), primary_key=True)
    hall_id = db.Column(db.Integer, db.ForeignKey('production_hall.id'), nullable=False)
    area_id = db.Column(db.Integer, db.ForeignKey('production_area.id'), nullable=False)
    work_team_id = db.Column(db.Integer, db.ForeignKey('work_team.id'), nullable=False)
    category = db.Column(Enum(WorkerCategoryEnum), nullable=False)

class Engineer(db.Model):
    __tablename__ = 'engineer'
    
    employee_id = db.Column(db.Integer, db.ForeignKey('employee.id'), primary_key=True)
    hall_id = db.Column(db.Integer, db.ForeignKey('production_hall.id'), nullable=False)
    area_id = db.Column(db.Integer, db.ForeignKey('production_area.id'), nullable=False)
    category_id = db.Column(db.Integer, db.ForeignKey('category_engineer.id'), nullable=False)

class WorkTeam(db.Model):
    __tablename__ = 'work_team'
    
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(255), nullable=False)
    area_id = db.Column(db.Integer, db.ForeignKey('production_area.id'), nullable=False)
    hall_id = db.Column(db.Integer, db.ForeignKey('production_hall.id'), nullable=False)
    team_leader_id = db.Column(db.Integer, db.ForeignKey('worker.employee_id'), nullable=True)
    
    # Связи
    workers = db.relationship('Worker', backref='work_team', lazy=True)
    team_leader = db.relationship('Worker', foreign_keys=[team_leader_id], backref='led_teams', lazy=True)

class WorkType(db.Model):
    __tablename__ = 'work_type'
    
    id = db.Column(db.Integer, primary_key=True)
    work_name = db.Column(db.String(255), nullable=False)
    area_id = db.Column(db.Integer, db.ForeignKey('production_area.id'), nullable=False)
    work_team_id = db.Column(db.Integer, db.ForeignKey('work_team.id'), nullable=False)

class ItemWorkType(db.Model):
    __tablename__ = 'item_work_type'
    
    id = db.Column(db.Integer, primary_key=True)
    seq_number = db.Column(db.Integer, nullable=False)
    item_id = db.Column(db.Integer, db.ForeignKey('item.id'), nullable=False)
    work_type_id = db.Column(db.Integer, db.ForeignKey('work_type.id'), nullable=False)
    start_date = db.Column(db.Date)
    end_date = db.Column(db.Date)
    
    __table_args__ = (db.UniqueConstraint('item_id', 'work_type_id'),)

class CompletedItem(db.Model):
    __tablename__ = 'completed_item'
    
    id = db.Column(db.Integer, primary_key=True)
    item_id = db.Column(db.Integer, db.ForeignKey('item.id'), nullable=False)
    production_start_date = db.Column(db.Date, nullable=False)
    production_completion_date = db.Column(db.Date, nullable=False)
    assembled_by_team_id = db.Column(db.Integer, db.ForeignKey('work_team.id'), nullable=False)
    assembled_in_hall_id = db.Column(db.Integer, db.ForeignKey('production_hall.id'), nullable=False)
    final_area_id = db.Column(db.Integer, db.ForeignKey('production_area.id'), nullable=False)
    quantity_produced = db.Column(db.Integer, nullable=False, default=1)
    status = db.Column(db.String(50), nullable=False, default='completed')
    notes = db.Column(db.Text)
    created_at = db.Column(db.Date, nullable=False)
    updated_at = db.Column(db.Date, nullable=False)

class LabEquip(db.Model):
    __tablename__ = 'lab_equip'
    
    id = db.Column(db.Integer, primary_key=True)
    lab_id = db.Column(db.Integer, db.ForeignKey('testing_laboratory.id'), nullable=False)
    name = db.Column(db.String(255), nullable=False)

class LabWorker(db.Model):
    __tablename__ = 'lab_worker'
    
    employee_id = db.Column(db.Integer, db.ForeignKey('employee.id'), primary_key=True)
    lab_id = db.Column(db.Integer, db.ForeignKey('testing_laboratory.id'), nullable=False)

class ItemTests(db.Model):
    __tablename__ = 'item_tests'
    
    id = db.Column(db.Integer, primary_key=True)
    item_id = db.Column(db.Integer, db.ForeignKey('item.id'), nullable=False)
    lab_worker_id = db.Column(db.Integer, db.ForeignKey('lab_worker.employee_id'), nullable=False)
    lab_equip_id = db.Column(db.Integer, db.ForeignKey('lab_equip.id'), nullable=False)
    test_date = db.Column(db.Date, nullable=False)
    result = db.Column(db.String(255))
    
    __table_args__ = (db.UniqueConstraint('item_id', 'lab_worker_id', 'lab_equip_id', 'test_date'),)

class CompletedItemTest(db.Model):
    __tablename__ = 'completed_item_test'
    
    id = db.Column(db.Integer, primary_key=True)
    completed_item_id = db.Column(db.Integer, db.ForeignKey('completed_item.id'), nullable=False)
    lab_id = db.Column(db.Integer, db.ForeignKey('testing_laboratory.id'), nullable=False)
    test_start_date = db.Column(db.Date, nullable=False)
    test_completion_date = db.Column(db.Date)
    test_result = db.Column(db.String(255))
    test_status = db.Column(db.String(50), nullable=False, default='in_progress')  # in_progress, passed, failed
    conducted_by_worker_id = db.Column(db.Integer, db.ForeignKey('lab_worker.employee_id'), nullable=False)
    notes = db.Column(db.Text)
    created_at = db.Column(db.Date, nullable=False)
    updated_at = db.Column(db.Date, nullable=False)
    
    __table_args__ = (db.UniqueConstraint('completed_item_id', 'lab_id', 'test_start_date'),)

class AreaBoss(db.Model):
    __tablename__ = 'area_boss'
    
    area_id = db.Column(db.Integer, db.ForeignKey('production_area.id'), primary_key=True)
    engineer_id = db.Column(db.Integer, db.ForeignKey('engineer.employee_id'), nullable=False)

class HallBosses(db.Model):
    __tablename__ = 'hall_bosses'
    
    hall_id = db.Column(db.Integer, db.ForeignKey('production_hall.id'), primary_key=True)
    engineer_id = db.Column(db.Integer, db.ForeignKey('engineer.employee_id'), nullable=False)

class Masters(db.Model):
    __tablename__ = 'masters'
    
    id = db.Column(db.Integer, primary_key=True)
    area_id = db.Column(db.Integer, db.ForeignKey('production_area.id'), nullable=False)
    engineer_id = db.Column(db.Integer, db.ForeignKey('engineer.employee_id'), nullable=False)
    
    __table_args__ = (db.UniqueConstraint('area_id', 'engineer_id'),)

class WorkerBoss(db.Model):
    __tablename__ = 'worker_boss'
    
    worker_id = db.Column(db.Integer, db.ForeignKey('worker.employee_id'), primary_key=True)

class AreasItems(db.Model):
    __tablename__ = 'areas_items'
    
    area_id = db.Column(db.Integer, db.ForeignKey('production_area.id'), primary_key=True)
    item_id = db.Column(db.Integer, db.ForeignKey('item.id'), primary_key=True)

class LabHall(db.Model):
    __tablename__ = 'lab_hall'
    
    lab_id = db.Column(db.Integer, db.ForeignKey('testing_laboratory.id'), primary_key=True)
    hall_id = db.Column(db.Integer, db.ForeignKey('production_hall.id'), primary_key=True)

class TestEquipmentUsage(db.Model):
    __tablename__ = 'test_equipment_usage'
    
    id = db.Column(db.Integer, primary_key=True)
    completed_item_test_id = db.Column(db.Integer, db.ForeignKey('completed_item_test.id'), nullable=False)
    lab_equip_id = db.Column(db.Integer, db.ForeignKey('lab_equip.id'), nullable=False)
    usage_date = db.Column(db.Date, nullable=False)
    duration_hours = db.Column(db.Numeric(5, 2))
    notes = db.Column(db.Text)
    created_at = db.Column(db.Date, nullable=False)
    
    # Relationships
    completed_item_test = db.relationship('CompletedItemTest', backref='equipment_usage')
    lab_equipment = db.relationship('LabEquip', backref='usage_records')

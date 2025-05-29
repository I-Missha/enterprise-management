from flask import Blueprint, render_template, request, redirect, url_for, flash, jsonify
from app import db
from app.models import *
from datetime import datetime
from sqlalchemy import func, and_, or_, distinct
from datetime import datetime, date

main = Blueprint('main', __name__)

@main.route('/')
def index():
    """Главная страница с навигацией"""
    return render_template('index.html')

# Маршруты для производственных цехов
@main.route('/halls')
def halls():
    halls = ProductionHall.query.all()
    return render_template('halls/list.html', halls=halls)

@main.route('/halls/add', methods=['GET', 'POST'])
def add_hall():
    if request.method == 'POST':
        name = request.form['name']
        hall = ProductionHall(name=name)
        db.session.add(hall)
        db.session.commit()
        flash('Цех добавлен успешно!', 'success')
        return redirect(url_for('main.halls'))
    return render_template('halls/add.html')

@main.route('/halls/edit/<int:id>', methods=['GET', 'POST'])
def edit_hall(id):
    hall = ProductionHall.query.get_or_404(id)
    if request.method == 'POST':
        hall.name = request.form['name']
        db.session.commit()
        flash('Цех обновлен успешно!', 'success')
        return redirect(url_for('main.halls'))
    return render_template('halls/edit.html', hall=hall)

@main.route('/halls/delete/<int:id>')
def delete_hall(id):
    hall = ProductionHall.query.get_or_404(id)
    db.session.delete(hall)
    db.session.commit()
    flash('Цех удален!', 'success')
    return redirect(url_for('main.halls'))

# Маршруты для производственных участков
@main.route('/areas')
def areas():
    areas = ProductionArea.query.all()
    return render_template('areas/list.html', areas=areas)

@main.route('/areas/add', methods=['GET', 'POST'])
def add_area():
    if request.method == 'POST':
        name = request.form['name']
        hall_id = request.form['hall_id']
        area = ProductionArea(name=name, hall_id=hall_id)
        db.session.add(area)
        db.session.commit()
        flash('Участок добавлен успешно!', 'success')
        return redirect(url_for('main.areas'))
    halls = ProductionHall.query.all()
    return render_template('areas/add.html', halls=halls)

@main.route('/areas/edit/<int:id>', methods=['GET', 'POST'])
def edit_area(id):
    area = ProductionArea.query.get_or_404(id)
    if request.method == 'POST':
        area.name = request.form['name']
        area.hall_id = request.form['hall_id']
        db.session.commit()
        flash('Участок обновлен успешно!', 'success')
        return redirect(url_for('main.areas'))
    halls = ProductionHall.query.all()
    return render_template('areas/edit.html', area=area, halls=halls)

@main.route('/areas/delete/<int:id>')
def delete_area(id):
    area = ProductionArea.query.get_or_404(id)
    db.session.delete(area)
    db.session.commit()
    flash('Участок удален!', 'success')
    return redirect(url_for('main.areas'))

# Маршруты для сотрудников
@main.route('/employees')
def employees():
    employees = Employee.query.all()
    return render_template('employees/list.html', employees=employees)

@main.route('/employees/add', methods=['GET', 'POST'])
def add_employee():
    if request.method == 'POST':
        name = request.form['name']
        hire_date = datetime.strptime(request.form['hire_date'], '%Y-%m-%d').date()
        current_position = request.form.get('current_position', '')
        
        employee = Employee(
            name=name,
            hire_date=hire_date,
            current_position=current_position
        )
        db.session.add(employee)
        db.session.commit()
        flash('Сотрудник добавлен успешно!', 'success')
        return redirect(url_for('main.employees'))
    return render_template('employees/add.html')

@main.route('/employees/edit/<int:id>', methods=['GET', 'POST'])
def edit_employee(id):
    employee = Employee.query.get_or_404(id)
    if request.method == 'POST':
        employee.name = request.form['name']
        employee.hire_date = datetime.strptime(request.form['hire_date'], '%Y-%m-%d').date()
        employee.current_position = request.form.get('current_position', '')
        db.session.commit()
        flash('Сотрудник обновлен успешно!', 'success')
        return redirect(url_for('main.employees'))
    return render_template('employees/edit.html', employee=employee)

@main.route('/employees/delete/<int:id>')
def delete_employee(id):
    employee = Employee.query.get_or_404(id)
    db.session.delete(employee)
    db.session.commit()
    flash('Сотрудник удален!', 'success')
    return redirect(url_for('main.employees'))

# Маршруты для рабочих
@main.route('/workers')
def workers():
    workers = db.session.query(Worker, Employee).join(
        Employee, Worker.employee_id == Employee.id
    ).all()
    return render_template('workers/list.html', workers=workers)

@main.route('/workers/add', methods=['GET', 'POST'])
def add_worker():
    if request.method == 'POST':
        employee_id = request.form['employee_id']
        hall_id = request.form['hall_id']
        area_id = request.form['area_id']
        work_team_id = request.form['work_team_id']
        category = request.form['category']
        
        worker = Worker(
            employee_id=employee_id,
            hall_id=hall_id,
            area_id=area_id,
            work_team_id=work_team_id,
            category=WorkerCategoryEnum(category)
        )
        db.session.add(worker)
        db.session.commit()
        flash('Рабочий добавлен успешно!', 'success')
        return redirect(url_for('main.workers'))
    
    employees = Employee.query.all()
    halls = ProductionHall.query.all()
    areas = ProductionArea.query.all()
    teams = WorkTeam.query.all()
    categories = WorkerCategoryEnum
    return render_template('workers/add.html', 
                         employees=employees, halls=halls, areas=areas, teams=teams, categories=categories)

@main.route('/workers/edit/<int:id>', methods=['GET', 'POST'])
def edit_worker(id):
    worker = Worker.query.filter_by(employee_id=id).first_or_404()
    if request.method == 'POST':
        worker.employee_id = request.form['employee_id']
        worker.hall_id = request.form['hall_id']
        worker.area_id = request.form['area_id']
        worker.work_team_id = request.form['work_team_id']
        worker.category = WorkerCategoryEnum(request.form['category'])
        db.session.commit()
        flash('Рабочий обновлен успешно!', 'success')
        return redirect(url_for('main.workers'))
    
    employees = Employee.query.all()
    halls = ProductionHall.query.all()
    areas = ProductionArea.query.all()
    teams = WorkTeam.query.all()
    categories = WorkerCategoryEnum
    return render_template('workers/edit.html', worker=worker,
                         employees=employees, halls=halls, areas=areas, teams=teams, categories=categories)

@main.route('/workers/delete/<int:id>')
def delete_worker(id):
    worker = Worker.query.filter_by(employee_id=id).first_or_404()
    db.session.delete(worker)
    db.session.commit()
    flash('Рабочий удален!', 'success')
    return redirect(url_for('main.workers'))

# Маршруты для инженеров
@main.route('/engineers')
def engineers():
    engineers = db.session.query(Engineer, Employee, CategoryEngineer).join(
        Employee, Engineer.employee_id == Employee.id
    ).join(
        CategoryEngineer, Engineer.category_id == CategoryEngineer.id
    ).all()
    return render_template('engineers/list.html', engineers=engineers)

@main.route('/engineers/add', methods=['GET', 'POST'])
def add_engineer():
    if request.method == 'POST':
        employee_id = request.form['employee_id']
        hall_id = request.form['hall_id']
        area_id = request.form['area_id']
        category_id = request.form['category_id']
        
        engineer = Engineer(
            employee_id=employee_id,
            hall_id=hall_id,
            area_id=area_id,
            category_id=category_id
        )
        db.session.add(engineer)
        db.session.commit()
        flash('Инженер добавлен успешно!', 'success')
        return redirect(url_for('main.engineers'))
    
    employees = Employee.query.all()
    halls = ProductionHall.query.all()
    areas = ProductionArea.query.all()
    categories = CategoryEngineer.query.all()
    return render_template('engineers/add.html', 
                         employees=employees, halls=halls, areas=areas, categories=categories)

@main.route('/engineers/edit/<int:id>', methods=['GET', 'POST'])
def edit_engineer(id):
    engineer = Engineer.query.filter_by(employee_id=id).first_or_404()
    if request.method == 'POST':
        engineer.employee_id = request.form['employee_id']
        engineer.hall_id = request.form['hall_id']
        engineer.area_id = request.form['area_id']
        engineer.category_id = request.form['category_id']
        db.session.commit()
        flash('Инженер обновлен успешно!', 'success')
        return redirect(url_for('main.engineers'))
    
    employees = Employee.query.all()
    halls = ProductionHall.query.all()
    areas = ProductionArea.query.all()
    categories = CategoryEngineer.query.all()
    return render_template('engineers/edit.html', engineer=engineer,
                         employees=employees, halls=halls, areas=areas, categories=categories)

@main.route('/engineers/delete/<int:id>')
def delete_engineer(id):
    engineer = Engineer.query.filter_by(employee_id=id).first_or_404()
    db.session.delete(engineer)
    db.session.commit()
    flash('Инженер удален!', 'success')
    return redirect(url_for('main.engineers'))

# Маршруты для изделий
@main.route('/items')
def items():
    items = db.session.query(Item, TypeItem, CategoryItem, ProductionHall).join(
        TypeItem, Item.type_id == TypeItem.id
    ).join(
        CategoryItem, TypeItem.category_id == CategoryItem.id
    ).join(
        ProductionHall, Item.hall_id == ProductionHall.id
    ).all()
    return render_template('items/list.html', items=items)

@main.route('/items/add', methods=['GET', 'POST'])
def add_item():
    if request.method == 'POST':
        name = request.form['name']
        type_id = request.form['type_id']
        hall_id = request.form['hall_id']
        status = request.form['status']
        
        item = Item(
            name=name,
            type_id=type_id,
            hall_id=hall_id,
            status=ItemStatusEnum(status)
        )
        db.session.add(item)
        db.session.commit()
        flash('Изделие добавлено успешно!', 'success')
        return redirect(url_for('main.items'))
    
    types = TypeItem.query.all()
    halls = ProductionHall.query.all()
    statuses = ItemStatusEnum
    return render_template('items/add.html', types=types, halls=halls, statuses=statuses)

@main.route('/items/edit/<int:id>', methods=['GET', 'POST'])
def edit_item(id):
    item = Item.query.get_or_404(id)
    if request.method == 'POST':
        item.name = request.form['name']
        item.type_id = request.form['type_id']
        item.hall_id = request.form['hall_id']
        item.status = ItemStatusEnum(request.form['status'])
        db.session.commit()
        flash('Изделие обновлено успешно!', 'success')
        return redirect(url_for('main.items'))
    
    types = TypeItem.query.all()
    halls = ProductionHall.query.all()
    statuses = ItemStatusEnum
    return render_template('items/edit.html', item=item, types=types, halls=halls, statuses=statuses)

@main.route('/items/delete/<int:id>')
def delete_item(id):
    item = Item.query.get_or_404(id)
    db.session.delete(item)
    db.session.commit()
    flash('Изделие удалено!', 'success')
    return redirect(url_for('main.items'))

# Маршруты для лабораторий
@main.route('/laboratories')
def laboratories():
    labs = TestingLaboratory.query.all()
    return render_template('laboratories/list.html', laboratories=labs)

@main.route('/laboratories/add', methods=['GET', 'POST'])
def add_laboratory():
    if request.method == 'POST':
        name = request.form['name']
        lab = TestingLaboratory(name=name)
        db.session.add(lab)
        db.session.commit()
        flash('Лаборатория добавлена успешно!', 'success')
        return redirect(url_for('main.laboratories'))
    return render_template('laboratories/add.html')

@main.route('/laboratories/edit/<int:id>', methods=['GET', 'POST'])
def edit_laboratory(id):
    laboratory = TestingLaboratory.query.get_or_404(id)
    if request.method == 'POST':
        laboratory.name = request.form['name']
        db.session.commit()
        flash('Лаборатория обновлена успешно!', 'success')
        return redirect(url_for('main.laboratories'))
    return render_template('laboratories/edit.html', laboratory=laboratory)

@main.route('/laboratories/delete/<int:id>')
def delete_laboratory(id):
    laboratory = TestingLaboratory.query.get_or_404(id)
    db.session.delete(laboratory)
    db.session.commit()
    flash('Лаборатория удалена!', 'success')
    return redirect(url_for('main.laboratories'))

# Маршруты для лабораторного оборудования
@main.route('/equipment')
def equipment():
    equipment_list = db.session.query(LabEquip, TestingLaboratory).join(
        TestingLaboratory, LabEquip.lab_id == TestingLaboratory.id
    ).all()
    return render_template('equipment/simple.html', equipment_list=equipment_list)

@main.route('/equipment/add', methods=['GET', 'POST'])
def add_equipment():
    if request.method == 'POST':
        name = request.form['name']
        lab_id = request.form['lab_id']
        
        equipment = LabEquip(name=name, lab_id=lab_id)
        db.session.add(equipment)
        db.session.commit()
        flash('Оборудование добавлено успешно!', 'success')
        return redirect(url_for('main.equipment'))
    
    laboratories = TestingLaboratory.query.all()
    return render_template('equipment/add.html', laboratories=laboratories)

@main.route('/equipment/edit/<int:id>', methods=['GET', 'POST'])
def edit_equipment(id):
    equipment = LabEquip.query.get_or_404(id)
    if request.method == 'POST':
        equipment.name = request.form['name']
        equipment.lab_id = request.form['lab_id']
        db.session.commit()
        flash('Оборудование обновлено успешно!', 'success')
        return redirect(url_for('main.equipment'))
    
    laboratories = TestingLaboratory.query.all()
    return render_template('equipment/edit.html', equipment=equipment, laboratories=laboratories)

@main.route('/equipment/delete/<int:id>')
def delete_equipment(id):
    equipment = LabEquip.query.get_or_404(id)
    db.session.delete(equipment)
    db.session.commit()
    flash('Оборудование удалено!', 'success')
    return redirect(url_for('main.equipment'))

# Маршруты для рабочих бригад
@main.route('/teams')
def teams():
    teams = db.session.query(WorkTeam, ProductionArea, ProductionHall).join(
        ProductionArea, WorkTeam.area_id == ProductionArea.id
    ).join(
        ProductionHall, WorkTeam.hall_id == ProductionHall.id
    ).all()
    return render_template('teams/list.html', teams=teams)

@main.route('/teams/add', methods=['GET', 'POST'])
def add_team():
    if request.method == 'POST':
        name = request.form['name']
        area_id = request.form['area_id']
        hall_id = request.form['hall_id']
        
        team = WorkTeam(name=name, area_id=area_id, hall_id=hall_id)
        db.session.add(team)
        db.session.commit()
        flash('Бригада добавлена успешно!', 'success')
        return redirect(url_for('main.teams'))
    
    areas = ProductionArea.query.all()
    halls = ProductionHall.query.all()
    return render_template('teams/add.html', areas=areas, halls=halls)

@main.route('/teams/edit/<int:id>', methods=['GET', 'POST'])
def edit_team(id):
    team = WorkTeam.query.get_or_404(id)
    if request.method == 'POST':
        team.name = request.form['name']
        team.area_id = request.form['area_id']
        team.hall_id = request.form['hall_id']
        db.session.commit()
        flash('Бригада обновлена успешно!', 'success')
        return redirect(url_for('main.teams'))
    
    areas = ProductionArea.query.all()
    halls = ProductionHall.query.all()
    return render_template('teams/edit.html', team=team, areas=areas, halls=halls)

@main.route('/teams/delete/<int:id>')
def delete_team(id):
    team = WorkTeam.query.get_or_404(id)
    db.session.delete(team)
    db.session.commit()
    flash('Бригада удалена!', 'success')
    return redirect(url_for('main.teams'))

# Маршруты для типов изделий
@main.route('/type_items')
def type_items():
    type_items = db.session.query(TypeItem, CategoryItem).join(
        CategoryItem, TypeItem.category_id == CategoryItem.id
    ).all()
    return render_template('type_items/list.html', type_items=type_items)

@main.route('/type_items/add', methods=['GET', 'POST'])
def add_type_item():
    if request.method == 'POST':
        name = request.form['name']
        category_id = request.form['category_id']
        
        type_item = TypeItem(
            name=name,
            category_id=category_id
        )
        
        try:
            db.session.add(type_item)
            db.session.commit()
            flash('Тип изделия добавлен успешно!', 'success')
            return redirect(url_for('main.type_items'))
        except Exception as e:
            db.session.rollback()
            flash(f'Ошибка при добавлении типа изделия: {str(e)}', 'danger')
    
    categories = CategoryItem.query.all()
    return render_template('type_items/add.html', categories=categories)

@main.route('/type_items/edit/<int:id>', methods=['GET', 'POST'])
def edit_type_item(id):
    type_item = TypeItem.query.get_or_404(id)
    
    if request.method == 'POST':
        type_item.name = request.form['name']
        type_item.category_id = request.form['category_id']
        
        try:
            db.session.commit()
            flash('Тип изделия обновлен успешно!', 'success')
            return redirect(url_for('main.type_items'))
        except Exception as e:
            db.session.rollback()
            flash(f'Ошибка при обновлении типа изделия: {str(e)}', 'danger')
    
    categories = CategoryItem.query.all()
    return render_template('type_items/edit.html', type_item=type_item, categories=categories)

@main.route('/type_items/delete/<int:id>', methods=['POST'])
def delete_type_item(id):
    type_item = TypeItem.query.get_or_404(id)
    
    try:
        db.session.delete(type_item)
        db.session.commit()
        flash('Тип изделия удален успешно!', 'success')
    except Exception as e:
        db.session.rollback()
        flash(f'Ошибка при удалении типа изделия: {str(e)}', 'danger')
    
    return redirect(url_for('main.type_items'))

# Маршруты для готовых изделий
@main.route('/completed_items')
def completed_items():
    completed_items = db.session.query(
        CompletedItem, Item, TypeItem, ProductionHall, WorkTeam, ProductionArea
    ).join(
        Item, CompletedItem.item_id == Item.id
    ).join(
        TypeItem, Item.type_id == TypeItem.id
    ).join(
        ProductionHall, CompletedItem.assembled_in_hall_id == ProductionHall.id
    ).join(
        WorkTeam, CompletedItem.assembled_by_team_id == WorkTeam.id
    ).join(
        ProductionArea, CompletedItem.final_area_id == ProductionArea.id
    ).all()
    return render_template('completed_items/list.html', completed_items=completed_items)

@main.route('/completed_items/add', methods=['GET', 'POST'])
def add_completed_item():
    if request.method == 'POST':
        item_id = request.form['item_id']
        production_start_date = datetime.strptime(request.form['production_start_date'], '%Y-%m-%d').date()
        production_completion_date = datetime.strptime(request.form['production_completion_date'], '%Y-%m-%d').date()
        assembled_by_team_id = request.form['assembled_by_team_id']
        assembled_in_hall_id = request.form['assembled_in_hall_id']
        final_area_id = request.form['final_area_id']
        quantity_produced = request.form['quantity_produced']
        notes = request.form.get('notes', '')
        
        completed_item = CompletedItem(
            item_id=item_id,
            production_start_date=production_start_date,
            production_completion_date=production_completion_date,
            assembled_by_team_id=assembled_by_team_id,
            assembled_in_hall_id=assembled_in_hall_id,
            final_area_id=final_area_id,
            quantity_produced=quantity_produced,
            notes=notes,
            created_at=datetime.now().date(),
            updated_at=datetime.now().date()
        )
        db.session.add(completed_item)
        db.session.commit()
        flash('Готовое изделие добавлено успешно!', 'success')
        return redirect(url_for('main.completed_items'))
    
    items = Item.query.all()
    teams = WorkTeam.query.all()
    halls = ProductionHall.query.all()
    areas = ProductionArea.query.all()
    return render_template('completed_items/add.html', items=items, teams=teams, halls=halls, areas=areas)

@main.route('/completed_items/edit/<int:id>', methods=['GET', 'POST'])
def edit_completed_item(id):
    completed_item = CompletedItem.query.get_or_404(id)
    if request.method == 'POST':
        completed_item.item_id = request.form['item_id']
        completed_item.production_start_date = datetime.strptime(request.form['production_start_date'], '%Y-%m-%d').date()
        completed_item.production_completion_date = datetime.strptime(request.form['production_completion_date'], '%Y-%m-%d').date()
        completed_item.assembled_by_team_id = request.form['assembled_by_team_id']
        completed_item.assembled_in_hall_id = request.form['assembled_in_hall_id']
        completed_item.final_area_id = request.form['final_area_id']
        completed_item.quantity_produced = request.form['quantity_produced']
        completed_item.notes = request.form.get('notes', '')
        completed_item.updated_at = datetime.now().date()
        db.session.commit()
        flash('Готовое изделие обновлено успешно!', 'success')
        return redirect(url_for('main.completed_items'))
    
    items = Item.query.all()
    teams = WorkTeam.query.all()
    halls = ProductionHall.query.all()
    areas = ProductionArea.query.all()
    return render_template('completed_items/edit.html', completed_item=completed_item, items=items, teams=teams, halls=halls, areas=areas)

@main.route('/completed_items/delete/<int:id>')
def delete_completed_item(id):
    completed_item = CompletedItem.query.get_or_404(id)
    db.session.delete(completed_item)
    db.session.commit()
    flash('Готовое изделие удалено!', 'success')
    return redirect(url_for('main.completed_items'))

# Маршруты для испытаний готовых изделий
@main.route('/tested_items')
def tested_items():
    tested_items = db.session.query(
        CompletedItemTest, CompletedItem, Item, TypeItem, TestingLaboratory, Employee
    ).join(
        CompletedItem, CompletedItemTest.completed_item_id == CompletedItem.id
    ).join(
        Item, CompletedItem.item_id == Item.id
    ).join(
        TypeItem, Item.type_id == TypeItem.id
    ).join(
        TestingLaboratory, CompletedItemTest.lab_id == TestingLaboratory.id
    ).join(
        Employee, CompletedItemTest.conducted_by_worker_id == Employee.id
    ).all()
    return render_template('tested_items/list.html', tested_items=tested_items)

@main.route('/tested_items/add', methods=['GET', 'POST'])
def add_tested_item():
    if request.method == 'POST':
        completed_item_id = request.form['completed_item_id']
        lab_id = request.form['lab_id']
        test_start_date = datetime.strptime(request.form['test_start_date'], '%Y-%m-%d').date()
        test_completion_date = None
        if request.form.get('test_completion_date'):
            test_completion_date = datetime.strptime(request.form['test_completion_date'], '%Y-%m-%d').date()
        test_result = request.form.get('test_result', '')
        test_status = request.form['test_status']
        conducted_by_worker_id = request.form['conducted_by_worker_id']
        notes = request.form.get('notes', '')
        
        tested_item = CompletedItemTest(
            completed_item_id=completed_item_id,
            lab_id=lab_id,
            test_start_date=test_start_date,
            test_completion_date=test_completion_date,
            test_result=test_result,
            test_status=test_status,
            conducted_by_worker_id=conducted_by_worker_id,
            notes=notes,
            created_at=datetime.now().date(),
            updated_at=datetime.now().date()
        )
        db.session.add(tested_item)
        db.session.commit()
        flash('Испытание готового изделия добавлено успешно!', 'success')
        return redirect(url_for('main.tested_items'))
    
    completed_items = db.session.query(CompletedItem, Item).join(Item, CompletedItem.item_id == Item.id).all()
    laboratories = TestingLaboratory.query.all()
    lab_workers = db.session.query(LabWorker, Employee).join(Employee, LabWorker.employee_id == Employee.id).all()
    
    return render_template('tested_items/add.html', 
                         completed_items=completed_items, 
                         laboratories=laboratories, 
                         lab_workers=lab_workers)

@main.route('/tested_items/edit/<int:id>', methods=['GET', 'POST'])
def edit_tested_item(id):
    tested_item = CompletedItemTest.query.get_or_404(id)
    if request.method == 'POST':
        tested_item.completed_item_id = request.form['completed_item_id']
        tested_item.lab_id = request.form['lab_id']
        tested_item.test_start_date = datetime.strptime(request.form['test_start_date'], '%Y-%m-%d').date()
        if request.form.get('test_completion_date'):
            tested_item.test_completion_date = datetime.strptime(request.form['test_completion_date'], '%Y-%m-%d').date()
        else:
            tested_item.test_completion_date = None
        tested_item.test_result = request.form.get('test_result', '')
        tested_item.test_status = request.form['test_status']
        tested_item.conducted_by_worker_id = request.form['conducted_by_worker_id']
        tested_item.notes = request.form.get('notes', '')
        tested_item.updated_at = datetime.now().date()
        db.session.commit()
        flash('Испытание готового изделия обновлено успешно!', 'success')
        return redirect(url_for('main.tested_items'))
    
    completed_items = db.session.query(CompletedItem, Item).join(Item, CompletedItem.item_id == Item.id).all()
    laboratories = TestingLaboratory.query.all()
    lab_workers = db.session.query(LabWorker, Employee).join(Employee, LabWorker.employee_id == Employee.id).all()
    
    return render_template('tested_items/edit.html', 
                         tested_item=tested_item, 
                         completed_items=completed_items, 
                         laboratories=laboratories, 
                         lab_workers=lab_workers)

@main.route('/tested_items/delete/<int:id>')
def delete_tested_item(id):
    tested_item = CompletedItemTest.query.get_or_404(id)
    db.session.delete(tested_item)
    db.session.commit()
    flash('Испытание готового изделия удалено!', 'success')
    return redirect(url_for('main.tested_items'))

# Маршруты для управления использованием оборудования в испытаниях
@main.route('/equipment_usage')
def equipment_usage():
    """Список использования оборудования в испытаниях"""
    usage_records = db.session.query(
        TestEquipmentUsage,
        CompletedItemTest,
        CompletedItem,
        Item,
        LabEquip,
        TestingLaboratory,
        Employee
    ).join(
        CompletedItemTest, TestEquipmentUsage.completed_item_test_id == CompletedItemTest.id
    ).join(
        CompletedItem, CompletedItemTest.completed_item_id == CompletedItem.id
    ).join(
        Item, CompletedItem.item_id == Item.id
    ).join(
        LabEquip, TestEquipmentUsage.lab_equip_id == LabEquip.id
    ).join(
        TestingLaboratory, LabEquip.lab_id == TestingLaboratory.id
    ).join(
        LabWorker, CompletedItemTest.conducted_by_worker_id == LabWorker.employee_id
    ).join(
        Employee, LabWorker.employee_id == Employee.id
    ).order_by(TestEquipmentUsage.usage_date.desc()).all()
    
    return render_template('equipment_usage/list.html', usage_records=usage_records)

@main.route('/equipment_usage/add', methods=['GET', 'POST'])
def add_equipment_usage():
    """Добавить новое использование оборудования"""
    if request.method == 'POST':
        usage = TestEquipmentUsage(
            completed_item_test_id=request.form['completed_item_test_id'],
            lab_equip_id=request.form['lab_equip_id'],
            usage_date=datetime.strptime(request.form['usage_date'], '%Y-%m-%d').date(),
            duration_hours=float(request.form['duration_hours']) if request.form['duration_hours'] else None,
            notes=request.form.get('notes'),
            created_at=date.today()
        )
        
        db.session.add(usage)
        db.session.commit()
        flash('Запись об использовании оборудования добавлена!', 'success')
        return redirect(url_for('main.equipment_usage'))
    
    # Получить все активные тесты
    active_tests = db.session.query(
        CompletedItemTest,
        CompletedItem,
        Item,
        TestingLaboratory
    ).join(
        CompletedItem, CompletedItemTest.completed_item_id == CompletedItem.id
    ).join(
        Item, CompletedItem.item_id == Item.id
    ).join(
        TestingLaboratory, CompletedItemTest.lab_id == TestingLaboratory.id
    ).filter(
        CompletedItemTest.test_status.in_(['in_progress', 'passed', 'failed'])
    ).all()
    
    # Получить все оборудование
    equipment = db.session.query(
        LabEquip,
        TestingLaboratory
    ).join(
        TestingLaboratory, LabEquip.lab_id == TestingLaboratory.id
    ).all()
    
    return render_template('equipment_usage/add.html', 
                         active_tests=active_tests, 
                         equipment=equipment)

@main.route('/equipment_usage/edit/<int:id>', methods=['GET', 'POST'])
def edit_equipment_usage(id):
    """Редактировать использование оборудования"""
    usage = TestEquipmentUsage.query.get_or_404(id)
    
    if request.method == 'POST':
        usage.completed_item_test_id = request.form['completed_item_test_id']
        usage.lab_equip_id = request.form['lab_equip_id']
        usage.usage_date = datetime.strptime(request.form['usage_date'], '%Y-%m-%d').date()
        usage.duration_hours = float(request.form['duration_hours']) if request.form['duration_hours'] else None
        usage.notes = request.form.get('notes')
        
        db.session.commit()
        flash('Запись об использовании оборудования обновлена!', 'success')
        return redirect(url_for('main.equipment_usage'))
    
    # Получить все активные тесты
    active_tests = db.session.query(
        CompletedItemTest,
        CompletedItem,
        Item,
        TestingLaboratory
    ).join(
        CompletedItem, CompletedItemTest.completed_item_id == CompletedItem.id
    ).join(
        Item, CompletedItem.item_id == Item.id
    ).join(
        TestingLaboratory, CompletedItemTest.lab_id == TestingLaboratory.id
    ).all()
    
    # Получить все оборудование
    equipment = db.session.query(
        LabEquip,
        TestingLaboratory
    ).join(
        TestingLaboratory, LabEquip.lab_id == TestingLaboratory.id
    ).all()
    
    return render_template('equipment_usage/edit.html', 
                         usage=usage,
                         active_tests=active_tests, 
                         equipment=equipment)

@main.route('/equipment_usage/delete/<int:id>')
def delete_equipment_usage(id):
    """Удалить запись об использовании оборудования"""
    usage = TestEquipmentUsage.query.get_or_404(id)
    db.session.delete(usage)
    db.session.commit()
    flash('Запись об использовании оборудования удалена!', 'success')
    return redirect(url_for('main.equipment_usage'))

# Маршруты для типов работ
@main.route('/work_types')
def work_types():
    work_types = db.session.query(WorkType, ProductionArea, WorkTeam).join(
        ProductionArea, WorkType.area_id == ProductionArea.id
    ).join(
        WorkTeam, WorkType.work_team_id == WorkTeam.id
    ).all()
    return render_template('work_types/list.html', work_types=work_types)

@main.route('/work_types/add', methods=['GET', 'POST'])
def add_work_type():
    if request.method == 'POST':
        work_name = request.form['work_name']
        area_id = request.form['area_id']
        work_team_id = request.form['work_team_id']
        
        work_type = WorkType(
            work_name=work_name,
            area_id=area_id,
            work_team_id=work_team_id
        )
        
        try:
            db.session.add(work_type)
            db.session.commit()
            flash('Тип работы добавлен успешно!', 'success')
            return redirect(url_for('main.work_types'))
        except Exception as e:
            db.session.rollback()
            flash(f'Ошибка при добавлении типа работы: {str(e)}', 'danger')
    
    areas = ProductionArea.query.all()
    teams = WorkTeam.query.all()
    return render_template('work_types/add.html', areas=areas, teams=teams)

@main.route('/work_types/edit/<int:id>', methods=['GET', 'POST'])
def edit_work_type(id):
    work_type = WorkType.query.get_or_404(id)
    
    if request.method == 'POST':
        work_type.work_name = request.form['work_name']
        work_type.area_id = request.form['area_id']
        work_type.work_team_id = request.form['work_team_id']
        
        try:
            db.session.commit()
            flash('Тип работы обновлен успешно!', 'success')
            return redirect(url_for('main.work_types'))
        except Exception as e:
            db.session.rollback()
            flash(f'Ошибка при обновлении типа работы: {str(e)}', 'danger')
    
    areas = ProductionArea.query.all()
    teams = WorkTeam.query.all()
    return render_template('work_types/edit.html', work_type=work_type, areas=areas, teams=teams)

@main.route('/work_types/delete/<int:id>', methods=['POST'])
def delete_work_type(id):
    work_type = WorkType.query.get_or_404(id)
    
    try:
        db.session.delete(work_type)
        db.session.commit()
        flash('Тип работы удален успешно!', 'success')
    except Exception as e:
        db.session.rollback()
        flash(f'Ошибка при удалении типа работы: {str(e)}', 'danger')
    
    return redirect(url_for('main.work_types'))

# Отчеты
# Дополнительные запросы из задания

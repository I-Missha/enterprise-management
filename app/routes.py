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

# API маршруты для динамической загрузки данных
@main.route('/api/areas/<int:hall_id>')
def get_areas_by_hall(hall_id):
    areas = ProductionArea.query.filter_by(hall_id=hall_id).all()
    return jsonify([{'id': area.id, 'name': area.name} for area in areas])

@main.route('/api/teams/<int:area_id>')
def get_teams_by_area(area_id):
    teams = WorkTeam.query.filter_by(area_id=area_id).all()
    return jsonify([{'id': team.id, 'name': team.name} for team in teams])

# Отчеты
# Дополнительные запросы из задания
@main.route('/queries')
def queries():
    """Страница со всеми запросами из задания"""
    return render_template('queries/index.html')

# 1. Получить перечень видов изделий отдельной категории и в целом, собираемых указанным цехом, предприятием
@main.route('/queries/items_by_hall')
def items_by_hall():
    hall_id = request.args.get('hall_id', type=int)
    category_id = request.args.get('category_id', type=int)
    
    query = db.session.query(Item, TypeItem, CategoryItem, ProductionHall).join(
        TypeItem, Item.type_id == TypeItem.id
    ).join(
        CategoryItem, TypeItem.category_id == CategoryItem.id
    ).join(
        ProductionHall, Item.hall_id == ProductionHall.id
    )
    
    if hall_id:
        query = query.filter(Item.hall_id == hall_id)
    if category_id:
        query = query.filter(CategoryItem.id == category_id)
    
    items = query.all()
    halls = ProductionHall.query.all()
    categories = CategoryItem.query.all()
    
    return render_template('queries/items_by_hall.html', 
                         items=items, halls=halls, categories=categories,
                         selected_hall=hall_id, selected_category=category_id)

# 2. Получить число и перечень изделий за определенный отрезок времени
@main.route('/queries/items_by_period')
def items_by_period():
    hall_id = request.args.get('hall_id', type=int)
    area_id = request.args.get('area_id', type=int)
    category_id = request.args.get('category_id', type=int)
    start_date = request.args.get('start_date')
    end_date = request.args.get('end_date')
    
    query = db.session.query(CompletedItem, Item, TypeItem, CategoryItem, ProductionHall, ProductionArea).join(
        Item, CompletedItem.item_id == Item.id
    ).join(
        TypeItem, Item.type_id == TypeItem.id
    ).join(
        CategoryItem, TypeItem.category_id == CategoryItem.id
    ).join(
        ProductionHall, CompletedItem.assembled_in_hall_id == ProductionHall.id
    ).join(
        ProductionArea, CompletedItem.final_area_id == ProductionArea.id
    )
    
    if hall_id:
        query = query.filter(CompletedItem.assembled_in_hall_id == hall_id)
    if area_id:
        query = query.filter(CompletedItem.final_area_id == area_id)
    if category_id:
        query = query.filter(CategoryItem.id == category_id)
    if start_date:
        query = query.filter(CompletedItem.production_completion_date >= start_date)
    if end_date:
        query = query.filter(CompletedItem.production_completion_date <= end_date)
    
    completed_items = query.all()
    halls = ProductionHall.query.all()
    areas = ProductionArea.query.all()
    categories = CategoryItem.query.all()
    
    return render_template('queries/items_by_period.html',
                         completed_items=completed_items, halls=halls, areas=areas, categories=categories,
                         selected_hall=hall_id, selected_area=area_id, selected_category=category_id,
                         start_date=start_date, end_date=end_date)

# 3. Получить данные о кадровом составе цеха, предприятия в целом
@main.route('/queries/personnel_by_hall')
def personnel_by_hall():
    hall_id = request.args.get('hall_id', type=int)
    category_type = request.args.get('category_type')  # worker, engineer, all
    
    workers_query = db.session.query(Worker, Employee, ProductionHall).join(
        Employee, Worker.employee_id == Employee.id
    ).join(
        ProductionHall, Worker.hall_id == ProductionHall.id
    )
    
    engineers_query = db.session.query(Engineer, Employee, ProductionHall, CategoryEngineer).join(
        Employee, Engineer.employee_id == Employee.id
    ).join(
        ProductionHall, Engineer.hall_id == ProductionHall.id
    ).join(
        CategoryEngineer, Engineer.category_id == CategoryEngineer.id
    )
    
    if hall_id:
        workers_query = workers_query.filter(Worker.hall_id == hall_id)
        engineers_query = engineers_query.filter(Engineer.hall_id == hall_id)
    
    workers = workers_query.all() if category_type != 'engineer' else []
    engineers = engineers_query.all() if category_type != 'worker' else []
    
    halls = ProductionHall.query.all()
    
    return render_template('queries/personnel_by_hall.html',
                         workers=workers, engineers=engineers, halls=halls,
                         selected_hall=hall_id, category_type=category_type)

# 4. Получить число и перечень участков указанного цеха и их начальников
@main.route('/queries/areas_and_bosses')
def areas_and_bosses():
    hall_id = request.args.get('hall_id', type=int)
    
    query = db.session.query(ProductionArea, ProductionHall).join(
        ProductionHall, ProductionArea.hall_id == ProductionHall.id
    )
    
    if hall_id:
        query = query.filter(ProductionArea.hall_id == hall_id)
    
    areas = query.all()
    
    # Получаем начальников участков
    area_bosses = db.session.query(AreaBoss, ProductionArea, Engineer, Employee).join(
        ProductionArea, AreaBoss.area_id == ProductionArea.id
    ).join(
        Engineer, AreaBoss.engineer_id == Engineer.employee_id
    ).join(
        Employee, Engineer.employee_id == Employee.id
    )
    
    if hall_id:
        area_bosses = area_bosses.filter(ProductionArea.hall_id == hall_id)
    
    area_bosses = area_bosses.all()
    
    halls = ProductionHall.query.all()
    
    return render_template('queries/areas_and_bosses.html',
                         areas=areas, area_bosses=area_bosses, halls=halls,
                         selected_hall=hall_id)

# 5. Получить перечень работ, которые проходит указанное изделие
@main.route('/queries/item_work_types')
def item_work_types():
    item_id = request.args.get('item_id', type=int)
    
    work_types = []
    if item_id:
        work_types = db.session.query(ItemWorkType, WorkType, ProductionArea, WorkTeam).join(
            WorkType, ItemWorkType.work_type_id == WorkType.id
        ).join(
            ProductionArea, WorkType.area_id == ProductionArea.id
        ).join(
            WorkTeam, WorkType.work_team_id == WorkTeam.id
        ).filter(ItemWorkType.item_id == item_id).order_by(ItemWorkType.seq_number).all()
    
    items = Item.query.all()
    
    return render_template('queries/item_work_types.html',
                         work_types=work_types, items=items, selected_item=item_id)

# 6. Получить состав бригад указанного участка, цеха
@main.route('/queries/teams_composition')
def teams_composition():
    hall_id = request.args.get('hall_id', type=int)
    area_id = request.args.get('area_id', type=int)
    
    teams_query = db.session.query(WorkTeam, ProductionArea, ProductionHall).join(
        ProductionArea, WorkTeam.area_id == ProductionArea.id
    ).join(
        ProductionHall, WorkTeam.hall_id == ProductionHall.id
    )
    
    if hall_id:
        teams_query = teams_query.filter(WorkTeam.hall_id == hall_id)
    if area_id:
        teams_query = teams_query.filter(WorkTeam.area_id == area_id)
    
    teams = teams_query.all()
    
    # Получаем состав бригад
    team_members = {}
    for team, area, hall in teams:
        members = db.session.query(Worker, Employee).join(
            Employee, Worker.employee_id == Employee.id
        ).filter(Worker.work_team_id == team.id).all()
        team_members[team.id] = members
    
    halls = ProductionHall.query.all()
    areas = ProductionArea.query.all()
    
    return render_template('queries/teams_composition.html',
                         teams=teams, team_members=team_members, halls=halls, areas=areas,
                         selected_hall=hall_id, selected_area=area_id)

# 7. Получить список мастеров указанного участка, цеха
@main.route('/queries/masters')
def masters():
    hall_id = request.args.get('hall_id', type=int)
    area_id = request.args.get('area_id', type=int)
    
    query = db.session.query(Masters, ProductionArea, Engineer, Employee, ProductionHall).join(
        ProductionArea, Masters.area_id == ProductionArea.id
    ).join(
        Engineer, Masters.engineer_id == Engineer.employee_id
    ).join(
        Employee, Engineer.employee_id == Employee.id
    ).join(
        ProductionHall, ProductionArea.hall_id == ProductionHall.id
    )
    
    if hall_id:
        query = query.filter(ProductionArea.hall_id == hall_id)
    if area_id:
        query = query.filter(Masters.area_id == area_id)
    
    masters = query.all()
    
    halls = ProductionHall.query.all()
    areas = ProductionArea.query.all()
    
    return render_template('queries/masters.html',
                         masters=masters, halls=halls, areas=areas,
                         selected_hall=hall_id, selected_area=area_id)

# 8. Получить перечень изделий, собираемых в настоящий момент
@main.route('/queries/current_items')
def current_items():
    hall_id = request.args.get('hall_id', type=int)
    area_id = request.args.get('area_id', type=int)
    category_id = request.args.get('category_id', type=int)
    
    query = db.session.query(Item, TypeItem, CategoryItem, ProductionHall).join(
        TypeItem, Item.type_id == TypeItem.id
    ).join(
        CategoryItem, TypeItem.category_id == CategoryItem.id
    ).join(
        ProductionHall, Item.hall_id == ProductionHall.id
    ).filter(Item.status == ItemStatusEnum.in_progress)
    
    if hall_id:
        query = query.filter(Item.hall_id == hall_id)
    if category_id:
        query = query.filter(CategoryItem.id == category_id)
    
    # Если указан участок, получаем изделия через связь с участками
    if area_id:
        query = query.join(AreasItems, Item.id == AreasItems.item_id).filter(AreasItems.area_id == area_id)
    
    current_items = query.all()
    
    halls = ProductionHall.query.all()
    areas = ProductionArea.query.all()
    categories = CategoryItem.query.all()
    
    return render_template('queries/current_items.html',
                         current_items=current_items, halls=halls, areas=areas, categories=categories,
                         selected_hall=hall_id, selected_area=area_id, selected_category=category_id)

# 9. Получить состав бригад, участвующих в сборке указанного изделия
@main.route('/queries/item_teams')
def item_teams():
    item_id = request.args.get('item_id', type=int)
    
    teams = []
    if item_id:
        teams = db.session.query(WorkTeam, ProductionArea, WorkType).join(
            WorkType, WorkTeam.id == WorkType.work_team_id
        ).join(
            ProductionArea, WorkTeam.area_id == ProductionArea.id
        ).join(
            ItemWorkType, WorkType.id == ItemWorkType.work_type_id
        ).filter(ItemWorkType.item_id == item_id).distinct().all()
    
    items = Item.query.all()
    
    return render_template('queries/item_teams.html',
                         teams=teams, items=items, selected_item=item_id)

# 10. Получить перечень испытательных лабораторий для изделия
@main.route('/queries/item_laboratories')
def item_laboratories():
    item_id = request.args.get('item_id', type=int)
    
    laboratories = []
    if item_id:
        laboratories = db.session.query(TestingLaboratory, ItemTests).join(
            ItemTests, TestingLaboratory.id == ItemTests.lab_equip_id
        ).join(
            LabEquip, ItemTests.lab_equip_id == LabEquip.id
        ).filter(ItemTests.item_id == item_id).distinct().all()
    
    items = Item.query.all()
    
    return render_template('queries/item_laboratories.html',
                         laboratories=laboratories, items=items, selected_item=item_id)

# 11. Получить перечень изделий, проходивших испытание в лаборатории за период
@main.route('/queries/lab_tested_items')
def lab_tested_items():
    lab_id = request.args.get('lab_id', type=int)
    category_id = request.args.get('category_id', type=int)
    start_date = request.args.get('start_date')
    end_date = request.args.get('end_date')
    
    query = db.session.query(Item, TypeItem, CategoryItem, ItemTests, TestingLaboratory).join(
        TypeItem, Item.type_id == TypeItem.id
    ).join(
        CategoryItem, TypeItem.category_id == CategoryItem.id
    ).join(
        ItemTests, Item.id == ItemTests.item_id
    ).join(
        LabEquip, ItemTests.lab_equip_id == LabEquip.id
    ).join(
        TestingLaboratory, LabEquip.lab_id == TestingLaboratory.id
    )
    
    if lab_id:
        query = query.filter(TestingLaboratory.id == lab_id)
    if category_id:
        query = query.filter(CategoryItem.id == category_id)
    if start_date:
        query = query.filter(ItemTests.test_date >= start_date)
    if end_date:
        query = query.filter(ItemTests.test_date <= end_date)
    
    tested_items = query.all()
    
    laboratories = TestingLaboratory.query.all()
    categories = CategoryItem.query.all()
    
    return render_template('queries/lab_tested_items.html',
                         tested_items=tested_items, laboratories=laboratories, categories=categories,
                         selected_lab=lab_id, selected_category=category_id,
                         start_date=start_date, end_date=end_date)

# 12. Получить список испытателей для изделий в лаборатории за период
@main.route('/queries/lab_testers')
def lab_testers():
    lab_id = request.args.get('lab_id', type=int)
    item_id = request.args.get('item_id', type=int)
    category_id = request.args.get('category_id', type=int)
    start_date = request.args.get('start_date')
    end_date = request.args.get('end_date')
    
    query = db.session.query(LabWorker, Employee, TestingLaboratory, ItemTests, Item, TypeItem, CategoryItem).join(
        Employee, LabWorker.employee_id == Employee.id
    ).join(
        TestingLaboratory, LabWorker.lab_id == TestingLaboratory.id
    ).join(
        ItemTests, LabWorker.employee_id == ItemTests.lab_worker_id
    ).join(
        Item, ItemTests.item_id == Item.id
    ).join(
        TypeItem, Item.type_id == TypeItem.id
    ).join(
        CategoryItem, TypeItem.category_id == CategoryItem.id
    )
    
    if lab_id:
        query = query.filter(LabWorker.lab_id == lab_id)
    if item_id:
        query = query.filter(ItemTests.item_id == item_id)
    if category_id:
        query = query.filter(CategoryItem.id == category_id)
    if start_date:
        query = query.filter(ItemTests.test_date >= start_date)
    if end_date:
        query = query.filter(ItemTests.test_date <= end_date)
    
    testers = query.distinct().all()
    
    laboratories = TestingLaboratory.query.all()
    items = Item.query.all()
    categories = CategoryItem.query.all()
    
    return render_template('queries/lab_testers.html',
                         testers=testers, laboratories=laboratories, items=items, categories=categories,
                         selected_lab=lab_id, selected_item=item_id, selected_category=category_id,
                         start_date=start_date, end_date=end_date)

# 13. Получить состав оборудования для испытаний за период
@main.route('/queries/lab_equipment_usage')
def lab_equipment_usage():
    lab_id = request.args.get('lab_id', type=int)
    item_id = request.args.get('item_id', type=int)
    category_id = request.args.get('category_id', type=int)
    start_date = request.args.get('start_date')
    end_date = request.args.get('end_date')
    
    query = db.session.query(LabEquip, TestingLaboratory, ItemTests, Item, TypeItem, CategoryItem).join(
        TestingLaboratory, LabEquip.lab_id == TestingLaboratory.id
    ).join(
        ItemTests, LabEquip.id == ItemTests.lab_equip_id
    ).join(
        Item, ItemTests.item_id == Item.id
    ).join(
        TypeItem, Item.type_id == TypeItem.id
    ).join(
        CategoryItem, TypeItem.category_id == CategoryItem.id
    )
    
    if lab_id:
        query = query.filter(LabEquip.lab_id == lab_id)
    if item_id:
        query = query.filter(ItemTests.item_id == item_id)
    if category_id:
        query = query.filter(CategoryItem.id == category_id)
    if start_date:
        query = query.filter(ItemTests.test_date >= start_date)
    if end_date:
        query = query.filter(ItemTests.test_date <= end_date)
    
    equipment_usage = query.distinct().all()
    
    laboratories = TestingLaboratory.query.all()
    items = Item.query.all()
    categories = CategoryItem.query.all()
    
    return render_template('queries/lab_equipment_usage.html',
                         equipment_usage=equipment_usage, laboratories=laboratories, items=items, categories=categories,
                         selected_lab=lab_id, selected_item=item_id, selected_category=category_id,
                         start_date=start_date, end_date=end_date)

# 14. Получить число и перечень изделий, собираемых в настоящее время (дублирует 8-й запрос)
# Используем тот же маршрут current_items

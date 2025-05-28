from functools import wraps
from flask import session, redirect, url_for, request, flash
from werkzeug.security import check_password_hash, generate_password_hash

# Simple authentication decorator
def login_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        if 'user_id' not in session:
            flash('Пожалуйста, войдите в систему для доступа к этой странице.', 'warning')
            return redirect(url_for('main.login', next=request.url))
        return f(*args, **kwargs)
    return decorated_function

# Simple users storage (in production, this should be in the database)
USERS = {
    'admin': {
        'password_hash': generate_password_hash('admin123'),
        'role': 'admin'
    },
    'manager': {
        'password_hash': generate_password_hash('manager123'),
        'role': 'manager'
    },
    'operator': {
        'password_hash': generate_password_hash('operator123'),
        'role': 'operator'
    }
}

def authenticate_user(username, password):
    """Authenticate user with username and password"""
    user = USERS.get(username)
    if user and check_password_hash(user['password_hash'], password):
        return {'username': username, 'role': user['role']}
    return None

def get_current_user():
    """Get current logged in user"""
    if 'user_id' in session:
        return {
            'username': session['user_id'],
            'role': session.get('user_role', 'operator')
        }
    return None

def has_permission(required_role):
    """Check if current user has required permission"""
    user = get_current_user()
    if not user:
        return False
    
    role_hierarchy = {
        'operator': 1,
        'manager': 2,
        'admin': 3
    }
    
    user_level = role_hierarchy.get(user['role'], 0)
    required_level = role_hierarchy.get(required_role, 3)
    
    return user_level >= required_level

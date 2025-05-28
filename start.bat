
python --version >nul 2>&1
if %errorlevel% neq 0 (
    pause
    exit /b 1
)

if not exist "venv" (
    python -m venv venv
)


call venv\Scripts\activate.bat


pip install -r requirements.txt



python run.py

pause

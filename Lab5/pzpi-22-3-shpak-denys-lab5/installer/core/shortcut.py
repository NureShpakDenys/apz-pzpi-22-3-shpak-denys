import pythoncom
from win32com.client import Dispatch
import os

def create_shortcut(target, arguments, shortcut_path, working_dir):
    pythoncom.CoInitialize()
    shell_dispatch = Dispatch("WScript.Shell")
    shortcut = shell_dispatch.CreateShortcut(str(shortcut_path))
    shortcut.TargetPath = str(target)
    shortcut.Arguments = arguments
    shortcut.WorkingDirectory = str(working_dir)
    shortcut.IconLocation = str(target)
    shortcut.Save()

def create_start_bat(install_path):
    bat_path = os.path.join(install_path, "start.bat")
    frontend_path = os.path.join(install_path, "frontend/wayra")
    main_exe_path = os.path.join(install_path, "main.exe")

    content = f"""
        @echo off
        start "" "{main_exe_path}" --config "{os.path.join(install_path, 'config', 'config.yaml')}"

        cd /d "{frontend_path}"

        if not exist node_modules (
            echo Installing npm dependencies...
            npm install
        )

        npm start
    """

    with open(bat_path, "w", encoding="utf-8") as f:
        f.write(content)

    return bat_path

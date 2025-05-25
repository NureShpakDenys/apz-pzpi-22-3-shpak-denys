import os
import yaml
from PyQt5.QtCore import QThread, pyqtSignal
from pathlib import Path
from core.downloader import download_server_exe, download_and_extract_frontend
from core.shortcut import create_shortcut, create_start_bat
from win32com.shell import shell, shellcon
from core.backup import restore_database

class InstallerThread(QThread):
    step_completed = pyqtSignal(str)
    finished = pyqtSignal()
    error = pyqtSignal(str)

    def __init__(self, config_data, install_path):
        super().__init__()
        self.config_data = config_data
        self.install_path = install_path

    def run(self):
        try:
            os.makedirs(self.install_path, exist_ok=True)

            frontend_path = os.path.join(self.install_path, "frontend")
            os.makedirs(frontend_path, exist_ok=True)

            download_and_extract_frontend(frontend_path)
            self.step_completed.emit("Frontend extracted ✅")

            download_server_exe(self.install_path)
            self.step_completed.emit("Downloading main.exe ✅")

            config_dir = os.path.join(self.install_path, "config")
            os.makedirs(config_dir, exist_ok=True)

            config_path = os.path.join(config_dir, "config.yaml")
            with open(config_path, 'w') as f:
                yaml.dump(self.config_data["config"], f)
            self.step_completed.emit("Saving config.yaml ✅")

            if self.config_data.get("restore_backup") and self.config_data.get("backup_path"):
                self.step_completed.emit("Restoring from backup...")

            db_cfg = self.config_data["config"]["database"]
            try:
                restore_database(
                    backup_dir=self.config_data["backup_path"],
                    db_user=db_cfg["user"],
                    db_name=db_cfg["name"],
                    db_password=db_cfg["password"],
                    encryption_key=self.config_data["config"]["encryption_key"]
                )
                self.step_completed.emit("Database restored from backup ✅")
            except Exception as e:
                raise Exception(f"Backup restore failed: {e}")

            bat_path = create_start_bat(self.install_path)
            self.step_completed.emit("Created start.bat ✅")

            desktop = Path(shell.SHGetFolderPath(0, shellcon.CSIDL_DESKTOP, None, 0))
            shortcut_path = desktop / "LaunchApp.lnk"

            create_shortcut(
                target=bat_path,
                arguments="",
                shortcut_path=shortcut_path,
                working_dir=self.install_path
            )
            self.step_completed.emit("Creating shortcut ✅")

            self.step_completed.emit("All dependencies downloaded ✅")
            self.finished.emit()
        except Exception as e:
            self.error.emit(str(e))

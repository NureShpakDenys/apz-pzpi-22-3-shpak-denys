from PyQt5.QtWidgets import QWidget, QVBoxLayout, QLabel, QProgressBar
from PyQt5.QtCore import Qt, pyqtSlot
from PyQt5.QtGui import QFont
from core import InstallerThread
from PyQt5.QtWidgets import (
    QWidget, QLabel, QVBoxLayout, QProgressBar
)
import subprocess
import threading

class InstallPage(QWidget):
    def __init__(self, parent, config_data, install_path):
        super().__init__(parent)
        self.config_data = config_data
        self.install_path = install_path

        self.layout = QVBoxLayout()
        self.title = QLabel("Installing Components")
        self.title.setFont(QFont("Arial", 32, QFont.Bold))
        self.title.setAlignment(Qt.AlignCenter)

        self.status_label = QLabel("Installation in progress...")
        self.status_label.setFont(QFont("Arial", 20))

        self.progress_bar = QProgressBar()
        self.progress_bar.setFont(QFont("Arial", 18))

        self.log_label = QLabel("")
        self.log_label.setFont(QFont("Arial", 16))
        self.log_label.setWordWrap(True)

        self.layout.addWidget(self.title)
        self.layout.addWidget(self.status_label)
        self.layout.addWidget(self.progress_bar)
        self.layout.addWidget(self.log_label)

        self.setLayout(self.layout)

        self.install_steps = [
            ("Installing Go", ["winget", "install", "-e", "--id", "GoLang.Go"]),
            ("Installing Node.js", ["winget", "install", "-e", "--id", "OpenJS.NodeJS.LTS"]),
            ("Installing PostgreSQL", ["winget", "install", "-e", "--id", "PostgreSQL.PostgreSQL"]),
        ]
        self.progress_bar.setMaximum(len(self.install_steps) + 6)
        self.current_step = 0
        self.has_started = False

    def start_installation(self):
        if self.has_started:
            return
        self.has_started = True
        threading.Thread(target=self.run_dependency_installation).start()

    def run_dependency_installation(self):
        for desc, cmd in self.install_steps:
            self.append_log(f"{desc}...")
            try:
                output = subprocess.check_output(cmd, stderr=subprocess.STDOUT, text=True, shell=True)
                self.append_log(output.strip())
            except subprocess.CalledProcessError as e:
                self.append_log(f"Error during {desc}:\n{e.output.strip()}")
            self.current_step += 1
            self.progress_bar.setValue(self.current_step)

        self.append_log("System dependencies installed. Proceeding with internal setup...")
        self.start_internal_installer()

    def start_internal_installer(self):
        self.installer = InstallerThread(self.config_data, self.install_path)
        self.installer.step_completed.connect(self.on_step_completed)
        self.installer.finished.connect(self.on_finished)
        self.installer.error.connect(self.on_error)
        self.installer.start()

    @pyqtSlot(str)
    def on_step_completed(self, message):
        self.append_log(message)
        self.current_step += 1
        self.progress_bar.setValue(self.current_step)

    @pyqtSlot()
    def on_finished(self):
        self.status_label.setText("Installation complete!")

    @pyqtSlot(str)
    def on_error(self, message):
        self.status_label.setText("Installation failed!")
        self.append_log(f"❌ Error: {message}")

    def append_log(self, text):
        print(text)
        self.log_label.setText(text)
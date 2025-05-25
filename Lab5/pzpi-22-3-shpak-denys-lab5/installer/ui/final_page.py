import os
import subprocess
from PyQt5.QtCore import Qt
from PyQt5.QtGui import QFont
from PyQt5.QtWidgets import (
    QApplication, QWidget, QLabel,
    QVBoxLayout, QCheckBox, QPushButton,
)

class FinalPage(QWidget):
    def __init__(self):
        super().__init__()
        self.install_path = None
        self.setLayout(self._build_layout())

        layout = QVBoxLayout()

        self.title = QLabel("Installation Complete")
        self.title.setFont(QFont("Arial", 32, QFont.Bold))
        self.title.setAlignment(Qt.AlignCenter)

        self.message = QLabel("Everything is ready.")
        self.message.setFont(QFont("Arial", 20))

        self.run_checkbox = QCheckBox("Launch server and web client")
        self.run_checkbox.setFont(QFont("Arial", 18))

        self.finish_button = QPushButton("Finish")
        self.finish_button.setFont(QFont("Arial", 18))
        self.finish_button.clicked.connect(self.finish)

        layout.addWidget(self.title)
        layout.addWidget(self.message)
        layout.addWidget(self.run_checkbox)
        layout.addWidget(self.finish_button)

        self.setLayout(layout)

    def _build_layout(self):
        layout = QVBoxLayout()
        self.title = QLabel("Installation Complete")
        self.title.setFont(QFont("Arial", 32, QFont.Bold))
        self.title.setAlignment(Qt.AlignCenter)

        self.message = QLabel("Everything is ready.")
        self.message.setFont(QFont("Arial", 20))

        self.run_checkbox = QCheckBox("Launch server and web client")
        self.run_checkbox.setFont(QFont("Arial", 18))

        self.finish_button = QPushButton("Finish")
        self.finish_button.setFont(QFont("Arial", 18))
        self.finish_button.clicked.connect(self.finish)

        layout.addWidget(self.title)
        layout.addWidget(self.message)
        layout.addWidget(self.run_checkbox)
        layout.addWidget(self.finish_button)
        return layout

    def finish(self):
        if self.run_checkbox.isChecked() and self.install_path:
            bat_path = os.path.join(self.install_path, "start.bat")
            subprocess.Popen([bat_path], shell=True)
        QApplication.quit()

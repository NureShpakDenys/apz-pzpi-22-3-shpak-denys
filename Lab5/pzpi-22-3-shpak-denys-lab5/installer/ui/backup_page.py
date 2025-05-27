from PyQt5.QtCore import Qt
from PyQt5.QtGui import QFont
from PyQt5.QtWidgets import (
    QWidget, QLabel, QFileDialog,
    QVBoxLayout, QCheckBox, QPushButton,
)

class BackupPage(QWidget):
    def __init__(self, parent):
        super().__init__()
        self.parent = parent
        self.install_path = ""
        self.backup_path = ""

        layout = QVBoxLayout()

        self.title = QLabel("Backup Settings")
        self.title.setFont(QFont("Arial", 32, QFont.Bold))
        self.title.setAlignment(Qt.AlignCenter)

        self.install_path_button = QPushButton("Choose installation path")
        self.install_path_button.setFont(QFont("Arial", 18))
        self.install_path_button.clicked.connect(self.select_install_path)

        self.install_path_label = QLabel("No install path selected")
        self.install_path_label.setFont(QFont("Arial", 16))

        self.backup_checkbox = QCheckBox("Install from backup")
        self.backup_checkbox.setFont(QFont("Arial", 20))
        self.backup_checkbox.stateChanged.connect(self.toggle_backup_path)

        self.path_button = QPushButton("Choose backup folder")
        self.path_button.setFont(QFont("Arial", 18))
        self.path_button.setEnabled(False)
        self.path_button.clicked.connect(self.select_backup_path)

        self.selected_path_label = QLabel("No backup path selected")
        self.selected_path_label.setFont(QFont("Arial", 16))

        self.next_button = QPushButton("Next")
        self.next_button.setFont(QFont("Arial", 18))
        self.next_button.clicked.connect(self.go_to_next)

        layout.addWidget(self.title)

        layout.addWidget(self.install_path_button)
        layout.addWidget(self.install_path_label)

        layout.addWidget(self.backup_checkbox)
        layout.addWidget(self.path_button)
        layout.addWidget(self.selected_path_label)
        layout.addWidget(self.next_button)

        self.setLayout(layout)

    def toggle_backup_path(self):
        self.path_button.setEnabled(self.backup_checkbox.isChecked())

    def select_backup_path(self):
        path = QFileDialog.getExistingDirectory(self, "Select backup folder")
        if path:
            self.backup_path = path
            self.selected_path_label.setText(path)

    def select_install_path(self):
        path = QFileDialog.getExistingDirectory(self, "Select installation folder")
        if path:
            self.install_path = path
            self.install_path_label.setText(path)

    def go_to_next(self):
        if not self.install_path:
            self.install_path_label.setText("❗ Select installation path")
            return
        self.parent.setCurrentIndex(2)
from PyQt5.QtWidgets import QWidget, QVBoxLayout, QLabel, QLineEdit, QPushButton
from PyQt5.QtCore import Qt
from PyQt5.QtGui import QFont
import requests

class LoginPage(QWidget):
    def __init__(self, parent):
        super().__init__()
        self.parent = parent
        layout = QVBoxLayout()

        self.title = QLabel("System Installer")
        self.title.setFont(QFont("Arial", 32, QFont.Bold))
        self.title.setAlignment(Qt.AlignCenter)

        self.instruction = QLabel("Enter your credentials to continue:")
        self.instruction.setFont(QFont("Arial", 20))
        self.instruction.setAlignment(Qt.AlignCenter)

        self.login_input = QLineEdit()
        self.login_input.setPlaceholderText("Username")
        self.login_input.setFont(QFont("Arial", 18))

        self.password_input = QLineEdit()
        self.password_input.setPlaceholderText("Password")
        self.password_input.setEchoMode(QLineEdit.Password)
        self.password_input.setFont(QFont("Arial", 18))

        self.login_button = QPushButton("Next")
        self.login_button.setFont(QFont("Arial", 18))
        self.login_button.clicked.connect(self.check_credentials)

        layout.addWidget(self.title)
        layout.addWidget(self.instruction)
        layout.addWidget(self.login_input)
        layout.addWidget(self.password_input)
        layout.addWidget(self.login_button)

        self.setLayout(layout)

    def check_credentials(self):
        login = self.login_input.text()
        password = self.password_input.text()
        if login and password:
            try:
                response = requests.post("http://localhost:8080/get-creds", json={
                    "username": login,
                    "password": password
                })
                response.raise_for_status()
                self.parent.config_data = response.json()
                self.parent.setCurrentIndex(1)
            except Exception as e:
                print(f"Login failed: {e}")
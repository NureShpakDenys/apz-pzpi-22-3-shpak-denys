from PyQt5.QtWidgets import QStackedWidget
from ui.login_page import LoginPage
from ui.backup_page import BackupPage
from ui.install_page import InstallPage
from ui.final_page import FinalPage

class InstallerApp(QStackedWidget):
    def __init__(self):
        super().__init__()

        self.login_page = LoginPage(self)
        self.backup_page = BackupPage(self)
        self.install_page = InstallPage(self)
        self.final_page = FinalPage()

        self.addWidget(self.login_page)
        self.addWidget(self.backup_page)
        self.addWidget(self.install_page)
        self.addWidget(self.final_page)

        self.setFixedSize(1200, 800)
        self.setWindowTitle("System Installer")

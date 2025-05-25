import sys
from PyQt5.QtWidgets import QApplication
from ui.installer_app import InstallerApp



if __name__ == '__main__':
    app = QApplication(sys.argv)
    installer = InstallerApp()
    installer.show()
    sys.exit(app.exec_())
import sys
from PyQt5.QtWidgets import QApplication
from ui.installer_app import InstallerApp

if '--help' in sys.argv or '-h' in sys.argv:
    print("Usage: python main.py [options]")
    print("Description: Installer for \"Wayra\" - logistic management company platform.")
    print("\nOptions:")
    print("  -h, --help         Show this help message and exit")
    sys.exit(0)

if __name__ == '__main__':
    app = QApplication(sys.argv)
    installer = InstallerApp()
    installer.show()
    sys.exit(app.exec_())
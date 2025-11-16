import os
import winreg

exe_path = os.path.join(os.getcwd(), "customclip.exe")

def add_to_startup(name, path):
    reg_path = r"Software\Microsoft\Windows\CurrentVersion\Run"
    key = winreg.OpenKey(winreg.HKEY_CURRENT_USER, reg_path, 0, winreg.KEY_SET_VALUE)
    winreg.SetValueEx(key, name, 0, winreg.REG_SZ, path)
    winreg.CloseKey(key)

add_to_startup("CustomClip", exe_path)
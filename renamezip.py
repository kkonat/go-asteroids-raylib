import subprocess
import re
import shutil

def get_version():
    try:
        output = subprocess.check_output(["bangbang.exe", "-v"], stderr=subprocess.STDOUT, shell=True)
        output = output.decode("utf-8")
        version = re.search(r"bang bang v(\d+)\.(\d+)\.(\d+)", output)
        return version.group(1), version.group(2), version.group(3)
    except Exception as e:
        print(e)
        return None, None, None

def rename_game_zip(version):
    try:
        shutil.move("game.zip", f"bangbang_{version[0]}.{version[1]}.{version[2]}.x64.zip")
    except Exception as e:
        print(e)

if __name__ == "__main__":
    version = get_version()
    if version[0] is not None:
        rename_game_zip(version)
    else:
        print("Failed to get version")
    

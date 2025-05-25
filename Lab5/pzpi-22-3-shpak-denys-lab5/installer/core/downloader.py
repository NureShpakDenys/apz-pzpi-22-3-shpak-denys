import os
import requests
import zipfile
import io

def download_server_exe(install_path):
    url = "http://localhost:8080/get-server-exe"
    response = requests.get(url, stream=True)
    if response.status_code == 200:
        exe_path = os.path.join(install_path, "main.exe")
        with open(exe_path, 'wb') as f:
            for chunk in response.iter_content(chunk_size=8192):
                f.write(chunk)
        print("main.exe downloaded successfully")
    else:
        print(f"Failed to download main.exe: {response.status_code}")
        response.raise_for_status()


def download_and_extract_frontend(destination_path):
    url = "http://localhost:8080/get-frontend-zip"
    response = requests.get(url)
    if response.status_code == 200:
        with zipfile.ZipFile(io.BytesIO(response.content)) as zip_ref:
            zip_ref.extractall(destination_path)
        print("Frontend extracted successfully")
    else:
        raise Exception(f"Failed to download frontend: {response.status_code}")
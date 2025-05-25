import os
import shutil
import subprocess
import tempfile
from cryptography.hazmat.primitives.ciphers.aead import AESGCM
from pathlib import Path

backup_tables = [
    "roles", "users", "companies", "routes", "deliveries",
    "product_categories", "products", "waypoints", "sensor_data", "user_companies"
]

def decrypt_file_to(source_path, dest_path, key: bytes):
    with open(source_path, "rb") as f:
        data = f.read()

    nonce_size = 12
    nonce = data[:nonce_size]
    ciphertext = data[nonce_size:]

    aesgcm = AESGCM(key)
    plaintext = aesgcm.decrypt(nonce, ciphertext, None)

    with open(dest_path, "wb") as f:
        f.write(plaintext)

def restore_database(backup_dir: str, db_user: str, db_name: str, db_password: str, encryption_key: str):
    temp_dir = tempfile.mkdtemp()
    os.environ["PGPASSWORD"] = db_password

    try:
        for table in backup_tables:
            encrypted = os.path.join(backup_dir, f"{table}.csv")
            decrypted = os.path.join(temp_dir, f"{table}.csv")

            decrypt_file_to(encrypted, decrypted, bytes(encryption_key, encoding="utf-8"))

            truncate_cmd = [
                "psql", "-U", db_user, "-d", db_name, "-h", "localhost", "-p", "5432",
                "-c", f'TRUNCATE TABLE {table} RESTART IDENTITY CASCADE;'
            ]
            subprocess.run(truncate_cmd, check=True)

            import_cmd = [
                "psql", "-U", db_user, "-d", db_name, "-h", "localhost", "-p", "5432",
                "-c", f"\\COPY {table} FROM '{decrypted}' WITH CSV HEADER"
            ]
            subprocess.run(import_cmd, check=True)

    finally:
        shutil.rmtree(temp_dir)
        del os.environ["PGPASSWORD"]

import os
import paramiko
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Retrieve SFTP credentials
sftp_host = os.getenv("SFTP_HOST")
sftp_user = os.getenv("SFTP_USER")
private_key_path = os.path.expanduser(os.getenv("SFTP_PRIVATE_KEY"))
sftp_port = int(os.getenv("SFTP_PORT", 22))

def upload_to_sftp(local_file_path, remote_file_path):
    """Uploads a file to AWS Transfer Family (SFTP)"""
    try:
        if not os.path.exists(local_file_path):
            raise FileNotFoundError(f"❌ File not found: {local_file_path}")

        print(f"📂 Uploading to: {remote_file_path}")

        # Load SSH private key
        private_key = paramiko.RSAKey(filename=private_key_path)

        # Establish SFTP connection
        transport = paramiko.Transport((sftp_host, sftp_port))
        transport.connect(username=sftp_user, pkey=private_key)
        sftp = paramiko.SFTPClient.from_transport(transport)

        # Simply upload to the full path - AWS Transfer Family will handle it
        sftp.put(local_file_path, remote_file_path)
        print(f"✅ File uploaded successfully: {remote_file_path}")

        # Close connection
        sftp.close()
        transport.close()

    except Exception as e:
        print(f"❌ SFTP Upload Failed: {e}")

# Upload all generated transaction logs
base_dir = "./transaction-logs/"

for client_id in os.listdir(base_dir):
    client_path = os.path.join(base_dir, client_id)

    if os.path.isdir(client_path):  # Ensure it's a directory
        for filename in os.listdir(client_path):
            local_file_path = os.path.join(client_path, filename)
            remote_file_path = f"logs/{client_id}/{filename}"  # Upload inside `/sftp_user/ClientID/`
            
            # Upload file
            upload_to_sftp(local_file_path, remote_file_path)

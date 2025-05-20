import requests

SERVER_URL = "https://localhost:8080/upload?filename=test.json"

CLIENT_CERT = ("client.crt", "client.key")
CA_CERT = "ca.crt"

with open("test.json", "rb") as f:
    files = {"file": ("test.json", f, "application/json")}
    try:
        resp = requests.post(
            SERVER_URL,
            files=files,
            cert=CLIENT_CERT,
            verify=CA_CERT,
            timeout=10
        )
        print("Status:", resp.status_code)
        print("Response:", resp.text)
    except requests.exceptions.SSLError as e:
        print("SSL Error:", e)
    except Exception as e:
        print("Error:", e)
import requests

HOST_PORT = "http://localhost:10168"

def service_find(name):
    requests.get(f"{HOST_PORT}/service-manager/find?name={name}")

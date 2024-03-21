import requests

HOST_PORT = "http://127.0.0.1:51126"


def service_find(name: str):
    response = requests.get(f"{HOST_PORT}/service-manager/find?name={name}")
    if response.status_code != 200:
        raise Exception(
            f"call kasos api /service-manager/find failed, code: {response.status_code}, msg: {response.text}"
        )
    return response.json()


def list_services(index: int):
    response = requests.get(f"{HOST_PORT}/service-manager/list?index={index}")
    if response.status_code != 200:
        raise Exception(
            f"call kasos api /service-manager/list failed, code: {response.status_code}, msg: {response.text}"
        )
    return response.json()


def register_service(name: str, tags: dict[str, str], model_name: str):
    response = requests.post(f"{HOST_PORT}/service-manager/register",
                             json={
                                 "name": name,
                                 "tags": tags,
                                 "model_name": model_name
                             })
    if response.status_code != 200:
        raise Exception(
            f"call kasos api /service-manager/register failed, code: {response.status_code}, msg: {response.text}"
        )
    return response.json()


def register_service_result(name: str):
    response = requests.get(
        f"{HOST_PORT}/service-manager/register-result?name={name}")
    if response.status_code != 200:
        raise Exception(
            f"call kasos api /service-manager/register-result failed, code: {response.status_code}, msg: {response.text}"
        )
    return response.json()


def delete_service(name: str):
    response = requests.post(f"{HOST_PORT}/service-manager/delete?name={name}")
    if response.status_code != 200:
        raise Exception(
            f"call kasos api /service-manager/delete failed, code: {response.status_code}, msg: {response.text}"
        )
    return response.json()


def find_model(name: str):
    response = requests.get(f"{HOST_PORT}/model-manager/find?name={name}")
    if response.status_code != 200:
        raise Exception(
            f"call kasos api /model-manager/find failed, code: {response.status_code}, msg: {response.text}"
        )
    return response.json()


def list_models(index: int):
    response = requests.get(f"{HOST_PORT}/model-manager/list?index={index}")
    if response.status_code != 200:
        raise Exception(
            f"call kasos api /model-manager/list failed, code: {response.status_code}, msg: {response.text}"
        )
    return response.json()


def register_model(name: str, train_script: str, infer_script: str):
    response = requests.post(f"{HOST_PORT}/model-manager/register",
                             json={
                                 "name": name,
                                 "train_script": train_script,
                                 "infer_script": infer_script
                             })
    if response.status_code != 200:
        raise Exception(
            f"call kasos api /model-manager/register failed, code: {response.status_code}, msg: {response.text}"
        )
    return response.json()


def register_model_result(name: str):
    response = requests.get(
        f"{HOST_PORT}/model-manager/register-result?name={name}")
    if response.status_code != 200:
        raise Exception(
            f"call kasos api /model-manager/register-result failed, code: {response.status_code}, msg: {response.text}"
        )
    return response.json()


def fetch_data(start_time: str, end_time: str, tags: str):
    response = requests.post(f"{HOST_PORT}/data-manager/fetch",
                             json={
                                 "start_time": start_time,
                                 "end_time": end_time,
                                 "tags": tags
                             })
    if response.status_code != 200:
        raise Exception(
            f"call kasos api /data-manager/fetch failed, code: {response.status_code}, msg: {response.text}"
        )
    print(response.headers)
    return response.json()

if __name__ == "__main__":
    # with open("./model_lstm_train.py", "r") as f:
    #     train_script = f.read()
    # with open("./model_lstm_infer.py", "r") as f:
    #     infer_script = f.read()
    # resp = register_model("lstm", train_script, infer_script)
    # print(resp)

    # resp = list_models(0)
    # print(resp)

    # resp = register_service("measure", {"auto_hpa": "on", "service_name": "measure"}, "lstm")
    # print(resp)

    # resp = list_services(0)
    # print(resp)

    # resp = delete_service("measure")
    # print(resp)

    resp = fetch_data("2024-03-21 09:00:00", "2024-03-21 12:00:00", 'auto_hpa="on",service_name="measure"')
    print(resp)
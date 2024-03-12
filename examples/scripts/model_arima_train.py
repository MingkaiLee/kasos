#!/usr/bin/python3
import argparse
from pmdarima import auto_arima
import joblib


def parse_args_train() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="train model")
    parser.add_argument("--new", action="store_true")
    parser.add_argument("-d",
                        "--data",
                        type=str,
                        required=True,
                        help="train data path")
    parser.add_argument("-m",
                        "--model",
                        type=str,
                        required=True,
                        help="train model path")

    return parser.parse_args()


def prepare_data(data_path: str) -> list[float]:
    r = list()
    with open(data_path, "r") as f:
        r = [float(line.strip().split("\t")[1]) for line in f.readlines()]

    return r


# 训练arima模型
def min_aic_train(data, max_p: int = 10, max_q: int = 10):
    model = auto_arima(data,
                       start_p=1,
                       start_q=1,
                       max_p=max_p,
                       max_q=max_q,
                       test='adf',
                       seasonal=False,
                       stepwise=True)

    return model


# 保存模型
def save_model(model, model_path: str):
    joblib.dump(model, model_path)


# 加载模型
def load_model(model_path: str):
    return joblib.load(model_path)


if __name__ == "__main__":
    args = parse_args_train()
    data = prepare_data(args.data)
    if args.new:
        model = min_aic_train(data)
    else:
        model = load_model(args.model)
        model.update(data)
    save_model(model, args.model)

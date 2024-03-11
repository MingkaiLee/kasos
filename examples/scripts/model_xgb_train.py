#!/usr/bin/python3
import argparse
import xgboost as xgb
from sklearn.model_selection import train_test_split
import numpy as np
from datetime import datetime


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


def prepare_data(data_path: str) -> list[tuple[float, float]]:
    r = list()
    with open(data_path, "r") as f:
        tmp = [line.strip().split("\t") for line in f.readlines()]
        for timestamp, value in tmp:
            t = datetime.strptime(timestamp, "%Y-%m-%d %H:%M:%S")
            t_ref = t.replace(hour=0, minute=0, second=0)
            diff_sec = (t - t_ref).total_seconds()
            r.append((diff_sec, float(value)))

    return r


def gen_dataset(data: list[tuple[float, float]], window_size: int=4):
    X, y = [], []
    for i in range(len(data) - window_size - 1):
        # TODO
        X.append(data[i:(i+window_size)])
        y.append(data[i+window_size])
    return np.array(X), np.array(y)


def prepare_xgb_dataset(X, y):
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, shuffle=False)
    dtrain = xgb.DMatrix(X, label=y)
    return dtrain


if __name__ == "__main__":
    args = parse_args_train()
    data = prepare_data(args.data)
    print(gen_dataset(data))
#!/usr/bin/python3
import argparse
from datetime import datetime
import pickle
import xgboost as xgb
import numpy as np


def parse_args_infer() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="infer model")
    parser.add_argument("-t",
                        "--timestamp",
                        type=str,
                        default=datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
    parser.add_argument("-v", "--value", type=float, required=True)
    parser.add_argument("-m", "--model", type=str, required=True)

    return parser.parse_args()


def prepare_data(timestamp: str, value: float):
    t = datetime.strptime(timestamp, "%Y-%m-%d %H:%M:%S")
    t_ref = t.replace(hour=0, minute=0, second=0)
    diff_sec = (t - t_ref).total_seconds()
    vec = list()
    for i in range(3):
        v = diff_sec - (3-i)*15
        if v < 0:
            vec.append(0)
        else:
            vec.append(v)
    vec.append(diff_sec)
    vec.append(value)
    return xgb.DMatrix(np.array([vec]))


def load_model(model_path: str) -> xgb.Booster:
    with open(model_path, "rb") as f:
        return pickle.load(f)

if __name__ == "__main__":
    args = parse_args_infer()
    data = prepare_data(args.timestamp, args.value)
    model = load_model(args.model)
    forecast = model.predict(data)
    print(forecast[0])
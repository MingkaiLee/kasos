#!/usr/bin/python3
import argparse
import xgboost as xgb
import numpy as np
from datetime import datetime
import pickle


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
            t_ref = t.replace(minute=0, second=0)
            diff_sec = (t - t_ref).total_seconds()
            r.append((diff_sec, float(value)))

    return r


def gen_dataset(data: list[tuple[float, float]], window_size: int = 4):
    X, y = [], []
    for i in range(len(data)):
        if i < window_size:
            x = [v[0] for v in data[(i - window_size):]]
            x.extend([v[0] for v in data[:i]])
        else:
            x = [v[0] for v in data[(i - window_size):i]]
        x.append(data[i - 1][1])
        X.append(x)
        y.append(data[i][1])
    return np.array(X), np.array(y)


def xgb_train(model, X, y) -> xgb.Booster:
    data = xgb.DMatrix(X, label=y)
    return xgb.train(params,
                     data,
                     num_boost_round=train_round,
                     xgb_model=model)


params = {
    'max_depth': 6,
    'eta': 0.1,
    'objective': 'count:poisson',
    'eval_metric': 'rmse',
}

train_round = 300

if __name__ == "__main__":
    args = parse_args_train()
    data = prepare_data(args.data)
    X, y = gen_dataset(data)
    pre_model = None
    if not args.new:
        with open(args.model, "rb") as f:
            pre_model = pickle.load(f)
    model = xgb_train(pre_model, X, y)
    with open(args.model, "wb+") as f:
        pickle.dump(model, f)

#!/usr/bin/python3
import argparse
from datetime import datetime
import joblib


def parse_args_infer() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="infer model")
    parser.add_argument("-t",
                        "--timestamp",
                        type=str,
                        default=datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
    parser.add_argument("-v", "--value", type=float, required=True)
    parser.add_argument("-m", "--model", type=str, required=True)

    return parser.parse_args()


def timestamp_2_index(timestamp: str) -> int:
    t = datetime.strptime(timestamp, "%Y-%m-%d %H:%M:%S")
    t_ref = t.replace(minute=0, second=0)
    diff_sec = (t - t_ref).total_seconds()

    return round(diff_sec / 15)


def load_model(model_path: str):
    return joblib.load(model_path)


if __name__ == "__main__":
    args = parse_args_infer()
    idx = timestamp_2_index(args.timestamp)
    model = load_model(args.model)
    forecast = model.predict_in_sample(end=idx)
    print(forecast[-1])

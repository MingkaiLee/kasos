#!/usr/bin/python3
from datetime import datetime, timedelta
import argparse
import random
import numpy as np


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="generate data for models")
    parser.add_argument(
        "-d",
        "--date",
        type=str,
        default=datetime.now().strftime("%Y-%m-%d"),
        help="date of qps time series data",
    )
    parser.add_argument(
        "-b",
        "--base",
        type=int,
        default=100,
        help="base value of max qps",
    )
    parser.add_argument(
        "-o",
        "--output",
        type=str,
        default="data.csv",
        help="output file name",
    )
    return parser.parse_args()


# 每组12个值, 表示上升趋势, 下降趋势, 先升后降, 先降后升, 平稳趋势
element = [[5, 10, 20, 30, 35, 40, 50, 60, 70, 75, 80, 85],
           [90, 80, 70, 65, 60, 50, 40, 35, 25, 15, 10, 5],
           [10, 25, 40, 65, 85, 95, 70, 60, 45, 30, 20, 10],
           [95, 70, 55, 45, 25, 10, 5, 20, 35, 55, 70, 90],
           [40, 55, 45, 60, 50, 55, 55, 45, 60, 50, 45, 60]]


def gen_data(date: str, base: int) -> list[tuple[str, float]]:
    date_start = datetime.strptime(date, "%Y-%m-%d")
    date_end = date_start + timedelta(days=1)
    res = list()
    while date_start < date_end:
        s = random.randint(0, 4)
        choice = element[s]
        for p in choice:
            date_str = date_start.strftime("%Y-%m-%d %H:%M:%S")
            v = float(base) * p / 100 * (1 + rand_factor())
            res.append((date_str, v))
            date_start += timedelta(seconds=15)
    return res


def rand_factor() -> float:
    v = np.random.normal(0, 1)
    while (v > 0.8 or v < -0.8):
        v = np.random.normal(0, 1)
    return v


if __name__ == "__main__":
    args = parse_args()
    data = gen_data(args.date, args.base)
    with open(args.output, "w+") as f:
        for d in data:
            f.write("{:s}\t{:.4f}\n".format(d[0], d[1]))

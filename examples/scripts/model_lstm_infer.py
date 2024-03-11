#!/usr/bin/python3
import argparse
from datetime import datetime


def parse_args_infer() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="infer model")
    parser.add_argument("-t",
                        "--timestamp",
                        type=str,
                        default=datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
    parser.add_argument("-v", "--value", type=float, required=True)
    parser.add_argument("-m", "--model", type=str, required=True)

    return parser.parse_args()
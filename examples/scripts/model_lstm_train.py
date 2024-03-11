#!/usr/bin/python3
import argparse


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
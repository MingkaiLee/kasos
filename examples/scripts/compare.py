#!/usr/bin/python3
import argparse
import numpy as np
from dtaidistance import dtw
from dtaidistance import dtw_visualisation as dtwvis


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="draw time series")
    parser.add_argument("-b",
                        "--observe",
                        type=str,
                        required=True,
                        help="observe data path")
    parser.add_argument("-d",
                        "--data",
                        type=str,
                        required=True,
                        help="infer data path")
    parser.add_argument("-o",
                        "--output",
                        type=str,
                        required=True,
                        help="output image file path")

    return parser.parse_args()


def prepare_data(data_path: str) -> list[float]:
    r = list()
    with open(data_path, "r") as f:
        r = [float(line.strip().split("\t")[1]) for line in f.readlines()]

    return r


start = 4 * 60 * 12
end = 4 * 60 * 12 + 4 * 30

if __name__ == "__main__":
    args = parse_args()
    ob_data = np.array(prepare_data(args.observe)[1:])
    infer_data = np.array(prepare_data(args.data)[:-1])
    distance = dtw.distance(ob_data, infer_data, use_pruning=True)
    print(f"dtw distance: {distance}")
    ob_data_vis = ob_data[start:end]
    infer_data_vis = infer_data[start:end]
    path = dtw.warping_path(ob_data_vis, infer_data_vis)
    dtwvis.plot_warping(ob_data_vis, infer_data_vis, path, args.output)

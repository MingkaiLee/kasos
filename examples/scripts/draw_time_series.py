#!/usr/bin/python3
import seaborn as sns
import pandas as pd
import matplotlib.pyplot as plt
import argparse
import os.path as osp


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="draw time series")
    parser.add_argument("-d",
                        "--data",
                        nargs="*",
                        type=str,
                        required=True,
                        help="data path")
    parser.add_argument("-o",
                        "--output",
                        type=str,
                        required=True,
                        help="output image file path")

    return parser.parse_args()


def read_data(data_path: str) -> pd.DataFrame:
    name = osp.split(data_path)[1].split(".")[0]
    return pd.read_csv(data_path,
                       sep="\t",
                       index_col=0,
                       header=None,
                       names=[name])


if __name__ == "__main__":
    args = parse_args()
    dfs: list[pd.DataFrame] = list()
    for data_path in args.data:
        df = read_data(data_path)
        dfs.append(df)
    data = dfs[0].join(dfs[1:], how="outer")
    data.index = [x.split(" ")[1] for x in data.index]
    plt_data = data.iloc[4 * 60 * 12 - 1:4 * 60 * 12 + 4 * 30]
    sns.set_style("darkgrid")
    sns.lineplot(data=plt_data)
    plt.xticks(rotation=45, fontsize=8)
    for index, label in enumerate(plt.gca().xaxis.get_ticklabels()):
        if index % 8 != 0:
            label.set_visible(False)
    plt.tight_layout()
    plt.legend()
    plt.savefig(args.output)

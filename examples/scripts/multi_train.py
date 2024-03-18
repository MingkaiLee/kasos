#!/usr/bin/python3
import argparse
import subprocess
import os
import os.path as osp


def parse_args_train() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="train model with data under a directory")
    parser.add_argument("-s",
                        "--script",
                        type=str,
                        required=True,
                        help="train script path")
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


if __name__ == "__main__":
    args = parse_args_train()
    file_names = [osp.join(args.data, x) for x in os.listdir(args.data)]
    file_names.sort()
    r = subprocess.run([
        "python3", args.script, "--new", "-d", file_names[0], "-m", args.model
    ])
    print(
        f"Init train done, code: {r.returncode}, msg: {r.stdout}, err: {r.stderr}"
    )
    for i, data_file in enumerate(file_names[1:]):
        r = subprocess.run(
            ["python3", args.script, "-d", data_file, "-m", args.model])
        print(
            f"Train epoch {i} done, code: {r.returncode}, msg: {r.stdout}, err: {r.stderr}"
        )

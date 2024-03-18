#!/usr/bin/python3
import argparse
import subprocess
from datetime import datetime, timedelta


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="infer")
    parser.add_argument("-s",
                        "--script",
                        type=str,
                        required=True,
                        help="infer script path")
    parser.add_argument("-d",
                        "--data",
                        type=str,
                        required=True,
                        help="infer data path")
    parser.add_argument("-m",
                        "--model",
                        type=str,
                        required=True,
                        help="infer model path")
    parser.add_argument("-o",
                        "--output",
                        type=str,
                        required=True,
                        help="infer output path")

    return parser.parse_args()


if __name__ == "__main__":
    args = parse_args()
    with open(args.data, "r") as f:
        data = f.readlines()
    result = []
    for line in data:
        vals = line.strip().split("\t")
        timestamp = vals[0]
        value = vals[1]
        cmd = subprocess.run(["python3", args.script, "-t", timestamp, "-v", value, "-m", args.model], stdout=subprocess.PIPE)
        infer_time = datetime.strptime(timestamp, "%Y-%m-%d %H:%M:%S") + timedelta(seconds=15)
        infer_timestamp = infer_time.strftime("%Y-%m-%d %H:%M:%S")
        infer_value = cmd.stdout.decode("utf-8").strip()
        result.append("{:s}\t{:s}".format(infer_timestamp, infer_value))
        print("{} {} {} {}".format(timestamp, value, infer_timestamp, infer_value))
    with open(args.output, "w+") as f:
        f.writelines(result)
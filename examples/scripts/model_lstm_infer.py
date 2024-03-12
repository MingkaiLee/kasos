#!/usr/bin/python3
import argparse
from datetime import datetime
import torch
import torch.nn as nn
from sklearn.preprocessing import MinMaxScaler
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


class lstm(nn.Module):

    def __init__(self,
                 input_size=5,
                 hidden_layer_size=100,
                 output_size=1,
                 *args,
                 **kwargs) -> None:
        super().__init__(*args, **kwargs)
        self.hidden_layer_size = hidden_layer_size
        self.model = nn.LSTM(input_size, hidden_layer_size)
        self.linear = nn.Linear(hidden_layer_size, output_size)
        self.hidden_cell = (torch.zeros(1, 1, hidden_layer_size),
                            torch.zeros(1, 1, hidden_layer_size))

    def forward(self, input_seq):
        output, _ = self.model(input_seq.view(len(input_seq), 1, -1),
                               self.hidden_cell)
        predictions = self.linear(output.view(len(output), -1))

        return predictions[-1]


def gen_data(timestamp: str, value: float) -> dict:
    t = datetime.strptime(timestamp, "%Y-%m-%d %H:%M:%S")
    t_ref = t.replace(hour=0, minute=0, second=0)
    diff_sec = (t - t_ref).total_seconds()
    t_range = [0, 24 * 60 * 60 - 15]
    for i in range(4):
        v = diff_sec - (3 - i) * 15
        if v < 0:
            t_range.append(24 * 60 * 60 + v)
        else:
            t_range.append(v)
    scalar = MinMaxScaler((-1, 1))
    data = scalar.fit_transform(np.array(t_range).reshape(
        -1, 1)).flatten().tolist()[2:]
    data.append(value / qps_scale)

    return torch.tensor([data])


qps_scale = 100

if __name__ == "__main__":
    args = parse_args_infer()
    X = gen_data(args.timestamp, args.value)
    model = torch.load(args.model)
    model.eval()
    with torch.no_grad():
        pred = model(X)
        print(float(pred[0] * qps_scale))

#!/usr/bin/python3
import argparse
from datetime import datetime
import torch
import torch.nn as nn
import numpy as np
from torch.utils.data import DataLoader, TensorDataset
from sklearn.preprocessing import MinMaxScaler


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
    scalar = MinMaxScaler((-1, 1))
    # 归一化处理
    timestamp = scalar.fit_transform(
        np.array([v[0] for v in data]).reshape(-1, 1)).flatten().tolist()
    # value = scalar.fit_transform(np.array([v[1] for v in data]).reshape(-1, 1)).flatten().tolist()
    value = [v[1] / qps_scale for v in data]
    X, y = [], []
    for i in range(len(data)):
        if i < window_size:
            x = timestamp[(i - window_size):]
            x.extend(timestamp[:i])
        else:
            x = timestamp[(i - window_size):i]
        x.append(value[i - 1])
        X.append(x)
        y.append(value[i])

    return X, y


def gen_data_loader(X, y, batch_size=1):
    dataset = TensorDataset(torch.tensor(X), torch.tensor(y))
    loader = DataLoader(dataset, batch_size, True)

    return loader


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


qps_scale = 10
epochs = 100

if __name__ == "__main__":
    args = parse_args_train()
    data = prepare_data(args.data)
    X, y = gen_dataset(data)
    loader = gen_data_loader(X, y)
    if args.new:
        model = lstm()
    else:
        model = torch.load(args.model)
    loss_fn = nn.MSELoss()
    optimizer = torch.optim.Adam(model.parameters(), lr=0.01)
    for i in range(epochs):
        for s, l in loader:
            optimizer.zero_grad()
            y_pred = model(s)
            loss = loss_fn(y_pred, l)
            loss.backward()
            optimizer.step()
        print(f'epoch: {i:3} loss: {loss.item():10.8f}')
    torch.save(model, args.model)

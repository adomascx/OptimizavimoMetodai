from collections import defaultdict

import matplotlib.pyplot as plt
import numpy as np


DATA_FILE = "plot_data.txt"


def parse_data(path):
    meta = {}
    grid_n = None
    x_vals = None
    y_vals = None
    z = None

    summaries = defaultdict(dict)
    paths = defaultdict(lambda: defaultdict(list))

    with open(path, "r", encoding="utf-8") as f:
        for raw in f:
            line = raw.strip()
            if not line or line.startswith("#"):
                continue

            parts = line.split()
            tag = parts[0]

            if tag == "META":
                meta[parts[1]] = float(parts[2])
            elif tag == "GRID_N":
                grid_n = int(parts[1])
                x_vals = np.zeros(grid_n)
                y_vals = np.zeros(grid_n)
                z = np.zeros((grid_n, grid_n))
            elif tag == "SUMMARY":
                algo = parts[1]
                start = parts[2]
                summaries[algo][start] = {
                    "end_x": float(parts[3]),
                    "end_y": float(parts[4]),
                    "dev": float(parts[5]),
                    "iter": int(parts[6]),
                    "fn": int(parts[7]),
                    "grad": int(parts[8]),
                }
            elif tag == "PATH":
                algo = parts[1]
                start = parts[2]
                idx = int(parts[3])
                x = float(parts[4])
                y = float(parts[5])
                paths[algo][start].append((idx, x, y))
            elif tag == "GRID":
                if x_vals is None or y_vals is None or z is None:
                    raise ValueError("GRID line found before GRID_N declaration")
                i = int(parts[1])
                j = int(parts[2])
                x = float(parts[3])
                y = float(parts[4])
                gm = float(parts[5])
                x_vals[i] = x
                y_vals[j] = y
                z[j, i] = gm

    if grid_n is None or x_vals is None or y_vals is None or z is None:
        raise ValueError("GRID_N or GRID data not found in plot_data.txt")

    for algo in paths:
        for start in paths[algo]:
            paths[algo][start].sort(key=lambda t: t[0])
            paths[algo][start] = [(t[1], t[2]) for t in paths[algo][start]]

    return meta, x_vals, y_vals, z, summaries, paths


def plot_algorithm(algo, out_file, x_vals, y_vals, z, summaries, paths):
    fig, ax = plt.subplots(figsize=(9, 7), dpi=150)

    im = ax.imshow(
        z,
        origin="lower",
        extent=(float(x_vals[0]), float(x_vals[-1]), float(y_vals[0]), float(y_vals[-1])),
        cmap="viridis",
        aspect="equal",
    )
    cbar = fig.colorbar(im, ax=ax)
    cbar.set_label("|grad f(x,y)|")

    color_map = {"X0": "tab:red", "X1": "tab:green", "Xm": "tab:blue"}

    for start in ["X0", "X1", "Xm"]:
        run = paths.get(algo, {}).get(start, [])
        if not run:
            continue
        arr = np.array(run)

        summary = summaries.get(algo, {}).get(start)
        if summary:
            label = f"{start} (dev={summary['dev']:.3e})"
        else:
            label = start

        ax.plot(
            arr[:, 0],
            arr[:, 1],
            color=color_map[start],
            linewidth=2.5,
            label=label,
        )
        ax.scatter(arr[0, 0], arr[0, 1], color=color_map[start], s=24, marker="o")
        ax.scatter(arr[-1, 0], arr[-1, 1], color=color_map[start], s=36, marker="x")

    ax.set_xlim(0.0, 1.0)
    ax.set_ylim(0.0, 1.0)
    ax.set_xlabel("x")
    ax.set_ylabel("y")
    ax.set_title(f"{algo}: optimization paths on |grad f(x,y)|")
    ax.grid(alpha=0.25)
    ax.legend(loc="upper right")

    fig.tight_layout()
    fig.savefig(out_file)
    plt.close(fig)


def main():
    _, x_vals, y_vals, z, summaries, paths = parse_data(DATA_FILE)

    plot_algorithm("GD", "gd.png", x_vals, y_vals, z, summaries, paths)
    plot_algorithm("SD", "sd.png", x_vals, y_vals, z, summaries, paths)
    plot_algorithm("SIMPLEX", "simplex.png", x_vals, y_vals, z, summaries, paths)

    print("Generated: gd.png, sd.png, simplex.png")


if __name__ == "__main__":
    main()

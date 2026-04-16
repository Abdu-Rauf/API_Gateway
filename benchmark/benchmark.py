import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

# 1. Load and clean the data
df = pd.read_csv('../results.csv')

def parse_latency(val):
    if isinstance(val, str):
        if 'ms' in val:
            return float(val.replace('ms', ''))
        elif 'us' in val:
            return float(val.replace('us', '')) / 1000.0
    return float(val)

df['Latency_ms'] = df['Latency'].apply(parse_latency)
df['TotalErrors'] = df['Timeouts'] + df['ConnErr']

node_data = df[df['Server'] == 'NodeJS']
go_data = df[df['Server'] == 'Go']

x_labels = ['10', '100', '1000', '5000']
x = np.arange(len(x_labels))

# Plot 1: Connections vs RPS
plt.figure(figsize=(10, 6))
plt.plot(x, node_data['ReqPerSec'], marker='o', label='Node.js', linewidth=2)
plt.plot(x, go_data['ReqPerSec'], marker='s', label='Go', linewidth=2)
plt.xticks(x, x_labels)
plt.xlabel('Concurrent Connections')
plt.ylabel('Requests Per Second (RPS)')
plt.title('Throughput: Connections vs Req/Sec')
plt.legend()
plt.grid(True)
plt.savefig('throughput.png')

# Plot 2: Connections vs Latency
plt.figure(figsize=(10, 6))
plt.plot(x, node_data['Latency_ms'], marker='o', label='Node.js', linewidth=2)
plt.plot(x, go_data['Latency_ms'], marker='s', label='Go', linewidth=2)
plt.xticks(x, x_labels)
plt.xlabel('Concurrent Connections')
plt.ylabel('Latency (ms)')
plt.title('Responsiveness: Connections vs Latency')
plt.legend()
plt.grid(True)
plt.savefig('latency.png')

# Plot 3: Error Rate
plt.figure(figsize=(10, 6))
width = 0.35
plt.bar(x - width/2, node_data['TotalErrors'], width, label='Node.js Errors')
plt.bar(x + width/2, go_data['TotalErrors'], width, label='Go Errors')
plt.xticks(x, x_labels)
plt.xlabel('Concurrent Connections')
plt.ylabel('Total Errors (Timeouts + ConnErr)')
plt.title('Stability: Error Rate under Load')
plt.legend()
plt.grid(True, axis='y')
plt.savefig('errors.png')

print("Graphs successfully generated and saved as PNG files!")
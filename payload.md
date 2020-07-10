Current metric payload:
[
    {"timestamp":1582898206000,"route":"/metric","uuid":"a2582b22-5e0e-11ea-bc5d-acde48001122","metric":"586b95c6","value":0.6046602879796196},
    ...
]

10_000 metrics/gateway = metric message size: 1,392,753, compressed size: 160,408, ratio: 0.12 

156 KB /sec

156 * 86400 KB /day

Cloud IoT upload limit 512 KB/s 
512 - 156 = 356 KB/s  (potential headroom)

it would take (156 * 86400) /356 /3600 = 10.5 hours to upload 24 hours of history


New metric payload encoded in protobuf:
{
    "timestamp":1582898206000,
    "metrics":[
        {"name":"586b95c6","value":0.6046602879796196, "timestamp":1582898206000},
        ...
    ]
}

10,000 metric message size (encoded in protobuf, uncompressed): 300,700
- or -
10,000 metric message size (encoded in protobuf, compressed): 165,500
162 KB /sec

162 * 86400 KB /day

Cloud IoT upload limit 512 KB/s 
512 - 162 = 350 KB/s  (potential headroom)
it would take (162 * 86400) /350 /3600 = 11.1 hours to upload 24 hours of history

But, if instead of one device, there were 10 devices, each with 1000 metrics.

1,000 metric message size (encoded in protobuf, compressed): 16,550
16.2 KB /sec

16.2 * 86400 KB /day

Cloud IoT upload limit 512 KB/s 
512 - 16.2 = 495.8 KB/s  (potential headroom)
it would take (16.2 * 86400) /495.8 /3600 = 47 minutes to upload 24 hours of history

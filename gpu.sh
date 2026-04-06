#!/bin/sh
apt-get update && apt-get install -y \
            ocl-icd-libopencl1 \
            opencl-headers \
            clinfo \
            && rm -rf /var/lib/apt/lists/*

mkdir -p /etc/OpenCL/vendors && \
            echo "libnvidia-opencl.so.1" > /etc/OpenCL/vendors/nvidia.icd

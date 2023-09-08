# window-post-bench

## Setup

```shell
sudo apt install -y mesa-opencl-icd ocl-icd-opencl-dev gcc git bzr jq pkg-config curl clang build-essential hwloc libhwloc-dev wget

# install golang && rust
# more details: https://lotus.filecoin.io/lotus/install/linux/
```

## Build 

```shell
git submodule update --init --recursive

make -C extern/filecoin-ffi/

go build .
```

## Run

```shell
./window-post-bench 
```
# window-post-bench

## Setup

```shell
sudo apt install -y hwloc
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
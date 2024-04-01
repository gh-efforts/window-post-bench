package main

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"github.com/docker/go-units"
	ffi "github.com/filecoin-project/filecoin-ffi"
	"github.com/filecoin-project/go-paramfetch"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
	lcli "github.com/filecoin-project/lotus/cli"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"math/big"
	"os"
	"time"
)

var log = logging.Logger("window-post-bench")

//go:embed proof.bin
var vp []byte

func main() {
	logging.SetDebugLogging()
	log.Info("Starting window-post-bench")

	app := &cli.App{
		Name:                      "window-post-bench",
		Usage:                     "Benchmark performance of lotus window post on your hardware",
		Version:                   build.BuildVersion,
		DisableSliceFlagSeparator: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "sector-size",
				Value: "32GiB",
			},
		},
		Action: func(c *cli.Context) error {
			sectorSizeInt, err := units.RAMInBytes(c.String("sector-size"))
			if err != nil {
				return err
			}
			sectorSize := abi.SectorSize(sectorSizeInt)

			ctx := lcli.ReqContext(c)

			if err := paramfetch.GetParams(ctx, build.ParametersJSON(), build.SrsJSON(), uint64(sectorSize)); err != nil {
				return xerrors.Errorf("get params: %w", err)
			}

			wpt, err := spt(sectorSize).RegisteredWindowPoStProof()
			if err != nil {
				return err
			}

			challenge := time.Now()
			var rand [32]byte // all zero
			proof, err := ffi.GenerateSinglePartitionWindowPoStWithVanilla(wpt, abi.ActorID(100), rand[:], [][]byte{vp}, 0)
			if err != nil {
				return xerrors.Errorf("generate post: %w", err)
			}
			end := time.Now()

			fmt.Printf("Proof %s (%s)\n", end.Sub(challenge), bps(sectorSize, 1, end.Sub(challenge)))
			fmt.Println(base64.StdEncoding.EncodeToString(proof.ProofBytes))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Warnf("%+v", err)
		return
	}
}

func spt(ssize abi.SectorSize) abi.RegisteredSealProof {
	spt, err := miner.SealProofTypeFromSectorSize(ssize, build.TestNetworkVersion, false)
	if err != nil {
		panic(err)
	}

	return spt
}

func bps(sectorSize abi.SectorSize, sectorNum int, d time.Duration) string {
	bdata := new(big.Int).SetUint64(uint64(sectorSize))
	bdata = bdata.Mul(bdata, big.NewInt(int64(sectorNum)))
	bdata = bdata.Mul(bdata, big.NewInt(time.Second.Nanoseconds()))
	bps := bdata.Div(bdata, big.NewInt(d.Nanoseconds()))
	return types.SizeStr(types.BigInt{Int: bps}) + "/s"
}

package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	aw "github.com/deanishe/awgo"
	cidpkg "github.com/ipfs/go-cid"
	mb "github.com/multiformats/go-multibase"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
)

// wf is the workflows main API.
var wf *aw.Workflow

func init() {
	wf = aw.New()
}

// run contains the workflow logic.
func run() {

	defer wf.SendFeedback()

	cidstr := wf.Args()[0]

	cid, err := cidpkg.Decode(cidstr)
	if err != nil {
		wf.NewItem("Invalid CID").Subtitle(err.Error())
		return
	}

	enc, err := cidpkg.ExtractEncoding(cidstr)
	if err != nil {
		wf.NewItem("Invalid encoding").Subtitle(err.Error())
		return
	}

	dmh, err := mh.Decode(cid.Hash())
	if err != nil {
		wf.NewItem("Invalid multihash").Subtitle(err.Error())
		return
	}

	encstr := mb.EncodingToStr[enc]
	verstr := fmt.Sprintf("cidv%d", cid.Prefix().Version)
	codstr := mc.Code(cid.Prefix().Codec).String()
	mhtype := mc.Code(dmh.Code).String()
	mhsize := fmt.Sprintf("%d", dmh.Length*8)
	digest := strings.ToUpper(hex.EncodeToString(dmh.Digest))

	wf.NewItem(encstr).Subtitle(fmt.Sprintf("Encoding (%s)", string(rune(enc)))).Arg(encstr).Valid(true)
	wf.NewItem(verstr).Subtitle("Version").Arg(verstr).Valid(true)
	wf.NewItem(codstr).Subtitle(fmt.Sprintf("Multicodec (0x%x)", cid.Prefix().Codec)).Arg(codstr).Valid(true)
	wf.NewItem(mhtype).Subtitle(fmt.Sprintf("Multihash Type (0x%x)", dmh.Code)).Arg(mhtype).Valid(true)
	wf.NewItem(mhsize).Subtitle("Multihash Size (bits)").Arg(mhsize).Valid(true)
	wf.NewItem(digest).Subtitle("Digest in Hex").Arg(digest).Valid(true)
}

func main() {
	wf.Run(run)
}

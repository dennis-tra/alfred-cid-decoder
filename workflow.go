package main

// Package is called aw
import (
	"encoding/hex"
	"fmt"
	"strings"

	aw "github.com/deanishe/awgo"
	cid "github.com/ipfs/go-cid"
	mb "github.com/multiformats/go-multibase"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
)

// Workflow is the main API
var wf *aw.Workflow

func init() {
	// Create a new Workflow using default settings.
	// Critical settings are provided by Alfred via environment variables,
	// so this *will* die in flames if not run in an Alfred-like environment.
	wf = aw.New()
}

// Your workflow starts here
func run() {

	args := wf.Args()
	if len(args) <= 0 {
		return
	}
	id := args[0]

	c, err := cid.Decode(id)
	if err != nil {
		it := wf.NewItem("Invalid CID")
		it.Subtitle(id)
		wf.SendFeedback()
		return
	}

	e, _ := cid.ExtractEncoding(id)

	wf.NewItem(mb.EncodingToStr[e]).Subtitle("Encoding")
	wf.NewItem(fmt.Sprintf("cidv%d", c.Prefix().Version)).Subtitle("CID Version")
	wf.NewItem(mc.Code(c.Prefix().Codec).String()).Subtitle(fmt.Sprintf("Multicodec (0x%x)", c.Prefix().Codec))
	dmh, _ := mh.Decode(c.Hash())
	wf.NewItem(mc.Code(dmh.Code).String()).Subtitle(fmt.Sprintf("Multihash Type (0x%x)", dmh.Code))
	wf.NewItem(fmt.Sprintf("%d", dmh.Length*8)).Subtitle("Multihash Size (bits)")
	wf.NewItem(strings.ToUpper(hex.EncodeToString(dmh.Digest))).Subtitle("Digest in Hex")

	// Send results to Alfred
	wf.SendFeedback()
}

func main() {
	// Wrap your entry point with Run() to catch and log panics and
	// show an error in Alfred instead of silently dying
	wf.Run(run)
}

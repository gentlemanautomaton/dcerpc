package dcerpc

import "context"

// Client is a DCE / RPC client that is capable of making procedure calls to
// one or more remote servers.
//
// Client implements the client side of the protocol described in the
// "DCE 1.1: Remote Procedure Call" technical standard.
type Client struct {
}

// Invoke will run the requested remote procedure.
//
// TODO: Figure out a good way to pass binding and interface information in.
//
// TODO: Figure out a good way to pass the RPC parameters in.
//
// TODO: Return some sort of call struct (or interface)?
func (c Client) Invoke(ctx context.Context) {

}

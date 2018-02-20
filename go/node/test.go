// Copyright 2018 The Keep Authors.  See LICENSE.md for details.
package node

import (
	"context"
	"log"
	host "github.com/libp2p/go-libp2p-host"
)

//Test will invoke the test handler and console log a message.
func Test(ha host.Host) {

	peers := ha.Peerstore().Peers()

	for _, p := range peers {
		if ha.ID().String() != p.String() {
			stream, err := ha.NewStream(context.Background(), p, "/test/1.0.0")
			if err != nil {
				log.Println("ERROR creating new stream", err)
				continue
			}

			for _, addr := range ha.Addrs() {
				if addr.String() != "" {
					_, err = stream.Write([]byte(addr.String() + "/ipfs/" + ha.ID().Pretty() + "\n"))
					if err != nil {
						log.Println("ERROR writing to stream", err)
					}
				}
			}
		}
	}
}

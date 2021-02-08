package chat

import (
	"github.com/hashicorp/memberlist"
	"log"
	"strings"
)
// We use a simple broadcast implementation in which items are never invalidated by others.
type simpleBroadcast []byte

func (b simpleBroadcast) Message() []byte                       { return []byte(b) }
func (b simpleBroadcast) Invalidates(memberlist.Broadcast) bool { return false }
func (b simpleBroadcast) Finished()                             {}

var broadcasts *memberlist.TransmitLimitedQueue

type Delegate struct{
	Messages *chan string
}

func (d *Delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

func (d *Delegate) NotifyMsg(b []byte) {
	*d.Messages <- string(b)
	log.Println("received from gossip: ",string(b))
}

func (d *Delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return broadcasts.GetBroadcasts(overhead, limit)
}
func (d *Delegate) LocalState(join bool) []byte {
	return []byte{'a'}
}
func (d *Delegate) MergeRemoteState(buf []byte, join bool) {

}


func Gossip(defaultConfig *memberlist.Config,clusterAddr *string)  {
	list, err := memberlist.Create(defaultConfig)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	// Join an existing cluster by specifying at least one known member.
	_, err = list.Join(	strings.Split(*clusterAddr,","))
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}

	// Ask for members of the cluster
	for _, member := range list.Members() {
		log.Printf("Member: %s %s\n", member.Name, member.Addr)
	}
	broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return list.NumMembers()
		},
		RetransmitMult: 0,
	}

}
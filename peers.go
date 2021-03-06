package cache

import pb "mini-cache/cachepb"

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer. 对应 HTTP 客户端
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}

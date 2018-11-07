package server

import "google.golang.org/grpc"

const defaultMethod = "unknown"

type RPCCallInfo struct {
	FullMethod string
}

func ServerInfoFromGrpc(gsi interface{}) *RPCCallInfo {
	si := &RPCCallInfo{FullMethod: defaultMethod}
	switch st := gsi.(type) {
	case grpc.UnaryServerInfo:
		si.FullMethod = st.FullMethod
	case grpc.StreamServerInfo:
		si.FullMethod = st.FullMethod
	}
	return si
}

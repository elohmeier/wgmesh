package meshservice

import (
	context "context"
	"errors"
	"fmt"
	"math/rand"
	"net"

	wgwrapper "github.com/aschmidt75/go-wg-wrapper/pkg/wgwrapper"
	log "github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

// BeginJoin allows other nodes to join by sending a JoinRequest
func (ms *MeshService) BeginJoin(ctx context.Context, req *JoinRequest) (*JoinResponse, error) {

	log.WithField("req", req).Trace("Got join request")

	// choose a random ip adress from the adress pool of this node
	// which has not been used before
	mip, _ := newIPInNet(ms.CIDRRange)

	targetWGIP := net.IPNet{
		mip,
		net.CIDRMask(32, 32),
	}

	// take public key and endpoint, add as peer to own wireguard interface
	p := wgwrapper.WireguardPeer{
		RemoteEndpointIP: req.EndpointIP,
		ListenPort:       int(req.EndpointPort),
		Pubkey:           req.Pubkey,
		AllowedIPs: []net.IPNet{
			targetWGIP,
		},
		Psk: nil,
	}
	log.WithField("peer", p).Trace("Adding peer")

	wg := wgwrapper.New()

	ok, err := wg.AddPeer(ms.WireguardInterface, p)
	if err != nil {
		return &JoinResponse{
			Result:       JoinResponse_ERROR,
			ErrorMessage: "Unable to add peer",
			JoinerMeshIP: "",
		}, nil
	}
	if !ok && err == nil {
		return &JoinResponse{
			Result:       JoinResponse_ERROR,
			ErrorMessage: "Peer already present",
			JoinerMeshIP: "",
		}, nil
	}

	return &JoinResponse{
		Result:       JoinResponse_OK,
		ErrorMessage: "",
		JoinerMeshIP: mip.String(),
		MeshCidr:     ms.CIDRRange.String(),
	}, nil
}

// Peers serves a list of all current peers, starting with this node
func (ms *MeshService) Peers(e *Empty, stream Mesh_PeersServer) error {
	err := stream.Send(&Peer{
		Pubkey:       ms.WireguardPubKey,
		EndpointIP:   ms.WireguardListenIP.String(),
		EndpointPort: int32(ms.WireguardListenPort),
		MeshIP:       ms.MeshIP.IP.String(),
	})

	return err
}

func newIPInNet(ipnet net.IPNet) (net.IP, error) {

	ipmask := ipnet.Mask
	log.WithField("ipmask", ipmask).Trace("dump")
	log.WithField("ip", ipnet.IP).Trace("dump")

	var newIP [4]byte
	if len(ipnet.IP) == 4 {
		newIP = [4]byte{
			(byte(rand.Intn(250)+2) & ^ipmask[0]) + ipnet.IP[0],
			(byte(rand.Intn(250)) & ^ipmask[1]) + ipnet.IP[1],
			(byte(rand.Intn(250)) & ^ipmask[2]) + ipnet.IP[2],
			(byte(rand.Intn(250)+1) & ^ipmask[3]) + ipnet.IP[3],
		}
	}
	if len(ipnet.IP) == 16 {
		newIP = [4]byte{
			(byte(rand.Intn(250)+2) & ^ipmask[0]) + ipnet.IP[12],
			(byte(rand.Intn(250)) & ^ipmask[1]) + ipnet.IP[13],
			(byte(rand.Intn(250)) & ^ipmask[2]) + ipnet.IP[14],
			(byte(rand.Intn(250)+1) & ^ipmask[3]) + ipnet.IP[15],
		}
	}
	log.WithField("newIP", newIP).Trace("newIPInNet.dump")

	return net.IPv4(newIP[0], newIP[1], newIP[2], newIP[3]), nil
}

// StartGrpcService ..
func (ms *MeshService) StartGrpcService() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ms.GrpcBindAddr, ms.GrpcBindPort))
	if err != nil {
		log.Errorf("failed to listen: %v", err)
		return errors.New("unable to start grpc mesh service")
	}

	ms.grpcServer = grpc.NewServer()
	RegisterMeshServer(ms.grpcServer, ms)
	if err := ms.grpcServer.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
		return errors.New("unable to start grpc mesh service")
	}

	return nil
}

// StopGrpcService ...
func (ms *MeshService) StopGrpcService() {

	log.Info("Stopping gRPC mesh service")
	ms.grpcServer.GracefulStop()
	log.Info("Stopped gRPC mesh service")
}
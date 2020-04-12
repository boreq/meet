package meet

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/boreq/errors"

	"github.com/pion/webrtc/v2"
)

const (
	rtcpPLIInterval = time.Second * 3
)

type WebRTCMeting struct {
	members []*Member
	mutex   sync.RWMutex
}

func NewWebRTCMeting() *WebRTCMeting {
	return &WebRTCMeting{}
}

func (m *WebRTCMeting) ReceivedSessionDescriptionProtocol(sdp string) (*Member, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	offer, err := m.decodeSessionDescriptionProtocol(sdp)
	if err != nil {
		return nil, errors.Wrap(err, "decoding session description protocol error")
	}

	member, err := NewMember(offer)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a member")
	}

	m.members = append(m.members, member)

	go func() {
		for {
			track, ok := <-member.Tracks()
			if !ok {
				return
			}
			m.trackReceived(track)
		}
	}()

	//// Since we are answering use PayloadTypes declared by offerer
	//mediaEngine := webrtc.MediaEngine{}
	//if err := mediaEngine.PopulateFromSDP(offer); err != nil {
	//	return errors.Wrap(err, "failed to populate media engine")
	//}
	//
	//// Create the API object with the MediaEngine
	//api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))
	//
	//peerConnectionConfig := webrtc.Configuration{
	//	ICEServers: []webrtc.ICEServer{
	//		{
	//			URLs: []string{"stun:stun.l.google.com:19302"},
	//		},
	//	},
	//}
	//
	//// Create a new RTCPeerConnection
	//peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
	//if err != nil {
	//	return errors.Wrap(err, "could not create a peer connection")
	//}
	//
	//// Allow us to receive 1 video track
	//if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
	//	return errors.Wrap(err, "could not add a transceiver")
	//}
	//
	//localTrackChan := make(chan *webrtc.Track)
	//
	//// Set a handler for when a new remote track starts, this just distributes all our packets
	//// to connected peers
	//peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
	//	// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
	//	// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
	//	go func() {
	//		ticker := time.NewTicker(rtcpPLIInterval)
	//		for range ticker.C {
	//			if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
	//				fmt.Println(rtcpSendErr)
	//			}
	//		}
	//	}()
	//
	//	// Create a local track, all our SFU clients will be fed via this track
	//	localTrack, newTrackErr := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
	//	if newTrackErr != nil {
	//		panic(newTrackErr)
	//	}
	//	localTrackChan <- localTrack
	//
	//	rtpBuf := make([]byte, 1400)
	//	for {
	//		i, readErr := remoteTrack.Read(rtpBuf)
	//		if readErr != nil {
	//			panic(readErr)
	//		}
	//
	//		// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
	//		if _, err = localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
	//			panic(err)
	//		}
	//	}
	//})
	//
	//// Set the remote SessionDescription
	//err = peerConnection.SetRemoteDescription(offer)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Create answer
	//answer, err := peerConnection.CreateAnswer(nil)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Sets the LocalDescription, and starts our UDP listeners
	//err = peerConnection.SetLocalDescription(answer)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Get the LocalDescription and take it to base64 so we can paste in browser
	//fmt.Println(Encode(answer))
	//
	//defer peerConnection.Close()

	return member, nil

	//localTrack := <-localTrackChan
	//for {
	//	fmt.Println("")
	//	fmt.Println("Curl an base64 SDP to start sendonly peer connection")
	//
	//	recvOnlyOffer := webrtc.SessionDescription{}
	//	Decode(<-sdpChan, &recvOnlyOffer)
	//
	//	// Create a new PeerConnection
	//	peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	_, err = peerConnection.AddTrack(localTrack)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	// Set the remote SessionDescription
	//	err = peerConnection.SetRemoteDescription(recvOnlyOffer)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	// Create answer
	//	answer, err := peerConnection.CreateAnswer(nil)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	// Sets the LocalDescription, and starts our UDP listeners
	//	err = peerConnection.SetLocalDescription(answer)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	// Get the LocalDescription and take it to base64 so we can paste in browser
	//	fmt.Println(Encode(answer))
	//}
}

func (m *WebRTCMeting) decodeSessionDescriptionProtocol(sdp string) (webrtc.SessionDescription, error) {
	offer := webrtc.SessionDescription{}
	if err := Decode(sdp, &offer); err != nil {
		return webrtc.SessionDescription{}, errors.Wrap(err, "decode failed")
	}
	return offer, nil
}

func (m *WebRTCMeting) trackReceived(track *webrtc.Track) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	fmt.Println(track)
}

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func Encode(obj webrtc.SessionDescription) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", errors.Wrap(err, "json marshal failed")
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func Decode(in string, obj *webrtc.SessionDescription) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return errors.Wrap(err, "base64 decoding failed")
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		return errors.Wrap(err, "json unmarshal failed")
	}

	return nil
}

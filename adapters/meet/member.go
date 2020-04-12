package meet

import (
	"fmt"
	"io"
	"time"

	"github.com/boreq/errors"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
)

type Member struct {
	answer    webrtc.SessionDescription
	trackChan chan *webrtc.Track
}

func NewMember(sessionDescription webrtc.SessionDescription) (*Member, error) {
	trackChan := make(chan *webrtc.Track)

	// Since we are answering use PayloadTypes declared by offerer
	mediaEngine := webrtc.MediaEngine{}
	if err := mediaEngine.PopulateFromSDP(sessionDescription); err != nil {
		return nil, errors.Wrap(err, "failed to populate media engine")
	}

	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a peer connection")
	}

	// Allow us to receive 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		return nil, errors.Wrap(err, "could not add a transceiver")
	}

	// Set a handler for when a new remote track starts, this just distributes all our packets
	// to connected peers
	peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
		go func() {
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		// Create a local track, all our SFU clients will be fed via this track
		localTrack, newTrackErr := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
		if newTrackErr != nil {
			panic(newTrackErr)
		}
		trackChan <- localTrack

		rtpBuf := make([]byte, 1400)
		for {
			i, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				panic(readErr)
			}

			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
			if _, err = localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
				panic(err)
			}
		}
	})

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(sessionDescription)
	if err != nil {
		panic(err)
	}

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		fmt.Println("state changed", state)
	})

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	return &Member{
		answer:    answer,
		trackChan: trackChan,
	}, nil
}

func (m *Member) Tracks() <-chan *webrtc.Track {
	return m.trackChan
}

func (m *Member) Answer() (string, error) {
	return Encode(m.answer)
}

//func zip(in []byte) []byte {
//	var b bytes.Buffer
//	gz := gzip.NewWriter(&b)
//	_, err := gz.Write(in)
//	if err != nil {
//		panic(err)
//	}
//	err = gz.Flush()
//	if err != nil {
//		panic(err)
//	}
//	err = gz.Close()
//	if err != nil {
//		panic(err)
//	}
//	return b.Bytes()
//}
//
//func unzip(in []byte) []byte {
//	var b bytes.Buffer
//	_, err := b.Write(in)
//	if err != nil {
//		panic(err)
//	}
//	r, err := gzip.NewReader(&b)
//	if err != nil {
//		panic(err)
//	}
//	res, err := ioutil.ReadAll(r)
//	if err != nil {
//		panic(err)
//	}
//	return res
//}

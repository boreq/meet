package webrtc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/boreq/errors"
	"github.com/boreq/meet/domain"
	"github.com/boreq/meet/internal/logging"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
)

type WebRTCParticipant struct {
	participant    *domain.Participant
	trackChan      chan *webrtc.Track
	peerConnection *webrtc.PeerConnection
	log            logging.Logger
}

func NewWebRTCParticipant(participant *domain.Participant) (*WebRTCParticipant, error) {
	return &WebRTCParticipant{
		participant: participant,
		trackChan:   make(chan *webrtc.Track),
		log:         logging.New(fmt.Sprintf("WebRTC member %s", participant.UUID())),
	}, nil
}

func (p *WebRTCParticipant) Connect(remoteDescription RemoteDescription) error {
	sessionDescription, err := decodeSdp(remoteDescription.String())
	if err != nil {
		return errors.Wrap(err, "failed to decode sdp description")
	}

	// Since we are answering use PayloadTypes declared by offerer
	mediaEngine := webrtc.MediaEngine{}
	if err := mediaEngine.PopulateFromSDP(sessionDescription); err != nil {
		return errors.Wrap(err, "failed to populate media engine")
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
		return errors.Wrap(err, "could not create a peer connection")
	}

	// Allow us to receive 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		return errors.Wrap(err, "could not add a transceiver")
	}

	// Set a handler for when a new remote track starts, this just distributes all our packets
	// to connected peers
	peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		p.log.Debug("received a track", "id", remoteTrack.ID(), "label", remoteTrack.Label())

		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
		go func() {
			// todo terminate this somehow
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		// Create a local track, all our SFU clients will be fed via this track
		localTrack, err := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
		if err != nil {
			panic(err)
		}

		p.trackChan <- localTrack

		rtpBuf := make([]byte, 1400)
		for {
			i, err := remoteTrack.Read(rtpBuf)
			if err != nil {
				panic(err)
			}

			p.log.Debug("read from remote track", "i", i)

			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
			if _, err = localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
				panic(err)
			}
		}
	})

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(sessionDescription)
	if err != nil {
		return errors.Wrap(err, "could not set the remote description")
	}

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		p.log.Debug("connection state change", "state", state)
	})

	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		p.log.Debug("ice candidate received", "candidate", candidate)
	})

	peerConnection.OnSignalingStateChange(func(state webrtc.SignalingState) {
		p.log.Debug("signaling state changed", "state", state)
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

}

func (m *WebRTCParticipant) AddTrack(track *webrtc.Track) error {
	m.log.Debug("adding a track", "id", track.ID(), "label", track.Label())

	_, err := m.peerConnection.AddTrack(track)
	if err != nil {
		return errors.Wrap(err, "could not add a track")
	}

	return nil
}

func (m *WebRTCParticipant) UUID() domain.ParticipantUUID {
	return m.uuid
}

func (m *WebRTCParticipant) String() string {
	return fmt.Sprintf("member (%s)", m.uuid)
}

func (m *WebRTCParticipant) Tracks() <-chan *webrtc.Track {
	return m.trackChan
}

func (m *WebRTCParticipant) Answer() (string, error) {
	return encodeSdp(m.answer)
}

func (m *WebRTCParticipant) OnRemoteDescription(sessionDescription RemoteDescription) {
}

func encodeSdp(obj webrtc.SessionDescription) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", errors.Wrap(err, "json marshal failed")
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func decodeSdp(sdp string) (webrtc.SessionDescription, error) {
	var offer webrtc.SessionDescription

	b, err := base64.StdEncoding.DecodeString(sdp)
	if err != nil {
		return offer, errors.Wrap(err, "base64 decoding failed")
	}

	if err = json.Unmarshal(b, &offer); err != nil {
		return offer, errors.Wrap(err, "json unmarshal failed")
	}

	return offer, nil
}

package meet

import (
	"fmt"
	"io"
	"time"

	"github.com/boreq/meet/internal/logging"

	"github.com/boreq/errors"
	"github.com/boreq/meet/domain"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
)

type Member struct {
	uuid           domain.ParticipantUUID
	answer         webrtc.SessionDescription
	trackChan      chan *webrtc.Track
	peerConnection *webrtc.PeerConnection
	log            logging.Logger
}

func NewMember(uuid domain.ParticipantUUID, sessionDescription webrtc.SessionDescription) (*Member, error) {
	trackChan := make(chan *webrtc.Track)
	log := logging.New("member " + uuid.String())

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
		log.Debug("received a track", "id", remoteTrack.ID(), "label", remoteTrack.Label())

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
		localTrack, err := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
		if err != nil {
			panic(err)
		}
		trackChan <- localTrack

		rtpBuf := make([]byte, 1400)
		for {
			i, err := remoteTrack.Read(rtpBuf)
			if err != nil {
				panic(err)
			}

			log.Debug("read from remote track", "i", i)

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
		uuid:           uuid,
		log:            log,
		peerConnection: peerConnection,
		answer:         answer,
		trackChan:      trackChan,
	}, nil
}

func (m *Member) AddTrack(track *webrtc.Track) error {
	m.log.Debug("adding a track", "id", track.ID(), "label", track.Label())

	_, err := m.peerConnection.AddTrack(track)
	if err != nil {
		return errors.Wrap(err, "could not add a track")
	}

	return nil
}

func (m *Member) UUID() domain.ParticipantUUID {
	return m.uuid
}

func (m *Member) String() string {
	return fmt.Sprintf("member (%s)", m.uuid)
}

func (m *Member) Tracks() <-chan *webrtc.Track {
	return m.trackChan
}

func (m *Member) Answer() (string, error) {
	return encodeSdp(m.answer)
}

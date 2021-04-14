package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-learn/p2p-video-conference/server/utils"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
)

const (
	rtcpPLIInterval = time.Second * 3
)

type Sdp struct {
	Sdp string
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	// sender to channel of track
	peerConnectionMap := make(map[string]chan *webrtc.Track)

	m := webrtc.MediaEngine{}

	// Setup the codecs you want to use.
	// Only support VP9(video compression), this makes our proxying code simpler
	m.RegisterCodec(webrtc.NewRTPVP9Codec(webrtc.DefaultPayloadTypeVP9, 90000))

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	peerConnectionConfig := webrtc.Configuration{
		ICEServers:           []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}},
		ICETransportPolicy:   0,
		BundlePolicy:         0,
		RTCPMuxPolicy:        0,
		PeerIdentity:         "",
		Certificates:         []webrtc.Certificate{},
		ICECandidatePoolSize: 0,
		SDPSemantics:         0,
	}

	router.POST("/webrtc/sdp/m/:meetingId/c/:userID/p/:peerId/s/:isSender", func(c *gin.Context) {
		isSender, _ := strconv.ParseBool(c.Param("isSender"))
		userID := c.Param("userID")
		peerID := c.Param("peerId")

		fmt.Println(c.Request)
		var session Sdp
		fmt.Printf("yolos --- %v", c.Request.Body)
		if err := c.ShouldBindJSON(&session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Fatal("sdp session json didn't bind: ", err.Error())
			return
		}
		offer := webrtc.SessionDescription{}
		utils.Decode(session.Sdp, &offer)

		// Create a new RTCPeerConnection
		// this is the gist of webrtc, generates and process SDP
		fmt.Println("Creating Connection")
		peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
		fmt.Println("Created Connection")
		if err != nil {
			log.Fatal("peer connection creation", err)
		}
		if !isSender {
			recieveTrack(peerConnection, peerConnectionMap, peerID)
		} else {
			createTrack(peerConnection, peerConnectionMap, userID)
		}
		// Set the SessionDescription of remote peer
		peerConnection.SetRemoteDescription(offer)

		// Create answer
		answer, err := peerConnection.CreateAnswer(nil)
		// fmt.Printf("Answer: %v", answer)
		if err != nil {
			log.Fatal("answer creationg error: ", err)
		}

		// Sets the LocalDescription, and starts our UDP listeners
		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, Sdp{Sdp: utils.Encode(answer)})
	})

	router.Run(":8080")
}

func createTrack(peerConnection *webrtc.PeerConnection,
	peerConnectionMap map[string]chan *webrtc.Track,
	currentUserID string) {

	fmt.Println("Creating Track")
	// Create a new RTCRtpTransceiver(SendRecv or RecvOnly) of RTPCodecTypeVideo kind and add it to the set of transceivers.
	if _, err := peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		log.Fatal(err)
	}

	// Register remote handler for remote track
	peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {

		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
		// The PictureLossIndication packet informs the encoder about the loss of an undefined amount of coded video data belonging to one or more pictures
		// SSRC - Synchronization Source - 32 bit random value that is generated in RTP to denote a specific source used to send media in an RTP connection.
		go func() {
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		// Create a local track, all our SFU clients will be fed via this track main track of the broadcaster
		// SSRC - Synchronization Source - 32 bit random value that is generated in RTP to denote a specific source used to send media in an RTP connection.
		localTrack, newTrackErr := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
		if newTrackErr != nil {
			log.Fatal(newTrackErr)
		}

		// the channel that will have the local track that is used by the sender
		// the localTrack needs to be fed to the reciever
		localTrackChan := make(chan *webrtc.Track, 1)
		localTrackChan <- localTrack
		if existingChan, ok := peerConnectionMap[currentUserID]; ok {
			// feed the exsiting track from user with this track
			existingChan <- localTrack
		} else {
			peerConnectionMap[currentUserID] = localTrackChan
		}

		rtpBuf := make([]byte, 1400)
		for { // for publisher only
			i, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				log.Fatal(readErr)
			}

			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
			if _, err := localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
				log.Fatal(err)
			}
		}

	})
}

// user is the caller of the method
// if user connects before peer: create channel and keep listening till track is added
// if peer connects before user: channel would have been created by peer and track can be added by getting the channel from cache
func recieveTrack(peerConnection *webrtc.PeerConnection,
	peerConnectionMap map[string]chan *webrtc.Track,
	peerID string) {
	fmt.Println("Receiving track")
	if _, ok := peerConnectionMap[peerID]; !ok {
		peerConnectionMap[peerID] = make(chan *webrtc.Track, 1)
	}
	localTrack := <-peerConnectionMap[peerID]
	fmt.Println("Received track")
	peerConnection.AddTrack(localTrack)
}

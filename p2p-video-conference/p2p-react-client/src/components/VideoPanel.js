import React, { useRef, useContext, useEffect } from 'react'
import { MeetingContext } from '../MeetingContext'

const urlPrefix = 'http://192.168.29.63:2223/webrtc/sdp/'

const VideoPanel = () => {
  const senderVideo = useRef()
  const recieverVideo = useRef()
  const { activeMeetingId, userId, peerId, isSender } = useContext(MeetingContext)

  const pcSender = new RTCPeerConnection({
    iceServers: [
      {
        urls: 'stun:stun.l.google.com:19302'
      }
    ]
  })

  const pcReciever = new RTCPeerConnection({
    iceServers: [
      {
        urls: 'stun:stun.l.google.com:19302'
      }
    ]
  })
  pcSender.onicecandidate = async (e) => {
    fetch(`${urlPrefix}m/${activeMeetingId}/c/${userId}/p/${peerId}/s/true`, {
      method: 'POST',
      body: JSON.stringify({ "sdp": btoa(JSON.stringify(pcSender.localDescription)) }),
    }).then(response => response.json())
      .then(response => {

        console.log('pcSender response - ', response)
        pcSender.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(response.Sdp))))
      })
  }

  pcReciever.onicecandidate = async (e) => {
    fetch(`${urlPrefix}m/${activeMeetingId}/c/${userId}/p/${peerId}/s/false`, {
      method: 'POST',
      body: JSON.stringify({ "sdp": btoa(JSON.stringify(pcReciever.localDescription)) }),
    }).then(response => response.json())
      .then(response => {

        console.log('pcReceiver response - ', response)
        pcReciever.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(response.Sdp))))
      })
  }

  useEffect(() => {
    if (!activeMeetingId) {
      senderVideo.current.style.visibility = 'hidden'
      recieverVideo.current.style.visibility = 'hidden'
      recieverVideo.current.srcObject = null

    } else {
      senderVideo.current.style.visibility = 'visible'
      recieverVideo.current.style.visibility = 'visible'

      // Start call
      navigator.mediaDevices.getUserMedia({ video: true, audio: true }).then((stream) => {
        senderVideo.current.srcObject = stream;
        var tracks = stream.getTracks();
        for (var i = 0; i < tracks.length; i++) {
          pcSender.addTrack(tracks[i])
        }
        pcSender.createOffer().then(d => pcSender.setLocalDescription(d))
      })

      // you can use event listner so that you inform he is connected!
      pcSender.addEventListener('connectionstatechange', event => {
        if (pcSender.connectionState === 'connected') {
          console.log("horray!")
        }
      });

      pcReciever.addTransceiver('video', { 'direction': 'recvonly' })

      pcReciever.createOffer().then(d => pcReciever.setLocalDescription(d))

      pcReciever.ontrack = (event) => {
        recieverVideo.current.srcObject = event.streams[0]
        recieverVideo.current.controls = true
        recieverVideo.current.autoplay = true
      }

    }
  }, [activeMeetingId])

  return (
    <div >
      <video autoPlay ref={recieverVideo} width="500" height="500" controls muted ></video>
      <div>
        <video autoPlay ref={senderVideo} width="160" height="120" controls muted ></video>
      </div>
    </div>
  )
}

export default VideoPanel

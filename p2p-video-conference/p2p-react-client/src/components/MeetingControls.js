import React from 'react'
import { useState, useContext } from 'react'
import { MeetingContext } from '../MeetingContext'

const MeetingControls = ({ onStartMeeting, onEndMeeting }) => {
  const meetingContext = useContext(MeetingContext)
  const [meetingId, setMeetingId] = useState('')
  const [userId, setUserId] = useState('')
  const [peerId, setPeerId] = useState('')
  const startMeetingHandler = (e) => {
    if (!meetingId) {
      alert('No Meeting Id')
    }
    // Call backend and start a new meeting
    console.log('start meeting', e)
    onStartMeeting(meetingId, userId, peerId, true)
  };
  const joinMeetingHandler = (e) => {
    if (!meetingId) {
      alert('No Meeting Id')
    }
    // Join existing meeting
    console.log('join meeting', e)
    onStartMeeting(meetingId, userId, peerId, false)
  };
  const endMeetinghandler = (e) => {
    console.log('End meeting', e)
    onEndMeeting()
  };

  if (meetingContext.activeMeetingId) {
    return (
      <div className='meeting-controls'>
        <h2 style={{ padding: '10px' }}>Meeting Id: {meetingContext.activeMeetingId}</h2>
        <button
          style={{ backgroundColor: 'red' }}
          onClick={endMeetinghandler}
        >End Meeting</button>
      </div>
    )
  } else {
    return (
      <div className='meeting-controls'>
        <input type='text' placeholder='Meeting Key' value={meetingId} onChange={(e) => setMeetingId(e.target.value)} />
        <input type='text' placeholder='User Id' value={userId} onChange={(e) => setUserId(e.target.value)} />
        <input type='text' placeholder='Peer Id' value={peerId} onChange={(e) => setPeerId(e.target.value)} />
        <button
          onClick={startMeetingHandler}
        >Start Meeting</button>
        <button
          onClick={joinMeetingHandler}
        >Join Meeting</button>
      </div>
    )
  }
}

export default MeetingControls

import './App.css';
import Header from './components/Header'
import MeetingControls from './components/MeetingControls';
import { useState } from 'react'
import { MeetingContext } from './MeetingContext';
import VideoPanel from './components/VideoPanel';


function App() {
  const [activeMeetingId, setActiveMeetingId] = useState(null)
  const [userId, setUserId] = useState(null)
  const [peerId, setPeerId] = useState(null)
  const [isSender, setIsSender] = useState(true)

  const onStartMeeting = (mId, userId, peerId, isSender) => {
    // Set active meeting id
    setActiveMeetingId(mId)
    setUserId(userId)
    setPeerId(peerId)
    setIsSender(isSender)
    // Call BE to start meeting
  }

  const onEndMeeting = () => {
    setActiveMeetingId(null)
  }

  return (
    <div className="App">
      <Header />
      <MeetingContext.Provider value={{ activeMeetingId, userId, peerId, isSender }}>
        <MeetingControls onStartMeeting={onStartMeeting} onEndMeeting={onEndMeeting} />
        <VideoPanel />
      </MeetingContext.Provider>
    </div>
  );
}

export default App;

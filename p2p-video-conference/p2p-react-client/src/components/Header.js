import PropTypes from 'prop-types'

const Header = (props) => {
  return (
    <header>
      <h1>Teaching Space</h1>
    </header>
  )
}


Header.defaultProps = {
  meetingKey: 'Meeting not started yet'
}

Header.propTypes = {
  meetingKey: PropTypes.string
}

export default Header

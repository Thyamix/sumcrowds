import { useEffect, useRef, useState } from 'react'
import '../App.css'
import { Navigate, useNavigate, useParams } from 'react-router-dom';
import { EnterPassword } from '../components/enterPassword';
import { fetchWithAuth } from '../utils/auth';

/** @type {string} */
const WSURL = import.meta.env.VITE_WSURL;
/** @type {string} */
const APIURL = import.meta.env.VITE_APIURL;

export function Counter() {
  /** @type {[boolean, (boolean) => void]} */
  const [loading, setLoading] = useState(true)
  /** @type {[boolean, (boolean) => void]} */
  const [isValid, setIsValid] = useState(null)
  /** @type {[boolean, (boolean) => void]} */
  const [hasAccess, setHasAccess] = useState(false)

  const { festivalCode } = useParams()
  const navigate = useNavigate()

  useEffect(() => {
    const checkFestival = async () => {
      const response = await fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/access")

      if (response.ok) {
        setHasAccess(true)
        setIsValid(true)
      }
      if (response.status == 404) {
        setIsValid(false)
      }
      if (response.status == 403) {
        setIsValid(true)
        setHasAccess(false)
      }
      setLoading(false)
    }
    checkFestival()
  }, [festivalCode, navigate])

  if (!loading) {
    if (!isValid) {
      return <Navigate to="/home" replace />
    }
    if (!hasAccess) {
      return <EnterPassword />
    }
    return <FestivalCountedPage />
  }
  return <div> loading ... </div>
}

function FestivalCountedPage() {
  /** @type {[int, () => null]} */
  const [total, setTotal] = useState("Loading...")
  /** @type {[int, () => null]} */
  const [maxJauge, setMaxJauge] = useState(0)

  /** @type {[string, () => null]} */
  const [colour, setColour] = useState("#ffffff")

  /** @type {React.RefObject < WebSocket >} */
  const socketRef = useRef(null)
  /** @type {string} */
  const { festivalCode } = useParams()

  useEffect(() => {
    socketRef.current = new WebSocket(WSURL + festivalCode)

    socketRef.current.onopen = () => {
      getTotal(socketRef.current)
      console.log("Open")
    }

    socketRef.current.onmessage = (event) => {
      handleMessage(event, setTotal, setMaxJauge)
    }

    socketRef.current.onclose = () => {
      console.log("Closed")
      location.reload()
    }
  }, [])

  function ColourSelector() {
    if (total >= maxJauge) {
      setColour("#ffa6a6")
    } else if (total >= (maxJauge * 0.9)) {
      setColour("#ffff87")
    } else {
      setColour("#ffffff")
    }
  }


  return (
    <div className='counter-page'>
      <ColourSelector total={total} max={maxJauge} />

      <div className="counter-main-container">
        <div className="counter-info-bar">
          <div className="counter-info-item">
            <span className="counter-info-label">CODE</span>
            <span className="counter-info-value">{festivalCode}</span>
          </div>
        </div>

        <div className="counter-display-section" style={{ background: colour }}>
          <div className="counter-current-value">{total}</div>
          <div className="counter-gauge-info">
            <span className="counter-gauge-label">GAUGE</span>
            <span className="counter-gauge-value">{maxJauge}</span>
          </div>
        </div>

        <div className="counter-controls">
          <div className="counter-column counter-column--decrease">
            <div className="counter-small-buttons">
              <button
                id="reduceThree"
                className="counter-button counter-button--small counter-button--decrease"
                onClick={() => {
                  handleMinus(3, socketRef.current)
                  if (total < 3) {
                    setTotal(0)
                  } else {
                    setTotal(total - 3)
                  }
                }}
              >
                −3
              </button>
              <button
                id="reduceTwo"
                className="counter-button counter-button--small counter-button--decrease"
                onClick={() => {
                  handleMinus(2, socketRef.current)
                  if (total < 2) {
                    setTotal(0)
                  } else {
                    setTotal(total - 2)
                  }
                }}
              >
                −2
              </button>
            </div>
            <div className="counter-small-buttons">
              <button
                id="reduceOne"
                className="counter-button counter-button--large counter-button--decrease"
                onClick={() => {
                  handleMinus(1, socketRef.current)
                  if (total < 1) {
                    setTotal(0)
                  } else {
                    setTotal(total - 1)
                  }
                }}
              >
                −1
              </button>
            </div>
          </div>

          <div className="counter-column counter-column--increase">
            <div className="counter-small-buttons">
              <button
                id="addTwo"
                className="counter-button counter-button--small counter-button--increase"
                onClick={() => {
                  handlePlus(2, socketRef.current)
                  setTotal(total + 2)
                }}
              >
                +2
              </button>
              <button
                id="addThree"
                className="counter-button counter-button--small counter-button--increase"
                onClick={() => {
                  handlePlus(3, socketRef.current)
                  setTotal(total + 3)
                }}
              >
                +3
              </button>
            </div>
            <div className="counter-small-buttons">
              <button
                id="addOne"
                className="counter-button counter-button--large counter-button--increase"
                onClick={() => {
                  handlePlus(1, socketRef.current)
                  setTotal(total + 1)
                }}
              >
                +1
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="counter-spacer" />
    </div>
  )
}


/**
 * @param {WebSocket} socket
        * @param {int} amount
        */
async function handlePlus(amount, socket) {
  socket.send(JSON.stringify({
    type: "inc",
    content: {
      amount: amount
    }
  }))
}

/**
        * @param {WebSocket} socket
        * @param {int} amount
        */
async function handleMinus(amount, socket) {
  socket.send(JSON.stringify({
    type: "dec",
    content: {
      amount: amount
    }
  }))
}

/**
        * @param {WebSocket} socket
        */
async function getTotal(socket) {
  socket.send(JSON.stringify({
    type: "getTotal",
  }))
}

/**
        * @param {Event} event
        * @param {() => null} setJauge
        * @param {() => null} setTotal
        **/
function handleMessage(event, setTotal, setJauge) {
  let result = JSON.parse(event.data)
  setTotal(result.total)
  setJauge(result.jauge)
}

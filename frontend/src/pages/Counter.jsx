import { useEffect, useRef, useState } from 'react'
import '../App.css'
import { Link, Navigate, useNavigate, useParams } from 'react-router-dom';
import { EnterPassword } from '../components/enterPassword';
import { fetchWithAuth } from '../utils/auth';
import { useTranslation } from 'react-i18next';
import { LanguageSwitcher } from '../components/languageSwitcher';
import HomeIcon from '../assets/home.svg?react';
import AdminIcon from '../assets/admin.svg?react';

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

  /** @type {string} */
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
  const { t } = useTranslation()

  /** @type {[int, () => null]} */
  const [total, setTotal] = useState("Loading...")
  /** @type {[int, () => null]} */
  const [maxJauge, setMaxJauge] = useState(0)

  /** @type {[string, () => null]} */
  const [colour, setColour] = useState("#ffffff")

  /** @type {string} */
  const { festivalCode } = useParams()

  /** @type {React.RefObject < WebSocket >} */
  const socketRef = useRef(null)
  /** @type {NodeJS.Timeout} */
  let heartbeat

  useEffect(() => {
    socketRef.current = new WebSocket(WSURL + festivalCode)

    socketRef.current.onopen = () => {
      getTotal(socketRef.current)
      console.log("Open")
      heartbeat = setInterval(() => {
        if (socketRef.current.readyState == WebSocket.OPEN) {
          socketRef.current.send(JSON.stringify({ type: "ping" }))
        }
      }, 10000)

    }

    socketRef.current.onmessage = (event) => {
      handleMessage(event, setTotal, setMaxJauge)
    }

    socketRef.current.onclose = () => {
      setTotal("Disconnected")
      clearInterval(heartbeat)
      console.log("Closed")
    }
  }, [])

  function ColourSelector() {
    useEffect(() => {
      if (total >= maxJauge) {
        setColour("#ffa6a6")
      } else if (total >= (maxJauge * 0.9)) {
        setColour("#ffff87")
      } else {
        setColour("#ffffff")
      }
    })
  }

  return (
    <div className='counter-page'>
      <ColourSelector />
      <div className="counter-main-container">
        <LanguageSwitcher />
        <div className="counter-info-bar">
          <Link to="/home" className="corner-button corner-button--left">
            <HomeIcon />
          </Link>
          <div className="counter-info-item">
            <span className="counter-info-label">{t("counter_code")}</span>
            <span className="counter-info-value">{festivalCode}</span>
          </div>
          <Link to={`/${festivalCode}/admin`} className="corner-button corner-button--right">
            <AdminIcon />
          </Link>
        </div>

        <div className="counter-display-section" style={{ background: colour }}>
          <div className="counter-current-value">{total}</div>
          <div className="counter-gauge-info">
            <span className="counter-gauge-label">{t("counter_gauge")}</span>
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
                  handleMinus(3, festivalCode)
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
                  handleMinus(2, festivalCode)
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
                  handleMinus(1, festivalCode)
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
                  handlePlus(2, festivalCode)
                  setTotal(total + 2)
                }}
              >
                +2
              </button>
              <button
                id="addThree"
                className="counter-button counter-button--small counter-button--increase"
                onClick={() => {
                  handlePlus(3, festivalCode)
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
                  handlePlus(1, festivalCode)
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
 * @param {int} amount
*/
async function handlePlus(amount, festivalCode) {
  const body = JSON.stringify({
    amount: amount
  })
  fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/inc", {
    method: "post",
    headers: {
      "Content-Type": "application/json",
    },
    body: body,
  })
}

/**
 * @param {int} amount
*/
async function handleMinus(amount, festivalCode) {
  const body = JSON.stringify({
    amount: amount
  })
  fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/dec", {
    method: "post",
    headers: {
      "Content-Type": "application/json",
    },
    body: body,
  })
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
  if (result.type == "pong") {
    return
  }
  setTotal(result.total)
  setJauge(result.jauge)
}

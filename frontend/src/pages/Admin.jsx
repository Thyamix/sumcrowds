import { useEffect, useState } from 'react';
import '../App.css'
import { useNavigate, useParams } from 'react-router-dom';
import { fetchWithAuth } from '../utils/auth';

/** @type {string} */
const APIURL = import.meta.env.VITE_APIURL;

export function Admin() {
  /** @type {[string, () => null]} */
  const [alert, setAlert] = useState(" ")
  /** @type {[string, () => null]} */
  const [inputValue, setInputValue] = useState("")

  const { festivalCode } = useParams()


  /** 
   * @param {Event} event 
  */
  function handleClick(event) {
    event.preventDefault()

    let valid = true

    for (let i = 0; i < inputValue.length; i++) {
      if (!("1234567890".includes(inputValue.at(i)))) {
        valid = false
        break
      }
    }
    if (inputValue.length == 0) {
      valid = false
    }

    if (!valid) {
      playAlert("Please only use numbers", setAlert)
    } else {
      onSetGaugePressed(inputValue)
      setInputValue("")
    }
  }

  /** 
   * @param {Event} event 
  */
  function handleInputValue(event) {
    const value = event.target.value
    if ("1234567890".includes(value.at(-1)) || value == "") {
      setInputValue(value)
    }
  }

  return (
    <div className='admin-page'>
      <div className="admin-main-container">
        <div className='admin-header'>
          Admin Panel
        </div>

        <form className='admin-form'>
          <Alert />

          <div className='admin-input-group'>
            <input
              type='text'
              name='maxGauge'
              value={inputValue}
              onChange={handleInputValue}
              placeholder='Max Gauge'
              className='admin-input'
            />
            <button
              className='admin-button admin-button--primary admin-button--large'
              onClick={handleClick}
            >
              Set Gauge
            </button>
          </div>

          <div className='admin-divider' />

          <div className='admin-section'>
            <div className='admin-section-title'>
              Current Event
            </div>
            <div className="admin-button-group">
              <button
                className='admin-button admin-button--danger admin-button--small'
                onClick={onArchivePressed}
              >
                Archive
              </button>
              <a
                href={APIURL + "v1/festival/" + festivalCode + "/download/activecsv"}
                className='admin-button admin-button--success admin-button--small'
              >
                Get CSV
              </a>
            </div>
          </div>
        </form>

        <Archive />
      </div >
    </div >
  )

  function Alert() {
    if (alert != " ") {
      return (<div className='admin-alert'>
        {alert}
      </div>)
    }

  }

  /**
   * @param {Object} param0 
   * @param {int} param0.id 
   * @param {int} param0.date 
   */
  function ArchivedData({ id, date }) {
    return (
      <a
        href={APIURL + "v1/festival/" + festivalCode + "/download/archivedcsv/" + id}
        className='archive-item'
      >
        <div className='archive-item-content'>
          <span className='archive-item-id'>ID: {id}</span>
          <span className='archive-item-date'>{date}</span>
        </div>
      </a>
    )
  }

  function Archive() {
    const [archives, setArchives] = useState([])

    useEffect(() => {
      fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/getarchivedevents")
        .then(res => res.json())
        .then(data => setArchives(data))
    }, [])

    function getDateTime(timestamp) {
      const time = new Date(timestamp * 1000)
      console.log(time)

      return time.toLocaleString().slice(0, 24)
    }

    return (
      <div className='archive-section'>
        <div className='archive-header'>
          Archived Data
        </div>
        <div className='archive-divider' />
        <div className='archive-list'>
          {archives.map((item) =>
            <ArchivedData key={item.id} id={item.id} date={getDateTime(item.time)} />
          )}
        </div>
      </div>
    )
  }
  /**
   * @param {string} alert 
   * @param {() => null} setAlert 
  */
  async function playAlert(alert, setAlert) {
    setAlert(alert)
    await new Promise(resolve => setTimeout(resolve, 7500))
    setAlert(" ")
  }

  /**
   * @param {Event} event
  */
  async function onArchivePressed(event) {
    event.preventDefault()
    const response = await fetch(APIURL + "v1/festival/" + festivalCode + "/archivecurrentevent", {
      method: "get",
      headers: {
        "Content-Type": "application/json"
      }
    })
    if (!response.ok) {
      throw new Error(`Response status:`, response.status)
    }
    location.reload()
  }

  /**
   * @param {() => null} newMax 
  */
  async function onSetGaugePressed(newMax) {
    const body = JSON.stringify({
      max: +newMax
    })
    const response = await fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/setgauge", {
      method: "post",
      headers: {
        "Content-Type": "application/json",
      },
      body: body,
    })
    if (!response.ok) {
      throw new Error("Response status:", response.status)
    }
  }

}

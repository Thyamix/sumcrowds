import { useEffect, useState } from 'react';
import '../App.css'
import { Link, Navigate, useNavigate, useParams } from 'react-router-dom';
import { fetchWithAuth } from '../utils/auth';
import { EnterPin } from '../components/enterPin';
import HomeIcon from '../assets/home.svg?react';
import BackIcon from '../assets/back.svg?react';
import { useTranslation } from 'react-i18next';
import { LanguageSwitcher } from '../components/languageSwitcher';

/** @type {string} */
const APIURL = import.meta.env.VITE_APIURL;

export function Admin() {
  /** @type {[boolean, (boolean) => void]} */
  const [loading, setLoading] = useState(true)
  /** @type {[boolean, (boolean) => void]} */
  const [isValid, setIsValid] = useState(null)
  /** @type {[boolean, (boolean) => void]} */
  const [hasAccess, setHasAccess] = useState(false)
  /** @type {[boolean, (boolean) => void]} */
  const [hasAdminAccess, setHasAdminAccess] = useState(false)

  const { festivalCode } = useParams()

  const navigate = useNavigate()



  useEffect(() => {
    const checkFestival = async () => {
      const response = await fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/admin/access")

      if (response.ok) {
        setHasAccess(true)
        setHasAdminAccess(true)
        setIsValid(true)
      }
      if (response.status == 422) {
        setIsValid(true)
        setHasAccess(true)
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
    } else if (!hasAdminAccess) {
      return <EnterPin />
    }
    return <FestivalAdminPage />
  }
  return <div> loading ... </div>
}


function FestivalAdminPage() {
  /** @type {[string, () => null]} */
  const [alert, setAlert] = useState(" ")
  /** @type {[string, () => null]} */
  const [inputValue, setInputValue] = useState("")

  const { t } = useTranslation()

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
        <LanguageSwitcher />
        <div className="counter-info-bar">
          <Link to="/home" className="corner-button corner-button--left">
            <HomeIcon className="corner-icon" />
          </Link>
          <div className='admin-header'>
            {t("admin_title")}
          </div>
          <Link to={`/${festivalCode}`} className="corner-button corner-button--right">
            <BackIcon className="corner-icon" />
          </Link>
        </div>
        <form className='admin-form'>
          <Alert />

          <div className='admin-input-group'>
            <input
              type='text'
              name='maxGauge'
              value={inputValue}
              onChange={handleInputValue}
              placeholder={t("admin_max_gauge")}
              className='admin-input'
            />
            <button
              className='admin-button admin-button--primary admin-button--large'
              onClick={handleClick}
            >
              {t("admin_set_gauge")}
            </button>
          </div>

          <div className='admin-divider' />

          <div className='admin-section'>
            <div className='admin-section-title'>
              {t("admin_current_event")}
            </div>
            <div className="admin-button-group">
              <button
                className='admin-button admin-button--danger admin-button--small'
                onClick={onArchivePressed}
              >
                {t("admin_archive")}
              </button>
              <a
                href={APIURL + "v1/festival/" + festivalCode + "/admin/download/activecsv"}
                className='admin-button admin-button--success admin-button--small'
              >
                {t("admin_get_csv")}
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
        href={APIURL + "v1/festival/" + festivalCode + "/admin/download/archivedcsv/" + id}
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
      const response = fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/admin/getarchivedevents")
        .then(res => res.json())
        .then(data => setArchives(data))
      if (!response.ok) {
        if (response.status == 422) {
          location.reload()
        }
      }
    }, [])

    function getDateTime(timestamp) {
      const time = new Date(timestamp * 1000)
      return time.toLocaleString().slice(0, 24)
    }

    return (
      <div className='archive-section'>
        <div className='archive-header'>
          {t("admin_archived_data")}
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
    const response = await fetch(APIURL + "v1/festival/" + festivalCode + "/admin/archivecurrentevent", {
      method: "get",
      headers: {
        "Content-Type": "application/json"
      }
    })
    if (!response.ok) {
      if (response.status == 422) {
        location.reload()
      }
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
    const response = await fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/admin/setgauge", {
      method: "post",
      headers: {
        "Content-Type": "application/json",
      },
      body: body,
    })
    if (!response.ok) {
      if (response.status == 422) {
        location.reload()
      }
      throw new Error("Response status:", response.status)
    }
  }

}

import { useEffect, useState } from "react"
import { auth, fetchWithAuth } from "../utils/auth"
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";

/** @type {string} */
const APIURL = import.meta.env.VITE_APIURL;

/**
 * @param {Object} param0 
 * @param {() => null} param0.close 
*/
export function JoinPopup({ close }) {
  const { t } = useTranslation()

  useEffect(() => { auth() }, [])
  /** @type {[string, () => null]} */
  const [codeInputValue, setCodeInputValue] = useState("")
  const [alert, setAlert] = useState(" ")

  const navigate = useNavigate()

  /** 
   * @param {Event} event 
  */
  function handleCodeInputValue(event) {
    /** @type {string} */
    const value = event.target.value
    if (("abcdefghijklmnopqrstuvwxyx1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ".includes(value.at(-1)) || value == "") && value.length <= 6) {
      setCodeInputValue(value.toUpperCase())
    }
  }

  function Alert() {
    if (alert != " ") {
      return (<div className='admin-alert'>
        {alert}
      </div>)
    }
  }

  /**
   * @param {Event} event 
  */
  async function handleJoin(event) {
    event.preventDefault()
    const response = await fetchWithAuth(APIURL + "v1/festival/" + codeInputValue + "/exists")
    if (response.ok) {
      navigate("/" + codeInputValue, { replace: true })
    } else {
      setAlert(t("joinpopup_alert"))
    }
  }

  return (
    <div className="join-modal">
      <div className="join-main-container">
        <div className="join-space" />
        <button className="join-close-button" onClick={close}>
          <b>×</b>
        </button>
        <div className="join-header">
          {t("joinpopup_header")}
        </div>
        <div className="join-spacer" />
        <div className="join-form">
          <div className="join-input-group">
            <Alert />
            <input
              type="text"
              name="Code"
              value={codeInputValue}
              onChange={handleCodeInputValue}
              placeholder={t("joinpopup_enter_code")}
              className="join-input"
            />
          </div>
          <button
            className="join-button join-button--large join-button--success"
            onClick={handleJoin}
          >
            {t("joinpopup_confirm")}
          </button>
        </div>
      </div>
    </div>)
}

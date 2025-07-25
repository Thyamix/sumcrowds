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
export function CreatePopup({ close }) {
  useEffect(() => { auth() }, [])
  /** @type {[string, () => null]} */
  const [passwordInputValue, setPasswordInputValue] = useState("")
  /** @type {[string, () => null]} */
  const [pinInputValue, setPinInputValue] = useState("")
  /** @type {[boolean, () => null]} */
  const [showPassword, setShowPassword] = useState(false);

  const { t } = useTranslation()

  const navigate = useNavigate()
  /** 
   * @param {Event} event 
  */
  function handlePinInputValue(event) {
    /** @type {string} */
    const value = event.target.value
    if (("1234567890".includes(value.at(-1)) || value == "") && value.length <= 4) {
      setPinInputValue(value)
    }
  }

  function togglePasswordVisibility() {
    setShowPassword(!showPassword);
  };

  /** 
   * @param {Event} event 
  */
  function handlePasswordInputValue(event) {
    /** @type {string} */
    const value = event.target.value
    if ("abcdefghijklmnopqrstuvwxyx1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-+_*".includes(value.at(-1)) || value == "") {
      setPasswordInputValue(value)
    }
  }

  /**
   * @param {Event} event 
  */
  async function handleCreate(event) {
    event.preventDefault()
    const body = JSON.stringify({
      password: passwordInputValue,
      pin: pinInputValue,
    })
    await fetchWithAuth(APIURL + "v1/create/festival", {
      method: "post",
      body: body,
    }).then(response => response.json())
      .then(data => {
        if (data.type == "festival code" && data.content != null) {
          navigate("/" + data.content, { replace: true })
        }
      })
  }

  return (
    <div className="create-modal">
      <div className="create-main-container">
        <div className="create-space" />
        <button className="create-close-button" onClick={close}>
          <b>×</b>
        </button>
        <div className="create-header">{t("createpopup_header")}</div>
        <div className="create-spacer" />
        <form className="create-form">
          <input
            type="text"
            name="email"
            autoComplete="username email"
            style={{ display: "none" }}
          />

          <div className="create-input-group">
            <input
              type="text"
              min={0}
              maxLength={4}
              step={1}
              name="pin"
              value={pinInputValue}
              onChange={handlePinInputValue}
              placeholder={t("createpopup_admin_pin")}
              className="create-input"
            />
          </div>

          <div className="create-input-group">
            <input
              id="create-password"
              type={showPassword ? 'text' : 'password'}
              maxLength={50}
              autoComplete="new-password"
              name="password"
              value={passwordInputValue}
              onChange={handlePasswordInputValue}
              placeholder={t("createpopup_password")}
              className="create-input"
            />
          </div>

          <div className="create-checkbox-container">
            <label className="create-checkbox-label">
              <input
                className="create-checkbox"
                type="checkbox"
                onChange={togglePasswordVisibility}
                name="show-password"
              />
              <span className="create-checkbox-text">{t("createpopup_show_password")}</span>
            </label>
          </div>
        </form>

        <button
          className="create-button create-button--large create-button--primary"
          onClick={handleCreate}
        >
          {t("createpopup_confirm")}
        </button>
      </div>
    </div>
  )
}


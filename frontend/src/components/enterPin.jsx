import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router-dom";
import { fetchWithAuth } from "../utils/auth";

/** @type {string} */
const APIURL = import.meta.env.VITE_APIURL;

export function EnterPin() {
  const { t } = useTranslation()

  /** @type {[string, React.Dispatch<React.SetStateAction<string>>]} */
  const [pinInputValue, setPinInputValue] = useState("");
  /** @type {[boolean, React.Dispatch<React.SetStateAction<boolean>>]} */
  const [showPin, setShowPin] = useState(false);
  /** @type {[string | null, React.Dispatch<React.SetStateAction<string | null>>]} */
  const [pinError, setPinError] = useState(null);


  const { festivalCode } = useParams()
  const navigate = useNavigate()

  /**
   * @param {Event} event
   */
  function handlePinInputValue(event) {
    /** @type {string} */
    const value = event.target.value;
    if (value.length > 4) {
      return
    }
    if (
      "0123456789".includes(
        value.at(-1)
      ) ||
      value === ""
    ) {
      setPinInputValue(value);
      setPinError(null);
    }
  }

  /**
   * @param {Event} event
   */
  async function handleConfirm(event) {
    event.preventDefault();
    setPinError(null);

    try {
      const response = await fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/admin/access", {
        headers: {
          "admin-pin": pinInputValue,
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        setPinError(t("pinpopup_alert"));
        console.error("API call failed:", response.statusText);
      } else {
        console.log("PIN confirmed and API call initiated!");
        setPinInputValue("");
        location.reload()
      }
    } catch (error) {
      console.error("Error making API call:", error);
      setPinError("An error occurred. Please try again later.");
    }
  }

  function goHome() {
    navigate("/", { replace: true })
  }

  return (
    <div className="join-modal">
      <div className="join-main-container">
        <div className="join-space" />
        <button className="join-close-button" onClick={goHome}>
          {" "}
          <b> x </b>{" "}
        </button>
        <div className="join-header"> {t("pinpopup_header")}</div>
        <div className="spacer" />
        <form onSubmit={handleConfirm}>
          <input
            type="text"
            name="email"
            autoComplete="username"
            style={{ display: "none" }}
          />{" "}
          {pinError && <p style={{ color: 'red' }}>{pinError}</p>} {/* Invalid PIN warning */}
          <input
            id="pin"
            type={showPin ? "text" : "password"}
            maxLength={50}
            name="pin"
            value={pinInputValue}
            onChange={handlePinInputValue}
            placeholder={t("pinpopup_pin")}
            className="join-input"
          />
          <div className="checkbox-container">
            <input
              className="checkbox"
              type="checkbox"
              checked={showPin}
              onChange={() => setShowPin(!showPin)}
              name="show-pin"
            />{t("pinpopup_show_pin")}</div>
          <button type="submit" className="join-button join-button--large join-button--success">
            {t("pinpopup_confirm")}
          </button>
        </form>
      </div>
    </div>
  );
}

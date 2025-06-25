import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

/** @type {string} */
const APIURL = import.meta.env.VITE_APIURL;

export function EnterPassword() {
  /** @type {[string, React.Dispatch<React.SetStateAction<string>>]} */
  const [passwordInputValue, setPasswordInputValue] = useState("");
  /** @type {[boolean, React.Dispatch<React.SetStateAction<boolean>>]} */
  const [showPassword, setShowPassword] = useState(false);
  /** @type {[string | null, React.Dispatch<React.SetStateAction<string | null>>]} */
  const [passwordError, setPasswordError] = useState(null);


  const { festivalCode } = useParams()
  const navigate = useNavigate()

  /**
   * @param {Event} event
   */
  function handlePasswordInputValue(event) {
    /** @type {string} */
    const value = event.target.value;
    if (
      "abcdefghijklmnopqrstuvwxyx1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-+_*".includes(
        value.at(-1)
      ) ||
      value === ""
    ) {
      setPasswordInputValue(value);
      setPasswordError(null); // Clear error when user starts typing again
    }
  }

  /**
   * @param {Event} event
   */
  async function handleConfirm(event) {
    event.preventDefault();
    setPasswordError(null); // Clear previous errors

    const body = JSON.stringify({
      password: passwordInputValue,
    });

    try {
      const response = await fetch(APIURL + "v1/festival/" + festivalCode + "/getaccess", {
        method: "post",
        headers: {
          "Content-Type": "application/json",
        },
        body: body,
      });

      if (!response.ok) {
        setPasswordError("Invalid password. Please try again.");
        console.error("API call failed:", response.statusText);
      } else {
        console.log("Password confirmed and API call initiated!");
        setPasswordInputValue(""); // Make password disappear on success
        location.reload()
      }
    } catch (error) {
      console.error("Error making API call:", error);
      setPasswordError("An error occurred. Please try again later.");
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
        <div className="join-header"> Enter Password </div>
        <div className="spacer" />
        <form onSubmit={handleConfirm}>
          <input
            type="text"
            name="email"
            autoComplete="username"
            style={{ display: "none" }}
          />{" "}
          {passwordError && <p style={{ color: 'red' }}>{passwordError}</p>} {/* Invalid password warning */}
          <input
            id="password"
            type={showPassword ? "text" : "password"}
            maxLength={50}
            autoComplete="new-password"
            name="password"
            value={passwordInputValue}
            onChange={handlePasswordInputValue}
            placeholder="Password"
            className="join-input"
          />
          <div className="checkbox-container">
            <input
              className="checkbox"
              type="checkbox"
              checked={showPassword} // Reflect the state
              onChange={() => setShowPassword(!showPassword)} // Toggle showPassword state
              name="show-password"
            />Show Password</div>
          <button type="submit" className="join-button join-button--large join-button--success">
            Confirm
          </button>
        </form>
      </div>
    </div>
  );
}

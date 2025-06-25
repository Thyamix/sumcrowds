import Popup from "reactjs-popup";
import { JoinPopup } from "../components/joinFestivalButton";
import { CreatePopup } from "../components/createFestivalButton";
import { useTranslation } from "react-i18next";
import LanguageSwitcher from "../components/languageSwitcher";

export function Home() {
  const { t } = useTranslation()

  return (
    <div className="home-page">
      <div className="home-main-container">
        <LanguageSwitcher />
        <div className="home-header">
          {t("home_home")}
        </div>
        <div className="home-form">
          <div className="home-section">
            <div className="home-section-title">{t("home_welcome")}</div>
            <p style={{
              fontSize: '18px',
              color: '#6b7280',
              marginBottom: '32px',
              textAlign: 'center'
            }}>
              {t("home_select_option")}
            </p>

            <div className="home-button-group" style={{ flexDirection: 'column', gap: '20px' }}>
              <Popup
                trigger={
                  <button className="home-button home-button--large home-button--success">
                    {t("home_join_button")}
                  </button>
                }
                position={"center center"}
                modal
                nested
                lockScroll>
                {(close) => <JoinPopup close={close} />}
              </Popup>

              <Popup
                trigger={
                  <button className="home-button home-button--large home-button--primary">
                    {t("home_create_button")}
                  </button>
                }
                position={"center center"}
                modal
                nested
                lockScroll>
                {(close) => <CreatePopup close={close} />}
              </Popup>
            </div>
          </div>
        </div>
      </div>


    </div>
  )
}



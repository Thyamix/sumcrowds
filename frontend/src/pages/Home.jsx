import Popup from "reactjs-popup";
import { JoinPopup } from "../components/joinFestivalButton";
import { CreatePopup } from "../components/createFestivalButton";

export function Home() {

  return (
    <div className="home-page">
      <div className="home-main-container">
        <div className="home-header">
          Home
        </div>
        <div className="home-form">
          <div className="home-section">
            <div className="home-section-title">Welcome</div>
            <p style={{
              fontSize: '18px',
              color: '#6b7280',
              marginBottom: '32px',
              textAlign: 'center'
            }}>
              Choose an option to get started
            </p>

            <div className="home-button-group" style={{ flexDirection: 'column', gap: '20px' }}>
              <Popup
                trigger={
                  <button className="home-button home-button--large home-button--success">
                    Join Session
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
                    Create Session
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



import { useEffect, useRef, useState } from 'react'
import { Link, Navigate, useParams } from 'react-router-dom';
import { EnterPassword } from '../components/enterPassword';
import { LeaveConfirmModal } from '../components/leaveConfirmModal';
import { fetchWithAuth } from '../utils/auth';
import { useTranslation } from 'react-i18next';
import { LanguageSwitcher } from '../components/languageSwitcher';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Home, Settings, Minus, Plus } from 'lucide-react';
import { cn } from '@/lib/utils';

const WSURL = import.meta.env.VITE_WSURL;
const APIURL = import.meta.env.VITE_APIURL;

export function Counter() {
  const [loading, setLoading] = useState(true)
  const [isValid, setIsValid] = useState(null)
  const [hasAccess, setHasAccess] = useState(false)
  const { festivalCode } = useParams()

  useEffect(() => {
    const checkFestival = async () => {
      const response = await fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/access")

      if (response.ok) {
        setHasAccess(true)
        setIsValid(true)
      }
      if (response.status === 404) {
        setIsValid(false)
      }
      if (response.status === 403) {
        setIsValid(true)
        setHasAccess(false)
      }
      setLoading(false)
    }
    checkFestival()
  }, [festivalCode])

  if (!loading) {
    if (!isValid) {
      return <Navigate to="/home" replace />
    }
    if (!hasAccess) {
      return <EnterPassword />
    }
    return <FestivalCountedPage />
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center">
      <div className="animate-pulse text-primary text-lg font-medium">Loading...</div>
    </div>
  )
}

function FestivalCountedPage() {
  const { t } = useTranslation()
  const [total, setTotal] = useState("...")
  const [maxJauge, setMaxJauge] = useState(0)
  const [status, setStatus] = useState("normal")
  const [showLeaveModal, setShowLeaveModal] = useState(false)
  const { festivalCode } = useParams()
  const socketRef = useRef(null)
  const heartbeatRef = useRef(null)

  useEffect(() => {
    socketRef.current = new WebSocket(WSURL + festivalCode)

    socketRef.current.onopen = () => {
      getTotal(socketRef.current)
      heartbeatRef.current = setInterval(() => {
        if (socketRef.current.readyState === WebSocket.OPEN) {
          socketRef.current.send(JSON.stringify({ type: "ping" }))
        }
      }, 10000)
    }

    socketRef.current.onmessage = (event) => {
      handleMessage(event, setTotal, setMaxJauge)
    }

    socketRef.current.onclose = () => {
      setTotal("Disconnected")
      clearInterval(heartbeatRef.current)
    }

    return () => {
      clearInterval(heartbeatRef.current)
      socketRef.current?.close()
    }
  }, [festivalCode])

  useEffect(() => {
    if (typeof total === 'number' && maxJauge > 0) {
      if (total >= maxJauge) {
        setStatus("danger")
      } else if (total >= maxJauge * 0.9) {
        setStatus("warning")
      } else {
        setStatus("normal")
      }
    }
  }, [total, maxJauge])

  const percentage = (typeof total === 'number' && maxJauge > 0) ? Math.min((total / maxJauge) * 100, 100) : 0

  return (
    <div className="min-h-screen bg-background flex flex-col items-center justify-center p-4">
      <Card className="w-full max-w-2xl overflow-hidden shadow-xl border-0">
        {/* Header */}
        <div className="bg-primary px-6 py-5 flex items-center justify-between">
          <Button
            variant="ghost"
            size="icon"
            className="text-white hover:bg-white/20"
            onClick={() => setShowLeaveModal(true)}
          >
            <Home className="h-5 w-5" />
          </Button>

          <div className="text-center">
            <p className="text-xs uppercase tracking-widest text-white/70 mb-1">{t("counter_code")}</p>
            <p className="text-2xl font-bold text-white font-mono tracking-wider">{festivalCode}</p>
          </div>

          <div className="flex items-center gap-2">
            <LanguageSwitcher />
            <Link to={`/${festivalCode}/admin`}>
              <Button variant="ghost" size="icon" className="text-white hover:bg-white/20">
                <Settings className="h-5 w-5" />
              </Button>
            </Link>
          </div>
        </div>

        {/* Counter Display */}
        <div className={cn(
          "py-12 px-6 text-center transition-colors duration-500",
          status === "danger" && "bg-destructive/10",
          status === "warning" && "bg-warning/10",
          status === "normal" && "bg-card"
        )}>
          <p className={cn(
            "text-8xl md:text-9xl font-black font-mono tabular-nums transition-colors",
            status === "danger" && "text-destructive",
            status === "warning" && "text-warning",
            status === "normal" && "text-foreground"
          )}>
            {total}
          </p>

          {/* Progress bar */}
          <div className="mt-8 max-w-md mx-auto">
            <div className="flex justify-between text-sm text-muted-foreground mb-2">
              <span>{t("counter_gauge")}</span>
              <span className="font-mono font-bold">{typeof total === 'number' ? total : 0} / {maxJauge}</span>
            </div>
            <div className="h-3 bg-muted rounded-full overflow-hidden">
              <div
                className={cn(
                  "h-full rounded-full transition-all duration-500",
                  status === "danger" && "bg-destructive",
                  status === "warning" && "bg-warning",
                  status === "normal" && "bg-success"
                )}
                style={{ width: `${percentage}%` }}
              />
            </div>
          </div>
        </div>

        {/* Controls */}
        <div className="grid grid-cols-2 gap-0 border-t">
          {/* Decrease column */}
          <div className="bg-destructive/5 p-6 border-r">
            <p className="text-center text-sm font-medium text-destructive mb-4 uppercase tracking-wide">{t("counter_exit")}</p>
            <div className="flex flex-col items-center gap-3">
              <div className="flex gap-2">
                <Button
                  variant="destructive"
                  size="counter"
                  className="shadow-lg hover:shadow-glow-destructive"
                  onClick={() => handleMinus(2, festivalCode)}
                >
                  −2
                </Button>
                <Button
                  variant="destructive"
                  size="counter"
                  className="shadow-lg hover:shadow-glow-destructive"
                  onClick={() => handleMinus(3, festivalCode)}
                >
                  −3
                </Button>
              </div>
              <Button
                variant="destructive"
                size="counterLg"
                className="shadow-lg hover:shadow-glow-destructive"
                onClick={() => handleMinus(1, festivalCode)}
              >
                <Minus className="w-6 h-6 mr-2" />
                1
              </Button>
            </div>
          </div>

          {/* Increase column */}
          <div className="bg-success/5 p-6">
            <p className="text-center text-sm font-medium text-success mb-4 uppercase tracking-wide">{t("counter_enter")}</p>
            <div className="flex flex-col items-center gap-3">
              <div className="flex gap-2">
                <Button
                  variant="success"
                  size="counter"
                  className="shadow-lg hover:shadow-glow-success"
                  onClick={() => handlePlus(2, festivalCode)}
                >
                  +2
                </Button>
                <Button
                  variant="success"
                  size="counter"
                  className="shadow-lg hover:shadow-glow-success"
                  onClick={() => handlePlus(3, festivalCode)}
                >
                  +3
                </Button>
              </div>
              <Button
                variant="success"
                size="counterLg"
                className="shadow-lg hover:shadow-glow-success"
                onClick={() => handlePlus(1, festivalCode)}
              >
                <Plus className="w-6 h-6 mr-2" />
                1
              </Button>
            </div>
          </div>
        </div>
      </Card>

      <LeaveConfirmModal
        open={showLeaveModal}
        onClose={() => setShowLeaveModal(false)}
      />
    </div>
  )
}

async function handlePlus(amount, festivalCode) {
  const body = JSON.stringify({ amount })
  try {
    await fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/inc", {
      method: "post",
      headers: { "Content-Type": "application/json" },
      body,
    })
  } catch (error) {
    console.error("Failed to increment:", error)
  }
}

async function handleMinus(amount, festivalCode) {
  const body = JSON.stringify({ amount })
  try {
    await fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/dec", {
      method: "post",
      headers: { "Content-Type": "application/json" },
      body,
    })
  } catch (error) {
    console.error("Failed to decrement:", error)
  }
}

async function getTotal(socket) {
  socket.send(JSON.stringify({ type: "getTotal" }))
}

function handleMessage(event, setTotal, setJauge) {
  const result = JSON.parse(event.data)
  if (result.type === "pong") return
  setTotal(result.total)
  setJauge(result.jauge)
}

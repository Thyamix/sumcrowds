import { useEffect, useState } from "react"
import { auth, fetchWithAuth } from "../utils/auth"
import { useNavigate } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { History, ChevronRight, Loader2 } from "lucide-react"

const APIURL = import.meta.env.VITE_APIURL

interface Session {
  code: string
  last_used_at: number
}

interface SessionsResponse {
  sessions: Session[]
  has_more: boolean
  page: number
}

interface RecentSessionsPopupProps {
  close: () => void
}

export function RecentSessionsPopup({ close }: RecentSessionsPopupProps) {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState("")
  const [page, setPage] = useState(0)
  const [hasMore, setHasMore] = useState(false)

  useEffect(() => {
    auth()
    loadSessions(0, true)
  }, [])

  async function loadSessions(pageNum: number, reset: boolean = false) {
    if (reset) {
      setLoading(true)
    } else {
      setLoadingMore(true)
    }
    setError("")

    try {
      const response = await fetchWithAuth(APIURL + `v1/user/sessions?page=${pageNum}`)
      if (response.ok) {
        const data: SessionsResponse = await response.json()
        if (reset) {
          setSessions(data.sessions || [])
        } else {
          setSessions(prev => [...prev, ...(data.sessions || [])])
        }
        setHasMore(data.has_more)
        setPage(data.page)
      } else {
        setError(t("error_generic"))
      }
    } catch (err) {
      setError(t("error_generic"))
    } finally {
      setLoading(false)
      setLoadingMore(false)
    }
  }

  function handleLoadMore() {
    if (!loadingMore && hasMore) {
      loadSessions(page + 1, false)
    }
  }

  function handleSelect(code: string) {
    navigate("/" + code, { replace: true })
  }

  function formatDate(timestamp: number): string {
    const date = new Date(timestamp * 1000)
    return date.toLocaleDateString(undefined, {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })
  }

  return (
    <Dialog open={true} onOpenChange={() => close()}>
      <DialogContent onClose={close} className="sm:max-w-md border-0 shadow-2xl overflow-hidden p-0">
        {/* Colored header */}
        <div className="bg-accent px-6 py-6 text-center">
          <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center mx-auto mb-3">
            <History className="w-7 h-7 text-white" />
          </div>
          <DialogTitle className="text-2xl font-bold text-white">
            {t("recent_sessions_title")}
          </DialogTitle>
        </div>

        <div className="p-6 space-y-4">
          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          {loading ? (
            <div className="flex justify-center py-8">
              <Loader2 className="w-8 h-8 animate-spin text-primary" />
            </div>
          ) : sessions.length === 0 ? (
            <p className="text-center text-muted-foreground py-8">
              {t("recent_sessions_empty")}
            </p>
          ) : (
            <div className="space-y-2 max-h-[300px] overflow-y-auto">
              {sessions.map((session, index) => (
                <button
                  key={`${session.code}-${index}`}
                  onClick={() => handleSelect(session.code)}
                  className="w-full flex items-center justify-between p-4 rounded-lg bg-secondary hover:bg-secondary/80 transition-colors text-left"
                >
                  <div>
                    <p className="font-semibold font-mono text-lg">{session.code}</p>
                    <p className="text-sm text-muted-foreground">
                      {formatDate(session.last_used_at)}
                    </p>
                  </div>
                  <ChevronRight className="w-5 h-5 text-muted-foreground" />
                </button>
              ))}

              {hasMore && (
                <Button
                  variant="outline"
                  className="w-full"
                  onClick={handleLoadMore}
                  disabled={loadingMore}
                >
                  {loadingMore ? (
                    <Loader2 className="w-4 h-4 animate-spin mr-2" />
                  ) : null}
                  {t("recent_sessions_load_more")}
                </Button>
              )}
            </div>
          )}

          <Button
            variant="outline"
            size="lg"
            className="w-full"
            onClick={close}
          >
            {t("recent_sessions_close")}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}

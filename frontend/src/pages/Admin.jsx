import { useEffect, useState } from 'react';
import { Link, Navigate, useParams } from 'react-router-dom';
import { fetchWithAuth } from '../utils/auth';
import { EnterPassword } from '../components/enterPassword';
import { EnterPin } from '../components/enterPin';
import { useTranslation } from 'react-i18next';
import { LanguageSwitcher } from '../components/languageSwitcher';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Home, ArrowLeft, Download, Archive, Settings, FileDown } from 'lucide-react';

const APIURL = import.meta.env.VITE_APIURL;

export function Admin() {
  const [loading, setLoading] = useState(true)
  const [isValid, setIsValid] = useState(null)
  const [hasAccess, setHasAccess] = useState(false)
  const [hasAdminAccess, setHasAdminAccess] = useState(false)
  const { festivalCode } = useParams()

  useEffect(() => {
    const checkFestival = async () => {
      const response = await fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/admin/access")

      if (response.ok) {
        setHasAccess(true)
        setHasAdminAccess(true)
        setIsValid(true)
      }
      if (response.status === 422) {
        setIsValid(true)
        setHasAccess(true)
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
    } else if (!hasAdminAccess) {
      return <EnterPin />
    }
    return <FestivalAdminPage />
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center">
      <div className="animate-pulse text-primary text-lg font-medium">Loading...</div>
    </div>
  )
}

function FestivalAdminPage() {
  const [alert, setAlert] = useState("")
  const [inputValue, setInputValue] = useState("")
  const { t } = useTranslation()
  const { festivalCode } = useParams()

  function handleClick(event) {
    event.preventDefault()

    let valid = true
    for (let i = 0; i < inputValue.length; i++) {
      if (!"1234567890".includes(inputValue.at(i))) {
        valid = false
        break
      }
    }
    if (inputValue.length === 0) valid = false

    if (!valid) {
      playAlert("Please only use numbers", setAlert)
    } else {
      onSetGaugePressed(inputValue)
      setInputValue("")
    }
  }

  function handleInputValue(event) {
    const value = event.target.value
    if ("1234567890".includes(value.at(-1)) || value === "") {
      setInputValue(value)
    }
  }

  return (
    <div className="min-h-screen bg-background p-4 flex items-start justify-center pt-8">
      <Card className="w-full max-w-2xl overflow-hidden shadow-xl border-0">
        {/* Header */}
        <div className="bg-accent px-6 py-5 flex items-center justify-between">
          <Link to="/home">
            <Button variant="ghost" size="icon" className="text-white hover:bg-white/20">
              <Home className="h-5 w-5" />
            </Button>
          </Link>

          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-white/20 rounded-xl flex items-center justify-center">
              <Settings className="w-5 h-5 text-white" />
            </div>
            <h1 className="text-2xl font-bold text-white">
              {t("admin_title")}
            </h1>
          </div>

          <div className="flex items-center gap-2">
            <LanguageSwitcher />
            <Link to={`/${festivalCode}`}>
              <Button variant="ghost" size="icon" className="text-white hover:bg-white/20">
                <ArrowLeft className="h-5 w-5" />
              </Button>
            </Link>
          </div>
        </div>

        <CardContent className="p-6 space-y-6">
          {alert && (
            <Alert variant="destructive">
              <AlertDescription>{alert}</AlertDescription>
            </Alert>
          )}

          {/* Set Gauge Section */}
          <div className="bg-primary/5 rounded-xl p-5 border border-primary/20">
            <h3 className="text-sm font-semibold text-primary uppercase tracking-wide mb-4">
              Set Maximum Capacity
            </h3>
            <form onSubmit={handleClick} className="flex gap-3">
              <Input
                type="text"
                name="maxGauge"
                value={inputValue}
                onChange={handleInputValue}
                placeholder={t("admin_max_gauge")}
                className="text-lg h-12 font-mono"
              />
              <Button type="submit" size="lg" className="px-8 shadow-lg hover:shadow-glow-primary">
                {t("admin_set_gauge")}
              </Button>
            </form>
          </div>

          {/* Current Event Actions */}
          <div className="bg-muted rounded-xl p-5">
            <h3 className="text-sm font-semibold text-foreground uppercase tracking-wide mb-4">
              {t("admin_current_event")}
            </h3>
            <div className="grid grid-cols-2 gap-3">
              <Button
                variant="destructive"
                size="lg"
                onClick={onArchivePressed}
                className="shadow-lg hover:shadow-glow-destructive"
              >
                <Archive className="h-5 w-5 mr-2" />
                {t("admin_archive")}
              </Button>
              <Button
                variant="success"
                size="lg"
                asChild
                className="shadow-lg hover:shadow-glow-success"
              >
                <a href={APIURL + "v1/festival/" + festivalCode + "/admin/download/activecsv"}>
                  <Download className="h-5 w-5 mr-2" />
                  {t("admin_get_csv")}
                </a>
              </Button>
            </div>
          </div>

          {/* Archives */}
          <ArchiveSection festivalCode={festivalCode} t={t} />
        </CardContent>
      </Card>
    </div>
  )

  async function onArchivePressed(event) {
    event.preventDefault()
    const response = await fetch(APIURL + "v1/festival/" + festivalCode + "/admin/archivecurrentevent", {
      method: "get",
      headers: { "Content-Type": "application/json" }
    })
    if (!response.ok) {
      if (response.status === 422) location.reload()
      throw new Error(`Response status:`, response.status)
    }
    location.reload()
  }

  async function onSetGaugePressed(newMax) {
    const body = JSON.stringify({ max: +newMax })
    const response = await fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/admin/setgauge", {
      method: "post",
      headers: { "Content-Type": "application/json" },
      body,
    })
    if (!response.ok) {
      if (response.status === 422) location.reload()
      throw new Error("Response status:", response.status)
    }
  }
}

function ArchiveSection({ festivalCode, t }) {
  const [archives, setArchives] = useState([])

  useEffect(() => {
    fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/admin/getarchivedevents")
      .then(res => res.json())
      .then(data => setArchives(data))
  }, [festivalCode])

  function getDateTime(timestamp) {
    const time = new Date(timestamp * 1000)
    return time.toLocaleString().slice(0, 24)
  }

  return (
    <div className="bg-secondary/50 rounded-xl p-5">
      <h3 className="text-sm font-semibold text-foreground uppercase tracking-wide mb-4">
        {t("admin_archived_data")}
      </h3>

      {archives.length === 0 ? (
        <p className="text-muted-foreground text-sm text-center py-6">
          No archived data available
        </p>
      ) : (
        <div className="space-y-2">
          {archives.map((item) => (
            <a
              key={item.id}
              href={APIURL + "v1/festival/" + festivalCode + "/admin/download/archivedcsv/" + item.id}
              className="flex items-center justify-between p-4 rounded-lg bg-card border hover:border-primary hover:shadow-md transition-all group"
            >
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
                  <FileDown className="w-5 h-5 text-primary" />
                </div>
                <span className="font-semibold">Archive #{item.id}</span>
              </div>
              <span className="text-muted-foreground text-sm font-mono">{getDateTime(item.time)}</span>
            </a>
          ))}
        </div>
      )}
    </div>
  )
}

async function playAlert(alert, setAlert) {
  setAlert(alert)
  await new Promise(resolve => setTimeout(resolve, 7500))
  setAlert("")
}

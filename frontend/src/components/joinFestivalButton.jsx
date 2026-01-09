import { useEffect, useState } from "react"
import { auth, fetchWithAuth } from "../utils/auth"
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Users, LogIn } from "lucide-react";

const APIURL = import.meta.env.VITE_APIURL;

export function JoinPopup({ close }) {
  const { t } = useTranslation()
  useEffect(() => { auth() }, [])

  const [codeInputValue, setCodeInputValue] = useState("")
  const [alert, setAlert] = useState("")
  const navigate = useNavigate()

  function handleCodeInputValue(event) {
    const value = event.target.value
    if (("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ".includes(value.at(-1)) || value === "") && value.length <= 6) {
      setCodeInputValue(value.toUpperCase())
    }
  }

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
    <Dialog open={true} onOpenChange={() => close()}>
      <DialogContent onClose={close} className="sm:max-w-md border-0 shadow-2xl overflow-hidden p-0">
        {/* Colored header */}
        <div className="bg-success px-6 py-6 text-center">
          <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center mx-auto mb-3">
            <Users className="w-7 h-7 text-white" />
          </div>
          <DialogTitle className="text-2xl font-bold text-white">
            {t("joinpopup_header")}
          </DialogTitle>
        </div>

        <form onSubmit={handleJoin} className="p-6 space-y-4">
          {alert && (
            <Alert variant="destructive">
              <AlertDescription>{alert}</AlertDescription>
            </Alert>
          )}

          <div>
            <label className="text-sm font-medium text-muted-foreground mb-2 block">
              Festival Code
            </label>
            <Input
              type="text"
              name="Code"
              value={codeInputValue}
              onChange={handleCodeInputValue}
              placeholder={t("joinpopup_enter_code")}
              className="text-center text-2xl h-14 uppercase tracking-[0.3em] font-mono font-bold"
              maxLength={6}
              autoFocus
            />
          </div>

          <Button type="submit" variant="success" size="xl" className="w-full shadow-lg hover:shadow-glow-success">
            <LogIn className="w-5 h-5 mr-2" />
            {t("joinpopup_confirm")}
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  )
}

import { useEffect, useState } from "react"
import { auth, fetchWithAuth } from "../utils/auth"
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Plus, Eye, EyeOff } from "lucide-react";

const APIURL = import.meta.env.VITE_APIURL;

export function CreatePopup({ close }) {
  useEffect(() => { auth() }, [])

  const [passwordInputValue, setPasswordInputValue] = useState("")
  const [pinInputValue, setPinInputValue] = useState("")
  const [showPassword, setShowPassword] = useState(false)
  const { t } = useTranslation()
  const navigate = useNavigate()

  function handlePinInputValue(event) {
    const value = event.target.value
    if (("1234567890".includes(value.at(-1)) || value === "") && value.length <= 4) {
      setPinInputValue(value)
    }
  }

  function handlePasswordInputValue(event) {
    const value = event.target.value
    if ("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-+_*".includes(value.at(-1)) || value === "") {
      setPasswordInputValue(value)
    }
  }

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
        if (data.type === "festival code" && data.content !== null) {
          navigate("/" + data.content, { replace: true })
        }
      })
  }

  return (
    <Dialog open={true} onOpenChange={() => close()}>
      <DialogContent onClose={close} className="sm:max-w-md border-0 shadow-2xl overflow-hidden p-0">
        {/* Colored header */}
        <div className="bg-primary px-6 py-6 text-center">
          <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center mx-auto mb-3">
            <Plus className="w-7 h-7 text-white" />
          </div>
          <DialogTitle className="text-2xl font-bold text-white">
            {t("createpopup_header")}
          </DialogTitle>
        </div>

        <form onSubmit={handleCreate} className="p-6 space-y-5">
          <input
            type="text"
            name="email"
            autoComplete="username email"
            className="hidden"
          />

          <div>
            <Label htmlFor="pin" className="text-sm font-medium text-muted-foreground mb-2 block">
              {t("createpopup_pin_label")}
            </Label>
            <Input
              id="pin"
              type="text"
              min={0}
              maxLength={4}
              name="pin"
              value={pinInputValue}
              onChange={handlePinInputValue}
              placeholder={t("createpopup_admin_pin")}
              className="text-center text-2xl h-14 tracking-[0.5em] font-mono font-bold"
              autoFocus
            />
          </div>

          <div>
            <Label htmlFor="create-password" className="text-sm font-medium text-muted-foreground mb-2 block">
              {t("createpopup_password_label")}
            </Label>
            <div className="relative">
              <Input
                id="create-password"
                type={showPassword ? 'text' : 'password'}
                maxLength={50}
                autoComplete="new-password"
                name="password"
                value={passwordInputValue}
                onChange={handlePasswordInputValue}
                placeholder={t("createpopup_password")}
                className="h-12 pr-12"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
              >
                {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>
          </div>

          <Button type="submit" size="xl" className="w-full shadow-lg hover:shadow-glow-primary">
            <Plus className="w-5 h-5 mr-2" />
            {t("createpopup_confirm")}
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  )
}

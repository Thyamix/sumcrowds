import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router-dom";
import { fetchWithAuth } from "../utils/auth";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { X, Lock, Eye, EyeOff, LogIn } from "lucide-react";

const APIURL = import.meta.env.VITE_APIURL;

export function EnterPassword() {
  const { t } = useTranslation()
  const [passwordInputValue, setPasswordInputValue] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [passwordError, setPasswordError] = useState(null);
  const { festivalCode } = useParams()
  const navigate = useNavigate()

  function handlePasswordInputValue(event) {
    const value = event.target.value;
    if (
      "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-+_*".includes(value.at(-1)) ||
      value === ""
    ) {
      setPasswordInputValue(value);
      setPasswordError(null);
    }
  }

  async function handleConfirm(event) {
    event.preventDefault();
    setPasswordError(null);

    const body = JSON.stringify({ password: passwordInputValue });

    try {
      const response = await fetchWithAuth(APIURL + "v1/festival/" + festivalCode + "/getaccess", {
        method: "post",
        headers: { "Content-Type": "application/json" },
        body: body,
      });

      if (!response.ok) {
        setPasswordError(t("pwpopup_alert"));
      } else {
        setPasswordInputValue("");
        location.reload()
      }
    } catch (error) {
      setPasswordError(t("error_generic"));
    }
  }

  function goHome() {
    navigate("/", { replace: true })
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <Card className="w-full max-w-md overflow-hidden border-0 shadow-xl">
        {/* Header */}
        <div className="bg-primary px-6 py-6 text-center relative">
          <Button
            variant="ghost"
            size="icon"
            className="absolute top-4 right-4 text-white hover:bg-white/20"
            onClick={goHome}
          >
            <X className="h-5 w-5" />
          </Button>

          <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center mx-auto mb-3">
            <Lock className="w-7 h-7 text-white" />
          </div>
          <h1 className="text-2xl font-bold text-white">
            {t("pwpopup_header")}
          </h1>
          <p className="text-white/70 text-sm mt-1">{t("pwpopup_festival")}: {festivalCode}</p>
        </div>

        <CardContent className="p-6">
          <form onSubmit={handleConfirm} className="space-y-5">
            <input
              type="text"
              name="email"
              autoComplete="username"
              className="hidden"
            />

            {passwordError && (
              <Alert variant="destructive">
                <AlertDescription>{passwordError}</AlertDescription>
              </Alert>
            )}

            <div>
              <Label htmlFor="password" className="text-sm font-medium text-muted-foreground mb-2 block">
                {t("pwpopup_password")}
              </Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  maxLength={50}
                  autoComplete="current-password"
                  name="password"
                  value={passwordInputValue}
                  onChange={handlePasswordInputValue}
                  placeholder={t("pwpopup_password")}
                  className="h-12 pr-12"
                  autoFocus
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

            <Button type="submit" variant="success" size="xl" className="w-full shadow-lg hover:shadow-glow-success">
              <LogIn className="w-5 h-5 mr-2" />
              {t("pwpopup_confirm")}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}

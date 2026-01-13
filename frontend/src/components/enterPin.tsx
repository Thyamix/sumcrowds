import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router-dom";
import { fetchWithAuth } from "../utils/auth";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { X, Shield, Eye, EyeOff, LogIn } from "lucide-react";

const APIURL = import.meta.env.VITE_APIURL;

export function EnterPin() {
  const { t } = useTranslation()
  const [pinInputValue, setPinInputValue] = useState("");
  const [showPin, setShowPin] = useState(false);
  const [pinError, setPinError] = useState<string | null>(null);
  const { festivalCode } = useParams()
  const navigate = useNavigate()

  function handlePinInputValue(event: React.ChangeEvent<HTMLInputElement>) {
    const value = event.target.value;
    if (value.length > 4) return
    if ("0123456789".includes(value.at(-1) || "") || value === "") {
      setPinInputValue(value);
      setPinError(null);
    }
  }

  async function handleConfirm(event: React.FormEvent) {
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
      } else {
        setPinInputValue("");
        location.reload()
      }
    } catch (error) {
      setPinError(t("error_generic"));
    }
  }

  function goBack() {
    navigate(`/${festivalCode}`, { replace: true })
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <Card className="w-full max-w-md overflow-hidden border-0 shadow-xl">
        {/* Header */}
        <div className="bg-accent px-6 py-6 text-center relative">
          <Button
            variant="ghost"
            size="icon"
            className="absolute top-4 right-4 text-white hover:bg-white/20"
            onClick={goBack}
          >
            <X className="h-5 w-5" />
          </Button>

          <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center mx-auto mb-3">
            <Shield className="w-7 h-7 text-white" />
          </div>
          <h1 className="text-2xl font-bold text-white">
            {t("pinpopup_header")}
          </h1>
          <p className="text-white/70 text-sm mt-1">{t("pinpopup_admin_access")}</p>
        </div>

        <CardContent className="p-6">
          <form onSubmit={handleConfirm} className="space-y-5">
            <input
              type="text"
              name="email"
              autoComplete="username"
              className="hidden"
            />

            {pinError && (
              <Alert variant="destructive">
                <AlertDescription>{pinError}</AlertDescription>
              </Alert>
            )}

            <div>
              <Label htmlFor="pin" className="text-sm font-medium text-muted-foreground mb-2 block">
                {t("pinpopup_label")}
              </Label>
              <div className="relative">
                <Input
                  id="pin"
                  type={showPin ? "text" : "password"}
                  maxLength={4}
                  name="pin"
                  value={pinInputValue}
                  onChange={handlePinInputValue}
                  placeholder="••••"
                  className="text-center text-3xl h-16 tracking-[0.3em] pr-12"
                  autoFocus
                />
                <button
                  type="button"
                  onClick={() => setShowPin(!showPin)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                >
                  {showPin ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                </button>
              </div>
            </div>

            <Button type="submit" variant="default" size="xl" className="w-full bg-accent hover:bg-accent/90 shadow-lg">
              <LogIn className="w-5 h-5 mr-2" />
              {t("pinpopup_confirm")}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}

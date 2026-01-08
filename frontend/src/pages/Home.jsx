import { useState } from "react";
import { useTranslation } from "react-i18next";
import { LanguageSwitcher } from "../components/languageSwitcher";
import { JoinPopup } from "../components/joinFestivalButton";
import { CreatePopup } from "../components/createFestivalButton";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Users, Plus } from "lucide-react";

export function Home() {
  const { t } = useTranslation();
  const [joinOpen, setJoinOpen] = useState(false);
  const [createOpen, setCreateOpen] = useState(false);

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <Card className="w-full max-w-md overflow-hidden shadow-xl border-0">
        {/* Colored Header */}
        <div className="bg-primary px-6 py-8 text-center relative">
          <div className="absolute top-4 right-4">
            <LanguageSwitcher />
          </div>
          <div className="w-16 h-16 bg-white/20 rounded-2xl flex items-center justify-center mx-auto mb-4">
            <Users className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-white">
            {t("home_home")}
          </h1>
        </div>

        <CardContent className="p-6 space-y-6">
          <div className="text-center">
            <h2 className="text-xl font-semibold text-foreground mb-2">
              {t("home_welcome")}
            </h2>
            <p className="text-muted-foreground">
              {t("home_select_option")}
            </p>
          </div>

          <div className="space-y-3">
            <Button
              variant="success"
              size="xl"
              className="w-full shadow-lg hover:shadow-glow-success transition-shadow"
              onClick={() => setJoinOpen(true)}
            >
              <Users className="w-5 h-5 mr-2" />
              {t("home_join_button")}
            </Button>

            <Button
              variant="default"
              size="xl"
              className="w-full shadow-lg hover:shadow-glow-primary transition-shadow"
              onClick={() => setCreateOpen(true)}
            >
              <Plus className="w-5 h-5 mr-2" />
              {t("home_create_button")}
            </Button>
          </div>
        </CardContent>
      </Card>

      {joinOpen && <JoinPopup close={() => setJoinOpen(false)} />}
      {createOpen && <CreatePopup close={() => setCreateOpen(false)} />}
    </div>
  );
}

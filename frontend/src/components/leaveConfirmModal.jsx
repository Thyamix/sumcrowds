import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogContent, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { AlertTriangle } from 'lucide-react';

export function LeaveConfirmModal({ open, onClose }) {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const handleConfirm = () => {
    onClose();
    navigate('/home');
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent onClose={onClose} className="sm:max-w-md border-0 shadow-2xl overflow-hidden p-0">
        <div className="bg-warning px-6 py-6 text-center">
          <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center mx-auto mb-3">
            <AlertTriangle className="w-7 h-7 text-white" />
          </div>
          <DialogTitle className="text-2xl font-bold text-white">
            {t("leave_confirm_title")}
          </DialogTitle>
        </div>
        <div className="p-6 space-y-4">
          <DialogDescription className="text-center text-base text-muted-foreground">
            {t("leave_confirm_message")}
          </DialogDescription>
          <DialogFooter className="flex gap-3 sm:justify-center">
            <Button
              variant="outline"
              size="lg"
              onClick={onClose}
              className="flex-1"
            >
              {t("leave_confirm_cancel")}
            </Button>
            <Button
              variant="destructive"
              size="lg"
              onClick={handleConfirm}
              className="flex-1"
            >
              {t("leave_confirm_leave")}
            </Button>
          </DialogFooter>
        </div>
      </DialogContent>
    </Dialog>
  );
}

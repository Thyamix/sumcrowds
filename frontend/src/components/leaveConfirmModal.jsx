import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogContent, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';

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
        <div className="bg-primary px-6 py-5 text-center">
          <DialogTitle className="text-xl font-semibold text-white">
            {t("leave_confirm_title")}
          </DialogTitle>
        </div>
        <div className="px-6 pt-4 pb-6">
          <DialogDescription className="text-center text-base text-muted-foreground mb-4">
            {t("leave_confirm_message")}
          </DialogDescription>
          <DialogFooter className="flex flex-row gap-3 justify-center">
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

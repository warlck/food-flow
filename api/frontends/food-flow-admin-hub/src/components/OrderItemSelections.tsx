import type { AdminOrderItem } from '@/lib/admin-api';

function formatCurrency(amount: number) {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount);
}

export function OrderItemSelections({ item }: { item: AdminOrderItem }) {
  return (
    <>
      {item.modifiers && item.modifiers.length > 0 && (
        <div className="mt-1 space-y-0.5">
          {item.modifiers.map((modifier) => (
            <div key={modifier.id} className="flex items-center justify-between gap-3 text-[11px] text-[#6B7280]">
              <span>+ {modifier.modifierOptionName} ({modifier.modifierGroupName})</span>
              <span className="text-[#FF4500]">
                {modifier.priceDelta === 0 ? '+$0.00' : `+${formatCurrency(modifier.priceDelta * item.quantity)}`}
              </span>
            </div>
          ))}
        </div>
      )}
      {item.addons && item.addons.length > 0 && (
        <div className="mt-1 space-y-0.5">
          {item.addons.map((addon) => (
            <div key={addon.id} className="flex items-center justify-between gap-3 text-[11px] text-[#6B7280]">
              <span>+ {addon.addonName} ×{addon.quantity}</span>
              <span className="text-[#FF4500]">
                +{formatCurrency(addon.addonPrice * addon.quantity * item.quantity)}
              </span>
            </div>
          ))}
        </div>
      )}
    </>
  );
}

import { useRef, useState } from 'react';
import { ImageOff, Loader2, Trash2, Upload } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  IMAGE_UPLOAD_ACCEPT,
  IMAGE_UPLOAD_MAX_BYTES,
  ImageEntityType,
  adminApi,
} from '@/lib/admin-api';

interface ImageFieldProps {
  /** Form field name; the resolved URL is submitted through a hidden input. */
  name: string;
  /** Restaurant the upload is scoped to. When undefined (e.g. a brand-new
   * restaurant that has no id yet), uploads are disabled and only URL
   * pasting is available. */
  restaurantId?: string;
  entityType: ImageEntityType;
  defaultValue?: string;
  variant?: 'default' | 'avatar' | 'banner';
  aspectHint?: string;
  placeholder?: string;
}

/**
 * ImageField lets admins upload an image directly to storage via a signed
 * URL minted by the sales API, while keeping the option to paste an
 * external URL. The final URL flows into the surrounding form through a
 * hidden input so existing save handlers stay unchanged.
 */
export function ImageField({
  name,
  restaurantId,
  entityType,
  defaultValue,
  variant = 'default',
  aspectHint,
  placeholder = 'https://images.unsplash.com/photo.jpg',
}: ImageFieldProps) {
  const [url, setUrl] = useState(defaultValue ?? '');
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState('');
  // Tracks an image uploaded during this dialog session so replacing or
  // removing it can delete the tracked object instead of orphaning it.
  const [trackedUpload, setTrackedUpload] = useState<{ imageId: string; publicUrl: string } | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const uploadsEnabled = Boolean(restaurantId);

  const dropTrackedUpload = (tracked: { imageId: string } | null) => {
    if (tracked) adminApi.deleteImage(tracked.imageId).catch(() => undefined);
  };

  const onUrlChange = (value: string) => {
    setUrl(value);
    setError('');
    if (trackedUpload && value !== trackedUpload.publicUrl) {
      dropTrackedUpload(trackedUpload);
      setTrackedUpload(null);
    }
  };

  const onFileSelected = async (file: File | undefined) => {
    if (!file) return;
    if (!IMAGE_UPLOAD_ACCEPT.split(',').includes(file.type)) {
      setError('Only JPEG, PNG, or WebP images are allowed.');
      return;
    }
    if (file.size > IMAGE_UPLOAD_MAX_BYTES) {
      setError(`Image must be ${Math.round(IMAGE_UPLOAD_MAX_BYTES / (1024 * 1024))} MB or smaller.`);
      return;
    }

    setUploading(true);
    setError('');
    try {
      const image = await adminApi.uploadEntityImage({ restaurantId: restaurantId!, entityType, file });
      dropTrackedUpload(trackedUpload);
      setTrackedUpload({ imageId: image.imageId, publicUrl: image.publicUrl });
      setUrl(image.publicUrl);
      toast.success('Image uploaded');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Image upload failed.';
      setError(message);
      toast.error(message);
    } finally {
      setUploading(false);
      if (fileInputRef.current) fileInputRef.current.value = '';
    }
  };

  const onRemove = () => {
    dropTrackedUpload(trackedUpload);
    setTrackedUpload(null);
    setUrl('');
    setError('');
  };

  if (variant === 'banner') {
    return (
      <div className="space-y-3">
        <input type="hidden" name={name} value={url} readOnly />

        {/* Widescreen Banner Preview Container */}
        <div className="relative flex h-36 w-full items-center justify-center overflow-hidden rounded-2xl border-2 border-dashed border-[#F3E1D9] bg-[#FFF7F3] shadow-inner transition-all hover:border-[#FF8C42]">
          {url ? (
            <img src={url} alt="Cover Banner Preview" className="h-full w-full object-cover" />
          ) : (
            <div className="flex flex-col items-center gap-1.5 text-center text-[#9CA3AF]">
              <ImageOff size={24} className="text-[#D6C0B4]" />
              <span className="text-xs font-semibold text-[#6B7280]">No cover banner uploaded</span>
              <span className="text-[11px] text-[#9CA3AF]">{aspectHint ?? 'Recommended: 1200×400px widescreen photo'}</span>
            </div>
          )}
          {url && (
            <div className="absolute bottom-2 right-2 flex gap-1.5 rounded-lg bg-black/60 p-1 backdrop-blur-md">
              <span className="px-2 py-0.5 text-[10px] font-medium text-white">Cover Banner</span>
            </div>
          )}
        </div>

        {/* Controls Toolbar */}
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
          <Input
            value={url}
            onChange={(event) => onUrlChange(event.target.value)}
            placeholder={placeholder}
            className="admin-input flex-1 text-xs"
            aria-label={`${name} url`}
          />
          <div className="flex shrink-0 items-center gap-2">
            <input
              ref={fileInputRef}
              type="file"
              accept={IMAGE_UPLOAD_ACCEPT}
              className="hidden"
              onChange={(event) => void onFileSelected(event.target.files?.[0])}
            />
            <Button
              type="button"
              variant="outline"
              size="sm"
              disabled={!uploadsEnabled || uploading}
              onClick={() => fileInputRef.current?.click()}
              className="h-9 gap-1.5 rounded-xl border-[#E5E7EB] bg-white px-3 text-xs font-semibold text-[#374151] hover:border-[#FF8C42] hover:bg-[#FFF7F3] hover:text-[#FF4500]"
            >
              {uploading ? <Loader2 size={14} className="animate-spin text-[#FF4500]" /> : <Upload size={14} className="text-[#FF4500]" />}
              {uploading ? 'Uploading…' : 'Upload Banner'}
            </Button>
            {url && (
              <Button type="button" variant="ghost" size="sm" onClick={onRemove} className="h-9 gap-1.5 rounded-xl px-2.5 text-xs text-[#B91C1C] hover:bg-[#FEF2F2]">
                <Trash2 size={14} />
                Remove
              </Button>
            )}
          </div>
        </div>

        {!uploadsEnabled && (
          <p className="text-[11px] text-[#9CA3AF]">Save restaurant first to enable direct bucket uploads, or paste an image URL.</p>
        )}
        {error && <p className="text-xs font-medium text-[#B91C1C]">{error}</p>}
      </div>
    );
  }

  if (variant === 'avatar') {
    return (
      <div className="space-y-3">
        <input type="hidden" name={name} value={url} readOnly />

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
          {/* Square 1:1 Avatar Preview Container */}
          <div className="relative flex h-24 w-24 shrink-0 items-center justify-center overflow-hidden rounded-2xl border-2 border-dashed border-[#F3E1D9] bg-[#FFF7F3] shadow-sm transition-all hover:border-[#FF8C42]">
            {url ? (
              <img src={url} alt="Brand Logo Preview" className="h-full w-full object-cover" />
            ) : (
              <div className="flex flex-col items-center text-[#9CA3AF]">
                <ImageOff size={20} className="text-[#D6C0B4]" />
                <span className="mt-1 text-[10px] font-semibold text-[#6B7280]">Logo 1:1</span>
              </div>
            )}
          </div>

          {/* Controls */}
          <div className="flex flex-1 flex-col gap-2">
            <Input
              value={url}
              onChange={(event) => onUrlChange(event.target.value)}
              placeholder={placeholder}
              className="admin-input text-xs"
              aria-label={`${name} url`}
            />
            <div className="flex flex-wrap items-center gap-2">
              <input
                ref={fileInputRef}
                type="file"
                accept={IMAGE_UPLOAD_ACCEPT}
                className="hidden"
                onChange={(event) => void onFileSelected(event.target.files?.[0])}
              />
              <Button
                type="button"
                variant="outline"
                size="sm"
                disabled={!uploadsEnabled || uploading}
                onClick={() => fileInputRef.current?.click()}
                className="h-8 gap-1.5 rounded-lg border-[#E5E7EB] bg-white px-3 text-xs font-semibold text-[#374151] hover:border-[#FF8C42] hover:bg-[#FFF7F3] hover:text-[#FF4500]"
              >
                {uploading ? <Loader2 size={13} className="animate-spin text-[#FF4500]" /> : <Upload size={13} className="text-[#FF4500]" />}
                {uploading ? 'Uploading…' : 'Upload Logo'}
              </Button>
              {url && (
                <Button type="button" variant="ghost" size="sm" onClick={onRemove} className="h-8 gap-1.5 rounded-lg px-2 text-xs text-[#B91C1C] hover:bg-[#FEF2F2]">
                  <Trash2 size={13} />
                  Remove
                </Button>
              )}
              <span className="text-[11px] text-[#9CA3AF]">{aspectHint ?? 'Recommended: 400×400px square PNG/SVG'}</span>
            </div>
            {!uploadsEnabled && (
              <p className="text-[11px] text-[#9CA3AF]">Save restaurant first to enable direct uploads, or paste a URL.</p>
            )}
          </div>
        </div>

        {error && <p className="text-xs font-medium text-[#B91C1C]">{error}</p>}
      </div>
    );
  }

  // Default compact variant
  return (
    <div className="space-y-2">
      <input type="hidden" name={name} value={url} readOnly />

      <div className="flex items-start gap-3">
        <div className="flex h-20 w-28 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-[#F3E1D9] bg-[#FFF7F3]">
          {url ? (
            <img src={url} alt="Preview" className="h-full w-full object-cover" />
          ) : (
            <ImageOff size={18} className="text-[#D6C0B4]" />
          )}
        </div>

        <div className="flex flex-1 flex-col gap-2">
          <Input
            value={url}
            onChange={(event) => onUrlChange(event.target.value)}
            placeholder={placeholder}
            className="admin-input text-xs"
            aria-label={`${name} url`}
          />
          <div className="flex items-center gap-2">
            <input
              ref={fileInputRef}
              type="file"
              accept={IMAGE_UPLOAD_ACCEPT}
              className="hidden"
              onChange={(event) => void onFileSelected(event.target.files?.[0])}
            />
            <Button
              type="button"
              variant="outline"
              size="sm"
              disabled={!uploadsEnabled || uploading}
              onClick={() => fileInputRef.current?.click()}
              className="gap-1.5"
            >
              {uploading ? <Loader2 size={14} className="animate-spin text-[#FF4500]" /> : <Upload size={14} className="text-[#FF4500]" />}
              {uploading ? 'Uploading…' : 'Upload'}
            </Button>
            {url && (
              <Button type="button" variant="ghost" size="sm" onClick={onRemove} className="gap-1.5 text-[#B91C1C]">
                <Trash2 size={14} />
                Remove
              </Button>
            )}
          </div>
          {!uploadsEnabled && (
            <p className="text-[11px] text-[#9CA3AF]">Save first to enable direct uploads, or paste a URL.</p>
          )}
        </div>
      </div>

      {error && <p className="text-xs font-medium text-[#B91C1C]">{error}</p>}
    </div>
  );
}


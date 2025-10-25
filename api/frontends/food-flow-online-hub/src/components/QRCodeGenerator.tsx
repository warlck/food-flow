
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { QrCode, Download, ExternalLink } from 'lucide-react';

const QRCodeGenerator: React.FC = () => {
  const mobileMenuUrl = `${window.location.origin}/mobile-menu`;
  
  // For demo purposes, using a QR code API service
  const qrCodeUrl = `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(mobileMenuUrl)}`;

  const handleDownload = () => {
    const link = document.createElement('a');
    link.href = qrCodeUrl;
    link.download = 'menu-qr-code.png';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const handleOpenMobileMenu = () => {
    window.open(mobileMenuUrl, '_blank');
  };

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <QrCode className="w-5 h-5" />
          QR Menu Code
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex justify-center">
          <img
            src={qrCodeUrl}
            alt="QR Code for Mobile Menu"
            className="border-2 border-gray-200 rounded-lg"
          />
        </div>
        
        <div className="text-center space-y-2">
          <p className="text-sm text-gray-600">
            Customers can scan this QR code to access the mobile menu
          </p>
          <p className="text-xs text-gray-500 font-mono break-all">
            {mobileMenuUrl}
          </p>
        </div>

        <div className="flex gap-2">
          <Button 
            onClick={handleDownload}
            variant="outline" 
            className="flex-1"
          >
            <Download className="w-4 h-4 mr-2" />
            Download
          </Button>
          <Button 
            onClick={handleOpenMobileMenu}
            className="flex-1"
          >
            <ExternalLink className="w-4 h-4 mr-2" />
            Preview
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

export default QRCodeGenerator;

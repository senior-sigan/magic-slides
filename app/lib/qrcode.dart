import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import 'package:slides/fetch.dart';
import 'package:slides/slider.dart';

const prefix = "slides://";

class QRCodePage extends StatelessWidget {
  const QRCodePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Slides Controller')),
      body: MobileScanner(
        controller: MobileScannerController(
          detectionTimeoutMs: 750,
          detectionSpeed: DetectionSpeed.normal,
        ),
        onDetect: (final capture) {
          var value = capture.barcodes[0].rawValue;
          if (value != null && value.startsWith(prefix)) {
            var code = value.substring(prefix.length);
            debugPrint('Barcode found! $code');  
            var res = dio.post("/api/start", data: {"code": code});
            res.then((value) {
                debugPrint("onStart");
                Navigator.push(
                  context, 
                  MaterialPageRoute(builder: (context) => SliderPage(code: code))
                );
            });
          }
        },
      ),
    );
  }
}

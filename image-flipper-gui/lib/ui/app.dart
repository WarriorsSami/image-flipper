import 'package:flutter/material.dart';
import 'package:image_flipper_gui/ui/pages/image_flipper.dart';

class ImageFlipperApp extends StatelessWidget {
  const ImageFlipperApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Image Flipper',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: Colors.lightGreen,
          secondary: Colors.lightBlue,
        ),
        useMaterial3: true,
      ),
      home: const ImageFlipperPage(),
    );
  }
}

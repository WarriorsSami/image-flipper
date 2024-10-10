import 'package:flutter/material.dart';

class ImageFlipperPage extends StatelessWidget {
  const ImageFlipperPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Image Flipper'),
      ),
      body: const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            Text('Flip your images!'),
          ],
        ),
      ),
    );
  }
}

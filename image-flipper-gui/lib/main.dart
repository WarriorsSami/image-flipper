import 'package:flutter/material.dart';
import 'package:image_flipper_gui/di.dart';
import 'package:image_flipper_gui/ui/app.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  setupDi();
  runApp(const ImageFlipperApp());
}

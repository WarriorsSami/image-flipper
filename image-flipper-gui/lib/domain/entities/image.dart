import 'dart:io';

import 'package:path/path.dart' as p;

class Image {
  final String path;
  final String name;

  Image({required this.path, required this.name});
  factory Image.file(File file) {
    final name = p.basename(file.path);

    return Image(
      path: file.path,
      name: name,
    );
  }
}

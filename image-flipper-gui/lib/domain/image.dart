import 'dart:io';

class Image {
  final String path;
  final String name;

  Image({required this.path, required this.name});
  factory Image.file(File file) {
    return Image(
      path: file.path,
      name: file.path.split('/').last,
    );
  }
}
